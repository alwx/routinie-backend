package utils

func matchesAny(s string, xs []string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

func OneOf(value string, avilableValues []string, defaultValue string) string {
	if matchesAny(value, avilableValues) {
		return value
	}
	return defaultValue
}
