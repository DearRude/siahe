package handlers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gotd/td/tg"

	in "github.com/DearRude/fumTheatreBot/internals"
)

func getCommandName(u in.UpdateMessage) string {
	text := u.Message.GetMessage()
	if len(text) <= 0 || text[0] != '/' {
		return ""
	}
	return strings.Split(text, " ")[0][1:]
}

func getCommandParams(u in.UpdateMessage) []string {
	text := u.Message.GetMessage()
	return strings.Split(text, " ")[1:]
}

func getCommandLines(u in.UpdateMessage) []string {
	text := u.Message.GetMessage()
	return strings.Split(text, "\n\n")[1:]
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
