package handlers

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"

	in "github.com/DearRude/fumTheatreBot/internals"
)

var (
	StateMap = in.NewUserStateMap()
	UserMap  = in.NewUserDataMap()
	EventMap = in.NewEventDataMap()

	sender           *message.Sender
	client           *tg.Client
	db               *gorm.DB
	adminPassword    string
	varificationChat *tg.InputPeerChat
)

func InitHandlers(database *gorm.DB, tgClient *tg.Client, messageSender *message.Sender, adminPass string, varifChat int) {
	client = tgClient
	sender = messageSender
	db = database
	adminPassword = adminPass
	varificationChat = &tg.InputPeerChat{ChatID: int64(varifChat)}
}

func HandleNewMessage(c context.Context, ent tg.Entities, u *tg.UpdateNewMessage) error {
	m, ok := u.GetMessage().(*tg.Message)
	if !ok || m.Out { // Outgoing message, not interesting.
		return nil
	}

	// Get sender user
	user, err := getSenderUser(m.GetPeerID(), ent)
	if err != nil {
		return err
	}

	updates := in.UpdateMessage{
		Ctx:      c,
		Ent:      ent,
		Unm:      u,
		Message:  m,
		PeerUser: user.AsInputPeer(),
	}

	// If new message is a command
	if command := getCommandName(updates); command != "" {
		if err := handleCommands(updates); err != nil {
			return fmt.Errorf("Error handle command %s: %w", command, err)
		}
	} else { // If not command, handle by state
		if err := handleMessageStates(updates); err != nil {
			return fmt.Errorf("Error handle by state: %w", err)
		}
	}
	return nil
}

func HandleCallbacks(ctx context.Context, ent tg.Entities, update *tg.UpdateBotCallbackQuery) error {
	// Answer read
	_, err := client.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: update.QueryID,
		Message: "حله!",
	})
	if err != nil {
		return err
	}

	// Construct update callback
	user := ent.Users[update.UserID]
	u := in.UpdateCallback{
		Ctx:      ctx,
		Ent:      ent,
		Ubc:      update,
		PeerUser: user.AsInputPeer(),
	}

	// Query from varification chat
	peerChat, isChat := update.GetPeer().(*tg.PeerChat)
	if isChat && (peerChat.GetChatID() == varificationChat.GetChatID()) {
		return varificationChatResponse(u)
	}

	// Query from user
	state, hasState := StateMap.Get(user.GetID())
	if hasState {
		switch state {
		case in.SignUpAskGender:
			return signUpAskGender(u)
		case in.SignUpAskIsFumStudent:
			return signUpAskIsFumStudent(u)
		case in.SignUpAskIsStudent:
			return signUpAskIsStudent(u)
		case in.SignUpAskIsMashhadStudent:
			return signUpAskIsMashhadStudent(u)
		case in.SignUpAskIsGraduate:
			return signUpAskIsGraduate(u)
		case in.SignUpAskIsStudentRelative:
			return signUpAskIsStudentRelative(u)
		case in.SignUpAskFumFaculty:
			return signUpAskFumFaculty(u)
		case in.SignUpAskIsMastPhd:
			return signUpAskMastPhd(u)
		case in.SignUpCheckInfo:
			return signUpCheckInfo(u)
		case in.GetTicketInit:
			return getTicketInit(u)
		}
	}
	return nil
}
