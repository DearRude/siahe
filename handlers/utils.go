package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/markup"
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

func getInputPeerChat(peer tg.PeerClass) (*tg.InputPeerChat, error) {
	peerChat, ok := peer.(*tg.PeerChat)
	if !ok {
		return nil, fmt.Errorf("peerclass could not reflect to peer chat")
	}
	return peerChat.AsInput(), nil
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

func seenMessage(u in.UpdateCallback, isAccpeted bool) error {
	seenMessage := "رد شد."
	if isAccpeted {
		seenMessage = "قبول شد."
	}

	inputPeer, err := getInputPeerChat(u.Ubc.GetPeer())
	if err != nil {
		return err
	}

	_, err = client.MessagesEditMessage(u.Ctx, &tg.MessagesEditMessageRequest{
		ID:          u.Ubc.MsgID,
		Peer:        inputPeer,
		ReplyMarkup: markup.InlineRow(markup.Callback(seenMessage, []byte("some_random_thing"))),
	},
	)

	return err
}

func exportUsers(u in.UpdateMessage) (*message.UploadedDocumentBuilder, error) {
	// Query data from the database
	var users []database.User
	res := db.Find(&users)
	if res.Error != nil || len(users) <= 0 {
		return nil, fmt.Errorf("error getting users from db")
	}

	// Create a buffer to hold CSV data
	var buf bytes.Buffer

	// Create a CSV writer
	writer := csv.NewWriter(&buf)

	// Write CSV headers
	headers := []string{
		"ID", "FirstName", "LastName", "Role", "PhoneNumber", "IsBoy",
		"IsFumStudent", "StudentNumber", "FumFaculty", "IsStudent",
		"IsMashhadStudent", "UniversityName", "EntranceYear", "IsMastPhd",
		"StudentMajor", "IsGraduateStudent", "IsStudentRelative",
	}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write headers: %w", err)
	}

	// Write data rows to CSV
	for _, user := range users {
		row := []string{
			strconv.FormatInt(user.ID, 10),
			user.FirstName,
			user.LastName,
			user.Role,
			user.PhoneNumber,
			strconv.FormatBool(user.IsBoy),
			strconv.FormatBool(user.IsFumStudent),
			user.StudentNumber,
			user.FumFaculty,
			strconv.FormatBool(user.IsStudent),
			strconv.FormatBool(user.IsMashhadStudent),
			user.UniversityName,
			user.EntranceYear,
			strconv.FormatBool(user.IsMastPhd),
			user.StudentMajor,
			strconv.FormatBool(user.IsGraduateStudent),
			strconv.FormatBool(user.IsStudentRelative),
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write row: %w", err)
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, err
	}

	up, err := upper.FromBytes(u.Ctx, "users.csv", buf.Bytes())
	if err != nil {
		return nil, err
	}

	return message.UploadedDocument(up).Filename("users.csv").MIME("text/csv"), nil
}

func exportTickets(eventID int, u in.UpdateMessage) (*message.UploadedDocumentBuilder, error) {
	// Query data from the database
	var tickets []database.Ticket
	res := db.Preload("User").Where("event_id = ?", eventID).Order("purchase_time").Find(&tickets)
	if res.Error != nil || len(tickets) <= 0 {
		return nil, fmt.Errorf("error getting users from db")
	}

	// Load Tehran timezone
	location, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return nil, err
	}

	// Create a buffer to hold CSV data
	var buf bytes.Buffer

	// Create a CSV writer
	writer := csv.NewWriter(&buf)

	// Write CSV headers
	headers := []string{"ID", "PurchaseTime", "Status",
		"UserID", "FirstName", "LastName", "Role", "PhoneNumber", "IsBoy",
		"IsFumStudent", "StudentNumber", "FumFaculty", "IsStudent",
		"IsMashhadStudent", "UniversityName", "EntranceYear", "IsMastPhd",
		"StudentMajor", "IsGraduateStudent", "IsStudentRelative",
	}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write headers: %w", err)
	}

	// Write data rows to CSV
	for _, ticket := range tickets {
		row := []string{
			strconv.FormatUint(uint64(ticket.ID), 10),
			ticket.PurchaseTime.In(location).Format(time.RFC3339),
			ticket.Status,
			strconv.FormatInt(ticket.User.ID, 10),
			ticket.User.FirstName,
			ticket.User.LastName,
			ticket.User.Role,
			ticket.User.PhoneNumber,
			strconv.FormatBool(ticket.User.IsBoy),
			strconv.FormatBool(ticket.User.IsFumStudent),
			ticket.User.StudentNumber,
			ticket.User.FumFaculty,
			strconv.FormatBool(ticket.User.IsStudent),
			strconv.FormatBool(ticket.User.IsMashhadStudent),
			ticket.User.UniversityName,
			ticket.User.EntranceYear,
			strconv.FormatBool(ticket.User.IsMastPhd),
			ticket.User.StudentMajor,
			strconv.FormatBool(ticket.User.IsGraduateStudent),
			strconv.FormatBool(ticket.User.IsStudentRelative),
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write row: %w", err)
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, err
	}

	up, err := upper.FromBytes(u.Ctx, "tickets.csv", buf.Bytes())
	if err != nil {
		return nil, err
	}

	return message.UploadedDocument(up).Filename("tickets.csv").MIME("text/csv"), nil
}

func toInputPeerUser(u database.User) tg.InputPeerUser {
	return tg.InputPeerUser{
		UserID:     u.ID,
		AccessHash: u.AccessHash,
	}
}
