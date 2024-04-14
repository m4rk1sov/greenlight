package data

import (
	"fmt"
	"strconv"
)

// declare Runtime type with int32
type Runtime int32

// implement MarshalJSON() method to satisfy interface
// value receive instead of pointer for versatility
func (r Runtime) MarshalJSON() ([]byte, error) {
	// generate string runtime required format
	jsonValue := fmt.Sprintf("%d mins", r)

	// strconv.Quote() to wrap it in double quotes to be valid for JSON
	quotedJSONValue := strconv.Quote(jsonValue)

	// convert string to byte slice and return value
	return []byte(quotedJSONValue), nil
}
