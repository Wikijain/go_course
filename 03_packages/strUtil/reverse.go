package strUtil

func Reverse(name string) string {
	n := len(name)
	runes := []rune(name)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}
	return string(runes)
}
