package handlers

import (
	"gorm.io/gorm"
	"strconv"
	"strings"

	"github.com/DearRude/siahe/database"
	in "github.com/DearRude/siahe/internals"
)

func handleCommands(u in.UpdateMessage) error {
	command := getCommandName(u)
	if command == "" {
		return nil
	}

	StateMap.Set(u.PeerUser.UserID, in.CommandState)
	UserMap.Delete(u.PeerUser.UserID)
	EventMap.Delete(u.PeerUser.UserID)

	// Handle user commands
	switch command {
	case "start":
		return startCommand(u)
	case "add_account":
		return addAccountCommand(u)
	case "delete_account":
		return deleteAccountCommand(u)
	case "get_account":
		return getAccountCommand(u)
	case "promote_me":
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
		case "get_user":
			return getUserCommand(u)
		case "delete_user":
			return deleteUserCommand(u)
		case "add_place":
			return addPlaceCommand(u)
		case "get_place":
			return getPlaceCommand(u)
		case "get_places":
			return getPlacesCommand(u)
		case "delete_place":
			return deletePlaceCommand(u)
		case "add_event":
			return addEventCommand(u)
		case "get_event":
			return getEventCommand(u)
		case "get_events":
			return getEventsCommand(u)
		case "delete_event":
			return deleteEventCommand(u)
		case "activate_event":
			return activateEventCommand(u)
		case "deactivate_event":
			return deactivateEventCommand(u)
		case "get_ticket":
			return getTicketCommand(u)
		case "attend_ticket":
			return attendTicketCommand(u)
		case "unattend_ticket":
			return unattendTicketCommand(u)
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

	// Check deeplink coomands
	params := getCommandParams(u)
	if len(params) == 1 { // deeplink initiated
		command := strings.Split(params[0], "_")[0]
		switch command {
		case "get_ticket":
			return getTicketDeepLink(u)
		}
	}

	return nil
}

func addAccountCommand(u in.UpdateMessage) error {
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
		}
		return err
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
	targetID, err := parseIDFromParam(u)
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
	targetID, err := parseIDFromParam(u)
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
	targetID, err := parseIDFromParam(u)
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
	targetID, err := parseIDFromParam(u)
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
	targetID, err := parseIDFromParam(u)
	if err != nil {
		return err
	}

	var place database.Place
	if err := db.First(&place, targetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No event found
			if err := reactToMessage(u, "👎"); err != nil {
				return err
			}
		}
		return err
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
	targetID, err := parseIDFromParam(u)
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

// line parameters: name, description, isPaid, maxTicketBatch, placeID
func addEventCommand(u in.UpdateMessage) error {
	params := getCommandLines(u)
	if len(params) != 5 { // only five parameters
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageAddEventeHelp()...); err != nil {
			return err
		}

		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageAddEventExample()...); err != nil {
			return err
		}
		return nil
	}

	isPaid, err := parseBoolFromText(params[2])
	if err != nil {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return err
	}

	maxTicketBatch, err := strconv.ParseUint(params[3], 10, 64)
	if err != nil {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return err
	}

	placeID, err := strconv.ParseUint(params[4], 10, 64)
	if err != nil {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return err
	}

	// Check if place exists
	var place database.Place
	if err := db.First(&place, placeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := reactToMessage(u, "👎"); err != nil {
				return err
			}
		}
		return err
	}

	event := database.Event{
		Name:           params[0],
		Description:    params[1],
		MaxTicketBatch: uint(maxTicketBatch),
		IsPaid:         isPaid,
		PlaceID:        uint(placeID),
		IsActive:       true, // An event is a active by default
	}

	res := db.Create(&event)
	if err := res.Error; err != nil {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return err
	}

	if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageEventAdded(event)...); err != nil {
		return err
	}
	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}

	return nil
}

// parameter: eventID
func getEventCommand(u in.UpdateMessage) error {
	targetID, err := parseIDFromParam(u)
	if err != nil {
		return err
	}

	var event database.Event
	if err := db.First(&event, targetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No event found
			if err := reactToMessage(u, "👎"); err != nil {
				return err
			}
		}
		return err
	}

	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}
	_, err = sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessagePrintEvent(event)...)
	return err
}

// parameter: eventID
func activateEventCommand(u in.UpdateMessage) error {
	targetID, err := parseIDFromParam(u)
	if err != nil {
		return err
	}

	res := db.Model(&database.Event{}).Where("id = ?", targetID).Update("is_active", true)
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

// parameter: eventID
func deactivateEventCommand(u in.UpdateMessage) error {
	targetID, err := parseIDFromParam(u)
	if err != nil {
		return err
	}

	res := db.Model(&database.Event{}).Where("id = ?", targetID).Update("is_active", false)
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

func getEventsCommand(u in.UpdateMessage) error {
	var events []database.Event
	if err := db.Find(&events).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No event found
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
	_, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessagePrintEvents(events)...)
	return err
}

// parameter: eventID
func deleteEventCommand(u in.UpdateMessage) error {
	targetID, err := parseIDFromParam(u)
	if err != nil {
		return err
	}

	res := db.Delete(&database.Event{}, targetID)
	if err := res.Error; err != nil {
		return err
	}

	// No event found
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

// parameter: ticketID
func getTicketCommand(u in.UpdateMessage) error {
	targetID, err := parseIDFromParam(u)
	if err != nil {
		return err
	}

	var ticket database.Ticket
	if err := db.Preload("User").Preload("Event").First(&ticket, targetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No event found
			if err := reactToMessage(u, "👎"); err != nil {
				return err
			}
		}
		return err
	}

	if err := reactToMessage(u, "👍"); err != nil {
		return err
	}
	_, err = sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessagePrintTicket(ticket)...)
	return err
}

// parameter: ticketID
func attendTicketCommand(u in.UpdateMessage) error {
	targetID, err := parseIDFromParam(u)
	if err != nil {
		return err
	}

	res := db.Model(&database.Ticket{}).Where("id = ?", targetID).Update("status", "attended")
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

// parameter: ticketID
func unattendTicketCommand(u in.UpdateMessage) error {
	targetID, err := parseIDFromParam(u)
	if err != nil {
		return err
	}

	res := db.Model(&database.Ticket{}).Where("id = ?", targetID).Update("status", "completed")
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
