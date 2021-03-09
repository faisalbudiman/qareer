package utils

func MapStringToString(s []string, fn func(string) string) []string {
	newS := make([]string, len(s))
	for i, v := range s {
		newS[i] = fn(v)
	}

	return newS
}

func MapStringToInterface(s []string, fn func(string) interface{}) []interface{} {
	newS := make([]interface{}, len(s))
	for i, v := range s {
		newS[i] = fn(v)
	}

	return newS
}
