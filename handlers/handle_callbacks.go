package handlers

import (
	"bytes"
	"fmt"

	in "github.com/DearRude/fumTheatreBot/internals"
)

func signUpAskGender(u in.UpdateCallback) error {
	// Update user
	user, _ := UserMap.Get(u.PeerUser.UserID)

	data, ok := u.Ubc.GetData()
	if !ok {
		return fmt.Errorf("Error getting callback 'gender' data")
	}

	var isGenderBoy bool
	if bytes.Equal(data, []byte("SignupAskGender_boy")) {
		isGenderBoy = true
	} else if bytes.Equal(data, []byte("SignupAskGender_girl")) {
		isGenderBoy = false
	} else {
		return fmt.Errorf("Invalid query data is sent")
	}

	user.IsBoy = isGenderBoy
	UserMap.Set(u.PeerUser.UserID, user)

	// Next state: Ask phone number
	// if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, messageAskPhone()...); err != nil {
	//	return err
	// }
	// StateMap.Set(u.PeerUser.UserID, SignUpAskPhoneNumber)

	return nil
}
