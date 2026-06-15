package utils

import "time"

func StringPtr(v string) *string { return &v }
func StringValue(v *string) string {
	if v != nil { return *v }
	return ""
}
func BoolPtr(v bool) *bool { return &v }
func BoolValue(v *bool) bool {
	if v != nil { return *v }
	return false
}
func IntPtr(v int) *int { return &v }
func IntValue(v *int) int {
	if v != nil { return *v }
	return 0
}
func Int64Ptr(v int64) *int64 { return &v }
func Int64Value(v *int64) int64 {
	if v != nil { return *v }
	return 0
}
func Float64Ptr(v float64) *float64 { return &v }
func Float64Value(v *float64) float64 {
	if v != nil { return *v }
	return 0
}
func Time(v time.Time) *time.Time { return &v }
func TimeValue(v *time.Time) time.Time {
	if v != nil { return *v }
	return time.Time{}
}
