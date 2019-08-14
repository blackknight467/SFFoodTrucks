package util

import (
	"errors"
	"fmt"
	"strings"
)

type YesNoBool bool

func (bit *YesNoBool) UnmarshalJSON(data []byte) error {
	asString := strings.ToLower(string(data))
	if asString == "y" || asString == "yes" {
		*bit = true
	} else if asString == "n" || asString == "no" {
		*bit = false
	} else {
		return errors.New(fmt.Sprintf("Boolean unmarshal error: invalid input %s", asString))
	}
	return nil
}
