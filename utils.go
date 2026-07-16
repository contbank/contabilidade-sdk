package contabilidade

// String returns a pointer to the given string value.
func String(v string) *string {
	return &v
}

// Bool returns a pointer to the given bool value.
func Bool(v bool) *bool {
	return &v
}

// Float64 returns a pointer to the given float64 value.
func Float64(v float64) *float64 {
	return &v
}

// Int64 returns a pointer to the given int64 value.
func Int64(v int64) *int64 {
	return &v
}

// Int32 returns a pointer to the given int32 value.
func Int32(v int32) *int32 {
	return &v
}
