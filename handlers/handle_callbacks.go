package handlers

import (
	"bytes"
	"fmt"

	in "github.com/DearRude/fumTheatreBot/internals"
)

func signUpAskGender(u in.UpdateCallback) error {
	user, _ := UserMap.Get(u.PeerUser.UserID)

	data, ok := u.Ubc.GetData()
	if !ok {
		return fmt.Errorf("Error getting callback 'gender' data")
	}

	var isGenderBoy bool
	if bytes.Equal(data, []byte("boy")) {
		isGenderBoy = true
	} else if bytes.Equal(data, []byte("girl")) {
		isGenderBoy = false
	} else {
		return fmt.Errorf("Invalid query data is sent")
	}

	user.IsBoy = isGenderBoy
	UserMap.Set(u.PeerUser.UserID, user)

	// Next state: Ask if current FUM student
	if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).Row(in.ButtonYesNo()...).StyledText(u.Ctx, in.MessageAskIsFUMStudent()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskIsFumStudent)

	return nil
}

func signUpAskIsFumStudent(u in.UpdateCallback) error {
	user, _ := UserMap.Get(u.PeerUser.UserID)

	isTrue, err := getYesNoButtonAnswer(u)
	if err != nil {
		return err
	}

	user.IsFumStudent = isTrue
	if isTrue {
		user.IsStudent = true
		user.IsMashhadStudent = true
		user.UniversityName = "دانشگاه فردوسی مشهد"
	}
	UserMap.Set(u.PeerUser.UserID, user)

	if isTrue {
		// Next state: Ask studentNumber
		if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).StyledText(u.Ctx, in.MessageAskStudentNumber()...); err != nil {
			return err
		}
		StateMap.Set(u.PeerUser.UserID, in.SignUpAskStudentNumber)
	} else {

		// Next state: Ask if is student
		if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).Row(in.ButtonYesNo()...).StyledText(u.Ctx, in.MessageAskIsStudent()...); err != nil {
			return err
		}
		StateMap.Set(u.PeerUser.UserID, in.SignUpAskIsStudent)
	}

	return nil
}

func signUpAskIsStudent(u in.UpdateCallback) error {
	user, _ := UserMap.Get(u.PeerUser.UserID)

	isTrue, err := getYesNoButtonAnswer(u)
	if err != nil {
		return err
	}

	user.IsStudent = isTrue
	UserMap.Set(u.PeerUser.UserID, user)

	if isTrue {
		// Next state: Ask is mashhad student
		if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).Row(in.ButtonYesNo()...).StyledText(u.Ctx, in.MessageAskIsMashhad()...); err != nil {
			return err
		}
		StateMap.Set(u.PeerUser.UserID, in.SignUpAskIsMashhadStudent)
	} else {
		// Next state: Ask if is graduate student
		if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).Row(in.ButtonYesNo()...).StyledText(u.Ctx, in.MessageAskIsGraduate()...); err != nil {
			return err
		}
		StateMap.Set(u.PeerUser.UserID, in.SignUpAskIsGraduate)
	}

	return nil
}

func signUpAskIsMashhadStudent(u in.UpdateCallback) error {
	user, _ := UserMap.Get(u.PeerUser.UserID)

	isTrue, err := getYesNoButtonAnswer(u)
	if err != nil {
		return err
	}

	user.IsMashhadStudent = isTrue
	UserMap.Set(u.PeerUser.UserID, user)

	// Next state: Ask university name
	if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).StyledText(u.Ctx, in.MessageAskUniversityName()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskUniversityName)

	return nil
}

func signUpAskIsGraduate(u in.UpdateCallback) error {
	user, _ := UserMap.Get(u.PeerUser.UserID)

	isTrue, err := getYesNoButtonAnswer(u)
	if err != nil {
		return err
	}

	user.IsGraduateStudent = isTrue
	UserMap.Set(u.PeerUser.UserID, user)

	// Next state: Ask student relative
	if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).Row(in.ButtonYesNo()...).StyledText(u.Ctx, in.MessageAskIsRelative()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskIsStudentRelative)

	return nil
}

func signUpAskIsStudentRelative(u in.UpdateCallback) error {
	user, _ := UserMap.Get(u.PeerUser.UserID)

	isTrue, err := getYesNoButtonAnswer(u)
	if err != nil {
		return err
	}

	user.IsStudentRelative = isTrue
	UserMap.Set(u.PeerUser.UserID, user)

	// signUp finished
	if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).StyledText(u.Ctx, in.MessageSignUpFinished(user.FirstName)...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.CommandState)

	return nil
}

func signUpAskFumFaculty(u in.UpdateCallback) error {
	user, _ := UserMap.Get(u.PeerUser.UserID)

	data, ok := u.Ubc.GetData()
	if !ok {
		return fmt.Errorf("Error getting callback data")
	}
	facultyName := string(data)

	user.FumFaculty = &facultyName
	UserMap.Set(u.PeerUser.UserID, user)

	// Next state: Is Master/Phd student
	if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).Row(in.ButtonYesNo()...).StyledText(u.Ctx, in.MessageAskIsMastPhd()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskIsMastPhd)

	return nil
}

func signUpAskMastPhd(u in.UpdateCallback) error {
	user, _ := UserMap.Get(u.PeerUser.UserID)

	isTrue, err := getYesNoButtonAnswer(u)
	if err != nil {
		return err
	}

	user.IsGraduateStudent = isTrue
	UserMap.Set(u.PeerUser.UserID, user)

	// Next state: Ask student major
	if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).StyledText(u.Ctx, in.MessageAskMajor()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskStudentMajor)

	return nil
}
