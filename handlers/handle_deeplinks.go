package handlers

import (
	"gorm.io/gorm"
	"strconv"
	"strings"

	"github.com/DearRude/fumTheatreBot/database"
	in "github.com/DearRude/fumTheatreBot/internals"
)

// parameter: eventID
func getTicketDeepLink(u in.UpdateMessage) error {
	// Check if user has account
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

	// Get target event ID
	param := strings.Split(getCommandParams(u)[0], "_")[1]
	targetID, err := strconv.Atoi(param)
	if err != nil {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return err
	}

	// Get Event
	var event database.Event
	if err := db.Preload("Place").First(&event, targetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No event found
			if err := reactToMessage(u, "👎"); err != nil {
				return err
			}
			if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageInvalidTicketLink()...); err != nil {
				return err
			}
		}
		return err
	}

	// Check user has already reserved request
	var resereved int64
	if err := db.Model(&database.Ticket{}).Where("user_id = ? AND event_id = ? AND status = ?", u.PeerUser.UserID, event.ID, "reserved").
		Count(&resereved).Error; err != nil {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return err
	}
	if resereved > 0 {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageTicketAlreadyReserving()...); err != nil {
			return err
		}
		return nil
	}

	// Check event being active
	if !event.IsActive {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageEventIsDeactive()...); err != nil {
			return err
		}
		return err
	}

	// Calculate reamining tickets for user
	var userTicketsBought int64
	if err := db.Model(&database.Ticket{}).Where("user_id = ? AND event_id = ?", u.PeerUser.UserID, event.ID).Count(&userTicketsBought).Error; err != nil {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return err
	}
	ticketRemains := int64(event.MaxTicketBatch) - userTicketsBought
	if ticketRemains <= 0 {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageMaxTicketIsReached()...); err != nil {
			return err
		}
		return err
	}
	event.MaxTicketBatch = uint(ticketRemains)

	// Find event free capacity
	var ticketsBought int64
	if err := db.Model(&database.Ticket{}).Where("event_id = ?", event.ID).Count(&ticketsBought).Error; err != nil {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		return err
	}
	capacityRemains := int64(event.Place.Capacity) - ticketsBought
	if capacityRemains <= 0 {
		if err := reactToMessage(u, "👎"); err != nil {
			return err
		}
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageEventIsFull()...); err != nil {
			return err
		}
		return err
	}

	// Set buyable tickets
	if event.MaxTicketBatch > uint(capacityRemains) {
		event.MaxTicketBatch = uint(capacityRemains)
	}
	EventMap.Set(u.PeerUser.UserID, event)

	// Next state: generate ticket
	if _, err := sender.Reply(u.Ent, u.Unm).Row(in.ButtonYesNo()...).StyledText(u.Ctx, in.MessageWantToGetTicket(event, event.Place, int(event.MaxTicketBatch))...); err != nil {
		return err
	}
	StateMap.Set(u.PeerUser.UserID, in.GetTicketInit)

	return nil
}
