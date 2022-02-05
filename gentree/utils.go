package main

func maxInt(x, y int) int {
	if x > y {
		return x
	}

	return y
}

func minInt(x, y int) int {
	if x < y {
		return x
	}

	return y
}

func containsStr(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}

	return false
}
