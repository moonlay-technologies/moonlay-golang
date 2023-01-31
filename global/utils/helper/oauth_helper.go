package helper

import "strings"

func GetAuthorizationValue(token string) string {
	authorizationToken := strings.Split(token, " ")
	return authorizationToken[1]
}
