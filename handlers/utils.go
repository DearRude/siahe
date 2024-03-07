package handlers

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gotd/td/tg"

	"github.com/DearRude/siahe/database"
	in "github.com/DearRude/siahe/internals"
)

func getCommandName(u in.UpdateMessage) string {
	text := u.Message.GetMessage()
	if len(text) <= 0 || text[0] != '/' {
		return ""
	}
	return strings.Split(strings.Split(text, "\n")[0], " ")[0][1:]
}

func getCommandParams(u in.UpdateMessage) []string {
	text := u.Message.GetMessage()
	return strings.Split(text, " ")[1:]
}

func getCommandLines(u in.UpdateMessage) []string {
	text := u.Message.GetMessage()
	return strings.Split(text, "\n-\n")[1:]
}

func getSenderUser(peer tg.PeerClass, ent tg.Entities) (*tg.User, error) {
	peerUser, ok := peer.(*tg.PeerUser)
	if !ok {
		return nil, fmt.Errorf("peerclass could not reflect to peer user")
	}
	user, ok := ent.Users[peerUser.GetUserID()]
	if !ok {
		return nil, fmt.Errorf("user not found in entities")
	}
	return user, nil
}

func getTextFromContact(u in.UpdateMessage) (string, error) {
	media, ok := u.Message.GetMedia()
	if !ok {
		return "", fmt.Errorf("message media is not present")
	}

	contact, ok := media.(*tg.MessageMediaContact)
	if !ok {
		return "", fmt.Errorf("message media is not a contact")
	}

	return contact.GetPhoneNumber(), nil
}

func getPhotoFromMedia(u in.UpdateMessage) (*tg.InputPhoto, error) {
	media, ok := u.Message.GetMedia()
	if !ok {
		return nil, fmt.Errorf("message media is not present")
	}

	mediaphoto, ok := media.(*tg.MessageMediaPhoto)
	if !ok {
		return nil, fmt.Errorf("message media is not a photo")
	}

	photo, ok := mediaphoto.Photo.(*tg.Photo)
	if !ok {
		return nil, fmt.Errorf("photo media is empty")
	}

	return photo.AsInput(), nil

}

func getYesNoButtonAnswer(u in.UpdateCallback) (bool, error) {
	data, ok := u.Ubc.GetData()
	if !ok {
		return false, fmt.Errorf("Error getting callback data")
	}

	var isTrue bool
	if bytes.Equal(data, []byte("yes")) {
		isTrue = true
	} else if bytes.Equal(data, []byte("no")) {
		isTrue = false
	} else {
		return false, fmt.Errorf("Invalid query data is sent")
	}
	return isTrue, nil
}

func reactToMessage(u in.UpdateMessage, emoji string) error {
	_, err := sender.To(u.PeerUser).Reaction(u.Ctx, u.Message.GetID(), &tg.ReactionEmoji{
		Emoticon: emoji,
	})
	return err
}

func isCallbackUserAdmin(u in.UpdateCallback) (bool, error) {
	user := database.User{ID: u.PeerUser.UserID}
	res := db.Select("role").First(&user)
	if res.Error != nil || user.Role == "" {
		return false, fmt.Errorf("Error finding the user")
	}
	return user.Role == "admin" || user.Role == "mod", nil
}

func isUserAdmin(u in.UpdateMessage) (bool, error) {
	user := database.User{ID: u.PeerUser.UserID}
	res := db.Select("role").First(&user)
	if res.Error != nil || user.Role == "" {
		return false, fmt.Errorf("Error finding the user")
	}
	return user.Role == "admin" || user.Role == "mod", nil
}

func isUserMod(u in.UpdateMessage) (bool, error) {
	user := database.User{ID: u.PeerUser.UserID}
	res := db.Select("role").First(&user)
	if res.Error != nil || user.Role == "" {
		return false, fmt.Errorf("Error finding the user")
	}
	return user.Role == "mod", nil
}

func parseIDFromParam(u in.UpdateMessage) (int, error) {
	params := getCommandParams(u)
	if len(params) != 1 { // only one parameter
		if err := reactToMessage(u, "👎"); err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("Not one parameter")
	}

	// Get target ID
	targetID, err := strconv.Atoi(params[0])
	if err != nil {
		if err := reactToMessage(u, "👎"); err != nil {
			return 0, err
		}
		return 0, err
	}

	return targetID, nil
}

func parseBoolFromText(text string) (bool, error) {
	switch text {
	case "بله":
		return true, nil
	case "خیر":
		return false, nil
	default:
		return false, fmt.Errorf("could not parse bool from text: %s", text)
	}
}

func generateTicket(u in.UpdateMessage, count int, status string) ([]uint, error) {
	event, _ := EventMap.Get(u.PeerUser.UserID)
	var ticketIDs []uint

	for i := 0; i < count; i++ {
		id := uint(rand.Intn(9000) + 1000) // generate ID between 1000 and 9999
		for db.First(&database.Ticket{}, id).RowsAffected != 0 {
			id = uint(rand.Intn(9000) + 1000)
		}

		ticket := database.Ticket{
			ID:           id,
			PurchaseTime: time.Now().UTC(),
			Status:       status,
			UserID:       u.PeerUser.UserID,
			EventID:      event.ID,
		}

		if err := db.Create(&ticket).Error; err != nil {
			return nil, err
		}
		ticketIDs = append(ticketIDs, id)
	}

	return ticketIDs, nil
}

func getMessageFromCallback(u in.UpdateCallback) (*tg.Message, error) {
	messagesClass, err := client.MessagesGetMessages(u.Ctx, []tg.InputMessageClass{&tg.InputMessageID{ID: u.Ubc.MsgID}})
	if err != nil {
		return nil, err
	}

	messagesMessages, ok := messagesClass.(*tg.MessagesMessages)
	if !ok {
		return nil, fmt.Errorf("no messagesMessages found")
	}

	message, ok := messagesMessages.Messages[0].(*tg.Message)
	if !ok {
		return nil, fmt.Errorf("message is not a normal message")
	}

	return message, nil
}

func getUserEventIDFromVarification(u in.UpdateCallback) (*tg.InputPeerUser, uint, error) {
	// Get eventID and userID
	message, err := getMessageFromCallback(u)
	if err != nil {
		return nil, 0, err
	}

	messageParams := strings.Split(strings.Split(message.GetMessage(), "\n")[0], " ")
	if len(messageParams) != 3 {
		return nil, 0, fmt.Errorf("invalid varification message")
	}

	eventID, err := strconv.ParseUint(messageParams[0], 10, 64)
	if err != nil {
		return nil, 0, err
	}

	userID, err := strconv.ParseInt(messageParams[1], 10, 64)
	if err != nil {
		return nil, 0, err
	}

	accessHash, err := strconv.ParseInt(messageParams[2], 10, 64)
	if err != nil {
		return nil, 0, err
	}

	peer := &tg.InputPeerUser{
		UserID:     userID,
		AccessHash: accessHash,
	}

	return peer, uint(eventID), nil
}

func deleteMessage(u in.UpdateCallback) error {
	_, err := client.MessagesDeleteMessages(u.Ctx, &tg.MessagesDeleteMessagesRequest{
		ID:     []int{u.Ubc.MsgID},
		Revoke: true,
	},
	)

	return err
}
