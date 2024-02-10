package main

import (
	"regexp"
)

func IsStringPersian(text string) bool {
	persianRegex := regexp.MustCompile(`[\x{0600}-\x{06FF}]+`)
	return persianRegex.MatchString(text)
}
