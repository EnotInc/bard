package services

func IsLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func IsNumber(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func IsLetterOrNumber(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_' || ('0' <= ch && ch <= '9')
}
