package service

import (
	"strconv"
)

// parseInt преобразует строку в int
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// intToString преобразует int в строку
func intToString(i int) string {
	return strconv.Itoa(i)
}
