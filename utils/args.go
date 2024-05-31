package utils

import "strings"

func GetArgs(i string) []string {
	return strings.Split(i, " ")
}
