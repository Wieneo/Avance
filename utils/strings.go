package utils

import "strings"

//IsEmpty returns true if string only consists of tabs / spaces
func IsEmpty(v string) bool{
	v = strings.Trim(v, " ")
	v = strings.Trim(v, "	")
	if len(v) == 0{
		return true
	}
	return false
}