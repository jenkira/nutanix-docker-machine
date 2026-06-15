package jmespath

import (
	"errors"
	"reflect"
	"unicode"
	"unicode/utf8"
)

/* This is a tree based interpreter.  It walks the AST and directly
   interprets the AST to search through a JSON document.
*/

type treeInterpreter struct {
	fCall *functionCaller
}

func newInterpreter() *treeInterpreter {
	interpreter := treeInterpreter{}
	interpreter.fCall = newFunctionCaller()
	return &interpreter
}

type expRef struct {
	ref ASTNode
}

// Execute takes an ASTNode and input data and interprets the AST directly.
// It will produce the result of applying the JMESPath expression associated
// with the ASTNode to the input data "value".
func (intr *treeInterpreter) Execute(node ASTNode, value interface{}) (interface{}, error) {
	switch node.nodeType {
	case ASTComparator:
		left, err := intr.Execute(node.children[0], value)
		if err != nil {
			return nil, err
		}
		right, err := intr.Execute(node.children[1], value)
		if err != nil {
			return nil, err
		}
		switch node.value {
		case tEQ:
			return objsEqual(left, right), nil
		case tNE:
			return !objsEqual(left, right), nil
		}
		leftNum, ok := left.(float64)
		if !ok {
			return nil, nil
		}
		rightNum, ok := right.(float64)
		if !ok {
			return nil, nil
		}
		switch node.value {
		case tGT:
			return leftNum > rightNum, nil
		case tGTE:
			return leftNum >= rightNum, nil
		case tLT:
			return leftNum < rightNum, nil
		case tLTE:
			return leftNum <= rightNum, nil
		}
	case ASTExpRef:
		return expRef{ref: node}, nil
	case ASTFunctionExpression:
		resolvedArgs := []interface{}{}
		for _, arg := range node.children {
			current, err := intr.Execute(arg, value)
			if err != nil {
				return nil, err
			}
			resolvedArgs = append(resolvedArgs, current)
		}
		return intr.fCall.CallFunction(node.value.(string), resolvedArgs, intr)
	case ASTField:
		if m, ok := value.(map[string]interface{}); ok {
			result := m[node.value.(string)]
			return result, nil
		}
		result := intr.fieldFromStruct(node.value.(string), value)
		return result, nil
	case ASTFilterProjection:
		left, err := intr.Execute(node.children[0], value)
		if err != nil {
			return nil, err
		}
		slice, ok := left.([]interface{})
		if !ok {
			return nil, nil
		}
		comparison := node.children[1]
		collected := []interface{}{}
		for _, element := range slice {
			result, err := intr.Execute(comparison, element)
			if err != nil {
				return nil, err
			}
			if !isFalse(result) {
				current, err := intr.Execute(node.children[2], element)
				if err != nil {
					return nil, err
				}
				if current != nil {
					collected = append(collected, current)
				}
			}
		}
		return collected, nil
	case ASTFlatten:
		left, err := intr.Execute(node.children[0], value)
		if err != nil {
			return nil, err
		}
		slice, ok := left.([]interface{})
		if !ok {
			// If it is not a slice, then return nil
			return nil, nil
		}
		collected := []interface{}{}
		for _, element := range slice {
			if elementSlice, ok := element.([]interface{}); ok {
				collected = append(collected, elementSlice...)
			} else if element != nil {
				collected = append(collected, element)
			}
		}
		return collected, nil
	case ASTIdentity, ASTCurrentNode:
		return value, nil
	case ASTIndex:
		if slice, ok := value.([]interface{}); ok {
			index := node.value.(int)
			if index < 0 {
				index += len(slice)
			}
			if index < len(slice) && index >= 0 {
				return slice[index], nil
			}
			return nil, nil
		}
		return nil, nil
	case ASTKeyValPair:
		return intr.Execute(node.children[0], value)
	case ASTLiteral:
		return node.value, nil
	case ASTMultiSelectHash:
		if value == nil {
			return nil, nil
		}
		collected := make(map[string]interface{})
		for _, child := range node.children {
			current, err := intr.Execute(child, value)
			if err != nil {
				return nil, err
			}
			key := child.value.(string)
			collected[key] = current
		}
		return collected, nil
	case ASTMultiSelectList:
		if value == nil {
			return nil, nil
		}
		collected := []interface{}{}
		for _, child := range node.children {
			current, err := intr.Execute(child, value)
			if err != nil {
				return nil, err
			}
			collected = append(collected, current)
		}
		return collected, nil
	case ASTOrExpression:
		matched, err := intr.Execute(node.children[0], value)
		if err != nil {
			return nil, err
		}
		if isFalse(matched) {
			matched, err = intr.Execute(node.children[1], value)
			if err != nil {
				return nil, err
			}
		}
		return matched, nil
	case ASTAndExpression:
		matched, err := intr.Execute(node.children[0], value)
		if err != nil {
			return nil, err
		}
		if isFalse(matched) {
			return matched, nil
		}
		return intr.Execute(node.children[1], value)
	case ASTNotExpression:
		matched, err := intr.Execute(node.children[0], value)
		if err != nil {
			return nil, err
		}
		if isFalse(matched) {
			return true, nil
		}
		return false, nil
	case ASTPipe:
		result := value
		var err error
		for _, child := range node.children {
			result, err = intr.Execute(child, result)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case ASTProjection:
		left, err := intr.Execute(node.children[0], value)
		if err != nil {
			return nil, err
		}
		slice, ok := left.([]interface{})
		if !ok {
			return nil, nil
		}
		collected := []interface{}{}
		for _, element := range slice {
			current, err := intr.Execute(node.children[1], element)
			if err != nil {
				return nil, err
			}
			if current != nil {
				collected = append(collected, current)
			}
		}
		return collected, nil
	case ASTSubexpression, ASTIndexExpression:
		left, err := intr.Execute(node.children[0], value)
		if err != nil {
			return nil, err
		}
		return intr.Execute(node.children[1], left)
	case ASTSlice:
		slice, ok := value.([]interface{})
		if !ok {
			return nil, nil
		}
		parts := node.value.([]*int)
		sliceParams := make([]sliceParam, 3)
		for i, part := range parts {
			if part != nil {
				sliceParams[i].Specified = true
				sliceParams[i].N = *part
			}
		}
		return slice(slice, sliceParams)
	case ASTValueProjection:
		left, err := intr.Execute(node.children[0], value)
		if err != nil {
			return nil, err
		}
		valMap, ok := left.(map[string]interface{})
		if !ok {
			return nil, nil
		}
		values := make([]interface{}, len(valMap))
		i := 0
		for _, val := range valMap {
			values[i] = val
			i++
		}
		collected := []interface{}{}
		for _, element := range values {
			current, err := intr.Execute(node.children[1], element)
			if err != nil {
				return nil, err
			}
			if current != nil {
				collected = append(collected, current)
			}
		}
		return collected, nil
	}
	return nil, errors.New("Unknown AST node")
}

func (intr *treeInterpreter) fieldFromStruct(key string, value interface{}) interface{} {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		var current string
		if tag := field.Tag.Get("json"); tag != "" {
			// If the field has a JSON tag, use that as the
			// key name
			current = strings.Split(tag, ",")[0]
		} else {
			// Otherwise, use the field name itself.
			current = field.Name
		}
		if current == key {
			return v.Field(i).Interface()
		}
	}
	return nil
}