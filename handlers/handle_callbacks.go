package handlers

import (
	"bytes"
	"fmt"

	"github.com/DearRude/fumTheatreBot/database"
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

	// Assign already-defined variables
	user.IsFumStudent = isTrue
	if isTrue {
		user.IsStudent = true
		user.IsMashhadStudent = true
		user.UniversityName = "دانشگاه فردوسی مشهد"
		user.IsGraduateStudent = false
		user.IsStudentRelative = false
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

	// Next state: Check info
	if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).Row(in.ButtonYesNo()...).StyledText(u.Ctx, in.MessageIsUserInfoCorrect(user)...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpCheckInfo)

	return nil
}

func signUpAskFumFaculty(u in.UpdateCallback) error {
	user, _ := UserMap.Get(u.PeerUser.UserID)

	data, ok := u.Ubc.GetData()
	if !ok {
		return fmt.Errorf("Error getting callback data")
	}

	user.FumFaculty = string(data)
	UserMap.Set(u.PeerUser.UserID, user)

	// Next state: Ask entrance year
	if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).StyledText(u.Ctx, in.MessageAskEntranceYear()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskEntraceYear)

	return nil
}

func signUpAskMastPhd(u in.UpdateCallback) error {
	user, _ := UserMap.Get(u.PeerUser.UserID)

	isTrue, err := getYesNoButtonAnswer(u)
	if err != nil {
		return err
	}

	user.IsMastPhd = isTrue
	UserMap.Set(u.PeerUser.UserID, user)

	// Next state: Ask student major
	if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).StyledText(u.Ctx, in.MessageAskMajor()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.SignUpAskStudentMajor)

	return nil
}

func signUpCheckInfo(u in.UpdateCallback) error {
	user, _ := UserMap.Get(u.PeerUser.UserID)

	isTrue, err := getYesNoButtonAnswer(u)
	if err != nil {
		return err
	}

	StateMap.Set(u.PeerUser.UserID, in.CommandState)

	// Cancel signup
	if !isTrue {
		if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).StyledText(u.Ctx, in.MessageCancelSignUp()...); err != nil {
			return err
		}
		return nil
	}

	// Add user to db
	res := db.Create(&user)
	if res.Error != nil {
		return err
	}

	// Finish singUp
	if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).StyledText(u.Ctx, in.MessageSignUpFinished(user.FirstName)...); err != nil {
		return err
	}

	return nil
}

func getTicketInit(u in.UpdateCallback) error {
	isTrue, err := getYesNoButtonAnswer(u)
	if err != nil {
		return err
	}

	// Cancel getTicket
	if !isTrue {
		if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).StyledText(u.Ctx, in.MessageGetTicketCancelled()...); err != nil {
			return err
		}
		return nil
	}

	if _, err := sender.To(u.PeerUser).Reply(u.Ubc.GetMsgID()).StyledText(u.Ctx, in.MessageAskTicketCount()...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.GetTicketCount)

	return nil
}

func varificationChatResponse(u in.UpdateCallback) error {
	// Check if user is admin
	isAdmin, err := isCallbackUserAdmin(u)
	if err != nil {
		return err
	}
	if !isAdmin {
		return nil
	}

	// Get eventID and userID
	peer, eventID, err := getUserEventIDFromVarification(u)
	if err != nil {
		return err
	}
	condition := db.Where("event_id = ? AND user_id = ? AND status = ?", eventID, peer.UserID, "reserved")

	// Check button
	accept, err := getYesNoButtonAnswer(u)
	if err != nil {
		return err
	}

	if accept {
		var tickets []uint
		if err := db.Model(&database.Ticket{}).Where("event_id = ? AND user_id = ? AND status = ?", eventID, peer.UserID, "reserved").Pluck("id", &tickets).Error; err != nil {
			return err
		}
		if err := condition.Model(&database.Ticket{}).Update("status", "completed").Error; err != nil {
			return err
		}
		if _, err := sender.To(peer).StyledText(u.Ctx, in.MessageTicketsBought(tickets)...); err != nil {
			return err
		}
	} else {
		if err := condition.Delete(&database.Ticket{}).Error; err != nil {
			return err
		}
		if _, err := sender.To(peer).StyledText(u.Ctx, in.MessageTicketNotAccepted()...); err != nil {
			return err
		}
	}

	// Delete the message
	if err := deleteMessage(u); err != nil {
		return err
	}

	return nil
}
