package handlers

import (
	"regexp"

	in "github.com/DearRude/fumTheatreBot/internals"
)

func IsStringPersian(text string) bool {
	// TODO: not include persian numbers
	persianRegex := regexp.MustCompile(`[\x{0600}-\x{06FF}]+`)
	return persianRegex.MatchString(text)
}

func IsStringPhoneNumber(text string) bool {
	phoneRegex := regexp.MustCompile(`^09\d{9}$`)
	return phoneRegex.MatchString(text)
}

func IsStringStudentNumber(text string) bool {
	phoneRegex := regexp.MustCompile(`^[49]\d{9}$`)
	return phoneRegex.MatchString(text)
}

func IsStringEntraceYear(text string) bool {
	phoneRegex := regexp.MustCompile(`^1[3-4][5-9][0-9]$`)
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

func CheckEntranceYear(u in.UpdateMessage) (bool, error) {
	text := u.Message.GetMessage()
	if text == "" {
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageHasNoText()...); err != nil {
			return false, err
		}
		return false, nil
	}
	if !IsStringEntraceYear(text) {
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageIsNotEntranceYear()...); err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func CheckPhoneText(u in.UpdateMessage) (bool, error) {
	text := u.Message.GetMessage()
	if text == "" { // photos, videos, sticker
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

func CheckStudentNumber(u in.UpdateMessage) (bool, error) {
	text := u.Message.GetMessage()
	if text == "" { // photos, videos, sticker
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageHasNoText()...); err != nil {
			return false, err
		}
		return false, nil
	}
	if !IsStringStudentNumber(text) {
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageIsNotStudentNumber()...); err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
