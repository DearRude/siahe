package handlers

import (
	"gorm.io/gorm"
	"strconv"

	"github.com/DearRude/fumTheatreBot/database"
	in "github.com/DearRude/fumTheatreBot/internals"
)

func handleCommands(u in.UpdateMessage) error {
	command := getCommandName(u)
	if command == "" {
		return nil
	}

	StateMap.Set(u.PeerUser.UserID, in.CommandState)

	// Handle user commands
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
	}

	// Handle mod commands
	ok, err := isUserMod(u)
	if err != nil {
		return err
	}
	if ok {
		switch command {
		case "promote":
			return promoteCommand(u)
		case "demote":
			return demoteCommand(u)
		}
	}

	// Handle admin commands
	ok, err = isUserAdmin(u)
	if err != nil {
		return err
	}
	if ok {
		switch command {
		case "getUser":
			return getUserCommand(u)
		case "deleteUser":
			return deleteUserCommand(u)
		case "addPlace":
			return addPlaceCommand(u)
		case "getPlace":
			return getPlaceCommand(u)
		case "getPlaces":
			return getPlacesCommand(u)
		case "deletePlace":
			return deletePlaceCommand(u)
		}
	}

	return nil
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
	if err := db.First(&user, u.PeerUser.UserID).Error; err != nil {
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
	if err := db.First(&user, u.PeerUser.UserID).Error; err != nil {
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

// parameters: userID
func promoteCommand(u in.UpdateMessage) error {
	targetID, err := getIDFromParam(u)
	if err != nil {
		return err
	}

	res := db.Model(&database.User{}).Where("id = ?", targetID).Update("role", "admin")
	if err := res.Error; err != nil || res.RowsAffected <= 0 {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return err
	}

	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}

	return nil
}

// parameters: userID
func demoteCommand(u in.UpdateMessage) error {
	targetID, err := getIDFromParam(u)
	if err != nil {
		return err
	}

	res := db.Model(&database.User{}).Where("id = ?", targetID).Update("role", "user")
	if err := res.Error; err != nil || res.RowsAffected <= 0 {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return err
	}

	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}

	return nil
}

// parameter: userID
func getUserCommand(u in.UpdateMessage) error {
	targetID, err := getIDFromParam(u)
	if err != nil {
		return err
	}

	var user database.User
	if err := db.First(&user, targetID).Error; err != nil {
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
	_, err = sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessagePrintUser(user)...)
	return err
}

// parameter: userID
func deleteUserCommand(u in.UpdateMessage) error {
	targetID, err := getIDFromParam(u)
	if err != nil {
		return err
	}

	res := db.Delete(&database.User{}, targetID)
	if err := res.Error; err != nil {
		return err
	}

	// No user found
	if res.RowsAffected <= 0 {
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

// line parameters: name, address, capacity
func addPlaceCommand(u in.UpdateMessage) error {
	params := getCommandLines(u)
	if len(params) != 3 { // only three parameters
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageAddPlaceHelp()...); err != nil {
			return err
		}

		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageAddPlaceExample()...); err != nil {
			return err
		}
		return nil
	}

	capacity, err := strconv.ParseUint(params[2], 10, 64)
	if err != nil {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return err
	}

	place := database.Place{
		Name:     params[0],
		Address:  params[1],
		Capacity: uint(capacity),
	}

	res := db.Create(&place)
	if err := res.Error; err != nil {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return err
	}

	if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessagePlaceAdded(place)...); err != nil {
		return err
	}
	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}

	return nil
}

// parameter: placeID
func getPlaceCommand(u in.UpdateMessage) error {
	targetID, err := getIDFromParam(u)
	if err != nil {
		return err
	}

	var place database.Place
	if err := db.First(&place, targetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No place found
			if err := reactToMessage(u, "👎"); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}
	_, err = sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessagePrintPlace(place)...)
	return err
}

func getPlacesCommand(u in.UpdateMessage) error {
	var places []database.Place
	if err := db.Find(&places).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No place found
			if err := reactToMessage(u, "👎"); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}
	_, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessagePrintPlaces(places)...)
	return err
}

// parameter: placeID
func deletePlaceCommand(u in.UpdateMessage) error {
	targetID, err := getIDFromParam(u)
	if err != nil {
		return err
	}

	res := db.Delete(&database.Place{}, targetID)
	if err := res.Error; err != nil {
		return err
	}

	// No place found
	if res.RowsAffected <= 0 {
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
