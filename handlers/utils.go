package handlers

import (
	"fmt"
	"strings"

	"github.com/gotd/td/tg"

	in "github.com/DearRude/fumTheatreBot/internals"
)

func getCommandName(m *tg.Message) string {
	text := m.GetMessage()
	if len(text) <= 0 || text[0] != '/' {
		return ""
	}
	return strings.Split(text, " ")[0][1:]
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

func reactToMessage(u in.UpdateMessage, emoji string) error {
	_, err := sender.To(u.PeerUser).Reaction(u.Ctx, u.Message.ID, &tg.ReactionEmoji{
		Emoticon: emoji,
	})
	return err
}
