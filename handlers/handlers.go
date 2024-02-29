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

	sender *message.Sender
	client *tg.Client
	db     *gorm.DB
)

func InitHandlers(database *gorm.DB, tgClient *tg.Client, messageSender *message.Sender) {
	client = tgClient
	sender = messageSender
	db = database
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
	if command := getCommandName(m); command != "" {
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

func HandleCallbacks(ctx context.Context, ent tg.Entities, u *tg.UpdateBotCallbackQuery) error {
	// Get sender user
	user, err := getSenderUser(u.GetPeer(), ent)
	if err != nil {
		return err
	}

	_, err = client.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: u.QueryID,
		Message: "حله!",
	})

	updates := in.UpdateCallback{
		Ctx:      ctx,
		Ent:      ent,
		Ubc:      u,
		PeerUser: user.AsInputPeer(),
	}

	state, hasState := StateMap.Get(user.GetID())
	if hasState {
		switch state {
		case in.SignUpAskGender:
			return signUpAskGender(updates)
		case in.SignUpAskIsFumStudent:
			return signUpAskIsFumStudent(updates)
		case in.SignUpAskIsStudent:
			return signUpAskIsStudent(updates)
		case in.SignUpAskIsMashhadStudent:
			return signUpAskIsMashhadStudent(updates)
		case in.SignUpAskIsGraduate:
			return signUpAskIsGraduate(updates)
		case in.SignUpAskIsStudentRelative:
			return signUpAskIsStudentRelative(updates)
		case in.SignUpAskFumFaculty:
			return signUpAskFumFaculty(updates)
		case in.SignUpAskIsMastPhd:
			return signUpAskMastPhd(updates)
		case in.SignUpCheckInfo:
			return signUpCheckInfo(updates)
		default:
			return nil
		}
	}
	return nil
}
