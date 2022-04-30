package burptest

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandString creates a random string of given length
func RandString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = SliceItem(letterRunes)
	}
	return string(b)
}
