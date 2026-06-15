package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// UnknownVariableValue well named
const UnknownVariableValue = "74D93920-ED26-11E3-AC10-0800200C9A66"

func Expand(m map[string]string, key string) interface{} {
	if v, ok := m[key]; ok {
		if v == "true" {
			return true
		} else if v == "false" {
			return false
		}
		return v
	}
	if v, ok := m[key+".#"]; ok {
		if v == UnknownVariableValue {
			return v
		}
		return expandArray(m, key)
	}
	prefix := key + "."
	for k := range m {
		if strings.HasPrefix(k, prefix) {
			return expandMap(m, prefix)
		}
	}
	return nil
}

func expandArray(m map[string]string, prefix string) []interface{} {
	num, err := strconv.ParseInt(m[prefix+".#"], 0, 0)
	if err != nil {
		panic(err)
	}
	if num == 0 {
		return []interface{}{}
	}
	keySet := map[int]bool{}
	computed := map[string]bool{}
	for k := range m {
		if !strings.HasPrefix(k, prefix+".") {
			continue
		}
		key := k[len(prefix)+1:]
		idx := strings.Index(key, ".")
		if idx != -1 {
			key = key[:idx]
		}
		if key == "#" {
			continue
		}
		if strings.HasPrefix(key, "~") {
			key = key[1:]
			computed[key] = true
		}
		k, err := strconv.Atoi(key)
		if err != nil {
			panic(err)
		}
		keySet[k] = true
	}
	keysList := make([]int, 0, num)
	for key := range keySet {
		keysList = append(keysList, key)
	}
	sort.Ints(keysList)
	result := make([]interface{}, len(keysList))
	for i, key := range keysList {
		keyString := strconv.Itoa(key)
		if computed[keyString] {
			keyString = "~" + keyString
		}
		result[i] = Expand(m, fmt.Sprintf("%s.%s", prefix, keyString))
	}
	return result
}

func expandMap(m map[string]string, prefix string) map[string]interface{} {
	if count, ok := m[prefix+"%"]; ok && count == "0" {
		return map[string]interface{}{}
	}
	result := make(map[string]interface{})
	for k := range m {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		key := k[len(prefix):]
		idx := strings.Index(key, ".")
		if idx != -1 {
			key = key[:idx]
		}
		if _, ok := result[key]; ok {
			continue
		}
		if key == "%" {
			continue
		}
		result[key] = Expand(m, k[:len(prefix)+len(key)])
	}
	return result
}
