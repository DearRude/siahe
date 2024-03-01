package handlers

import (
	"gorm.io/gorm"

	"github.com/DearRude/fumTheatreBot/database"
	in "github.com/DearRude/fumTheatreBot/internals"
)

func handleCommands(u in.UpdateMessage) error {
	command := getCommandName(u)
	if command == "" {
		return nil
	}

	StateMap.Set(u.PeerUser.UserID, in.CommandState)

	switch command {
	case "start":
		return startCommand(u)
	case "signup":
		return signupCommand(u)
	case "deleteAccount":
		return deleteAccountCommand(u)
	case "getAccount":
		return getAccountCommand(u)
	case "promoteMe":
		return promoteMeCommand(u)
	default:
		return nil
	}
}

func startCommand(u in.UpdateMessage) error {
	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}
	if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageStart(u.PeerUser.UserID)...); err != nil {
		return err
	}
	return nil
}

func signupCommand(u in.UpdateMessage) error {
	var user database.User
	if err := db.Model(&database.User{}).First(&user, u.PeerUser.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Accepted. Ask for first name
			if err := reactToMessage(u, "👍"); err != nil {
				return err
			}
			if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageAskFirstName()...); err != nil {
				return err
			}
			StateMap.Set(u.PeerUser.UserID, in.SignUpAskFirstName)
			return nil
		} else {
			return err
		}
	}

	// User is found. They can't sign up again

	if err := reactToMessage(u, "👎"); err != nil {
		return err
	}
	_, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageYouAlreadySignedUp(user.FirstName)...)
	return err
}

func deleteAccountCommand(u in.UpdateMessage) error {
	res := db.Delete(&database.User{}, u.PeerUser.UserID)
	if err := res.Error; err != nil {
		return err
	}

	// User does not have an account
	if res.RowsAffected <= 0 {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageUserHasNoAccount()...); err != nil {
			return err
		}
		return nil
	}

	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}
	_, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageAccountDeleted()...)
	return err
}

func getAccountCommand(u in.UpdateMessage) error {
	var user database.User
	if err := db.Model(&database.User{}).First(&user, u.PeerUser.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User has no account
			if err := reactToMessage(u, "👎"); err != nil {
				return err
			}
			_, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageUserHasNoAccount()...)
			return err
		} else {
			return err
		}
	}

	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}
	_, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessagePrintUser(user)...)
	return err
}

// parameters: password
func promoteMeCommand(u in.UpdateMessage) error {
	params := getCommandParams(u)
	if len(params) != 1 { // only one parameter
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return nil
	}
	givenPass := params[0]

	if givenPass != adminPassword {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return nil
	}

	res := db.Model(&database.User{}).Where("id = ?", u.PeerUser.UserID).Update("role", "mod")
	if err := res.Error; err != nil {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return nil
	}

	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}

	return nil
}
