package helper

import "strings"

func GetUsernameFromEmail(email string) string {
	return strings.Split(email, "@")[0]
}
