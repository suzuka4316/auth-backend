package utils

import (
	"errors"
	"fmt"
	"strings"
)

func FormatError(err string) error {
	fmt.Printf("FormatError():: err %v", err)
	if strings.Contains(err, "email") {
		return errors.New(EmailTaken)
	}
	if strings.Contains(err, "not found") {
		return errors.New(EmailNotFound)
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New(IncorrectPassword)
	}
	return errors.New(IncorrectDetails)
}