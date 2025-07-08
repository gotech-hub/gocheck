package utils

import "log"

func CheckErr(err error) {
	if err != nil {
		log.Println("Error:", err)
	}
}

// Helper function to detect commented-out code (simple heuristic)
func IsCommentedOutCode(comment string) bool {
	// Remove comment markers
	if len(comment) > 2 && (comment[:2] == "//" || comment[:2] == "/*") {
		comment = comment[2:]
	}
	comment = TrimSpace(comment)
	// Heuristic: contains ";", "{", "}", or looks like Go code
	if len(comment) > 0 && (ContainsAny(comment, []string{"=", ";", "{", "}", "func ", "var ", "if ", "for ", "return ", "package ", "import "})) {
		return true
	}
	return false
}

func TrimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}

func ContainsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if len(sub) > 0 && Contains(s, sub) {
			return true
		}
	}
	return false
}

func Contains(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) && (s == sub || (len(s) > len(sub) && (s[:len(sub)] == sub || Contains(s[1:], sub))))
}
