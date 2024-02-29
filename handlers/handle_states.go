package handlers

import (
	"github.com/DearRude/fumTheatreBot/database"
	in "github.com/DearRude/fumTheatreBot/internals"
)

func handleMessageStates(u in.UpdateMessage) error {
	state, hasState := StateMap.Get(u.PeerUser.UserID)
	if !hasState {
		return nil // No states found
	}

	switch state {
	case in.SignUpAskFirstName:
		return signUpAskFirstName(u)
	case in.SignUpAskLastName:
		return signUpAskLastName(u)
	case in.SignUpAskPhoneNumber:
		return signUpAskPhoneNumber(u)
	case in.SignUpAskStudentNumber:
		return signUpAskStudentNumber(u)
	case in.SignUpAskStudentMajor:
		return signUpAskStudentMajor(u)
	case in.SignUpAskUniversityName:
		return signUpAskUniversityName(u)
	case in.SignUpAskEntraceYear:
		return signUpAskEntranceYear(u)
	default:
		return nil
	}
}

func signUpAskFirstName(u in.UpdateMessage) error {
	user := database.User{
		ID:   u.PeerUser.UserID,
		Role: "user",
	}

	ok, err := CheckPersianText(u)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// Add user to map
	user.FirstName = u.Message.Message
	UserMap.Set(u.PeerUser.UserID, user)

	// React ok
	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}

	// Next state: Ask last name
	if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageAskLastName()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskLastName)

	return nil
}

func signUpAskLastName(u in.UpdateMessage) error {
	ok, err := CheckPersianText(u)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// Update user
	user, _ := UserMap.Get(u.PeerUser.UserID)
	user.LastName = u.Message.Message
	UserMap.Set(u.PeerUser.UserID, user)

	// React ok
	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}

	// Next state: Ask phone number
	if _, err := sender.Reply(u.Ent, u.Unm).Markup(in.ButtonAskPhone()).StyledText(u.Ctx, in.MessageAskPhone()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskPhoneNumber)

	return nil
}

func signUpAskPhoneNumber(u in.UpdateMessage) error {
	phone, err := getTextFromContact(u)
	if err != nil {
		return err
	}

	// Update user
	user, _ := UserMap.Get(u.PeerUser.UserID)
	user.PhoneNumber = phone
	UserMap.Set(u.PeerUser.UserID, user)

	// React ok
	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}

	// Next state: Ask gender
	if _, err := sender.Reply(u.Ent, u.Unm).Row(in.ButtonAskGender()...).StyledText(u.Ctx, in.MessageAskGender()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskGender)

	return nil
}

func signUpAskStudentNumber(u in.UpdateMessage) error {
	ok, err := CheckStudentNumber(u)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// Update user
	user, _ := UserMap.Get(u.PeerUser.UserID)
	user.StudentNumber = &u.Message.Message
	UserMap.Set(u.PeerUser.UserID, user)

	// React ok
	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}

	// Next state: Ask FUM faculty
	if _, err := sender.Reply(u.Ent, u.Unm).Markup(in.ButtonAskFumFaculty()).StyledText(u.Ctx, in.MessageAskFumFaculty()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskFumFaculty)

	return nil
}

func signUpAskStudentMajor(u in.UpdateMessage) error {
	ok, err := CheckPersianText(u)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// Update user
	user, _ := UserMap.Get(u.PeerUser.UserID)
	user.StudentMajor = u.Message.Message
	UserMap.Set(u.PeerUser.UserID, user)

	// React ok
	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}

	// signUp finished
	if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageSignUpFinished(user.FirstName)...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.CommandState)

	return nil
}

func signUpAskUniversityName(u in.UpdateMessage) error {
	ok, err := CheckPersianText(u)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// Update user
	user, _ := UserMap.Get(u.PeerUser.UserID)
	user.UniversityName = u.Message.Message
	UserMap.Set(u.PeerUser.UserID, user)

	// React ok
	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}

	// Next state: ask entance year
	if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageAskEntranceYear()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskEntraceYear)

	return nil
}

func signUpAskEntranceYear(u in.UpdateMessage) error {
	ok, err := CheckEntranceYear(u)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// Update user
	user, _ := UserMap.Get(u.PeerUser.UserID)
	user.EntraceYear = &u.Message.Message
	UserMap.Set(u.PeerUser.UserID, user)

	// React ok
	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}

	// Next state: is student master/phd
	if _, err := sender.Reply(u.Ent, u.Unm).Row(in.ButtonYesNo()...).StyledText(u.Ctx, in.MessageAskIsMastPhd()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskIsMastPhd)

	return nil
}
