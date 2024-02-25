package handlers

import (
	"regexp"

	in "github.com/DearRude/fumTheatreBot/internals"
)

func IsStringPersian(text string) bool {
	persianRegex := regexp.MustCompile(`[\x{0600}-\x{06FF}]+`)
	return persianRegex.MatchString(text)
}

func IsStringPhoneNumber(text string) bool {
	phoneRegex := regexp.MustCompile(`^09\d{9}$`)
	return phoneRegex.MatchString(text)
}

func CheckPersianText(u in.UpdateMessage) (bool, error) {
	text := u.Message.GetMessage()
	if text == "" {
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageHasNoText()...); err != nil {
			return false, err
		}
		return false, nil
	}
	if !IsStringPersian(text) {
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageIsNotPersian()...); err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func CheckPhoneText(u in.UpdateMessage) (bool, error) {
	text := u.Message.GetMessage()
	if text == "" {
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageHasNoText()...); err != nil {
			return false, err
		}
		return false, nil
	}
	if !IsStringPhoneNumber(text) {
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageIsNotPhone()...); err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
