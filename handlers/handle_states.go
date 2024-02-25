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
		return signUpAskFirstNameState(u)
	case in.SignUpAskLastName:
		return signUpAskLastNameState(u)
	case in.SignUpAskPhoneNumber:
		return signUpAskPhoneNumber(u)
	default:
		return nil
	}
}

func signUpAskFirstNameState(u in.UpdateMessage) error {
	ok, err := CheckPersianText(u)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// Add user to map
	user := database.User{FirstName: u.Message.Message}
	UserMap.Set(u.PeerUser.UserID, user)

	// Next state: Ask last name
	if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageAskLastName()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskLastName)

	return nil
}

func signUpAskLastNameState(u in.UpdateMessage) error {
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

	// Next state: Ask gender
	if _, err := sender.Reply(u.Ent, u.Unm).Row(in.ButtonAskGender()...).StyledText(u.Ctx, in.MessageAskGender()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskGender)

	return nil
}
