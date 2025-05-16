package utils

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func IntPtr(i int) *int {
	return &i
}