package main

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/time/rate"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

func main() {
	c := GenConfig()
	ctx := context.Background()

	// Init zap logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// Init database
	database := DbConfig{Path: c.SqlitePath}
	if err := database.InitDatabase(); err != nil {
		sugar.Panicf("Error init database: %w", err)
	}
	db = database.Db

	// Init telegram Client
	dispatcher := tg.NewUpdateDispatcher()
	waiter := floodwait.NewWaiter().WithCallback(func(ctx context.Context, wait floodwait.FloodWait) {
		sugar.Warn("Flood wait", zap.Duration("wait", wait.Duration))
	})

	client := telegram.NewClient(c.AppID, c.AppHash, telegram.Options{
		UpdateHandler:  dispatcher,
		Logger:         logger,
		SessionStorage: &telegram.FileSessionStorage{Path: c.SessionPath},
		Middlewares: []telegram.Middleware{
			ratelimit.New(rate.Every(c.RateLimit), c.RateBurst),
			waiter,
		},
	})

	if err := waiter.Run(ctx, func(ctx context.Context) error {
		stop, err := bg.Connect(client)
		if err != nil {
			sugar.Panicf("Cant connect to Telegram server: %w", err)
		}
		defer func() { _ = stop() }()

		// Authrozation
		if _, err := client.Auth().Bot(ctx, c.BotToken); err != nil {
			sugar.Panicf("Unable to authorize: %w", err)
		}

		api := tg.NewClient(client)
		sender = message.NewSender(api)

		// Setting up handler for incoming message.
		dispatcher.OnNewMessage(handleNewMessage)
		dispatcher.OnBotCallbackQuery(handleCallbacks)

		select {}
	}); err != nil {
		sugar.Fatalf("Error running client: %w", err)
	}
}

func handleCallbacks(ctx context.Context, ent tg.Entities, u *tg.UpdateBotCallbackQuery) error {
	// Get sender user
	user, err := getSenderUser(u.GetPeer(), ent)
	if err != nil {
		return err
	}

	updates := UpdateCallback{
		Ctx:      ctx,
		Ent:      ent,
		Ubc:      u,
		PeerUser: user.AsInputPeer(),
	}

	state, hasState := StateMap.Get(user.GetID())
	if hasState {
		switch state {
		case SignUpAskGender:
			return signUpAskGender(updates)
		default:
			return nil
		}
	}
	return nil
}

func signUpAskGender(u UpdateCallback) error {
	// TODO: Read data from callback query
	return nil
}

func handleNewMessage(c context.Context, ent tg.Entities, u *tg.UpdateNewMessage) error {
	m, ok := u.GetMessage().(*tg.Message)
	if !ok || m.Out { // Outgoing message, not interesting.
		return nil
	}

	updates := UpdateMessage{
		Ctx:     c,
		Ent:     ent,
		Unm:     u,
		Message: m,
	}
	// Get sender user
	user, err := getSenderUser(m.GetPeerID(), ent)
	if err != nil {
		return err
	}
	updates.SenderID = user.GetID()

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

func handleMessageStates(u UpdateMessage) error {
	state, hasState := StateMap.Get(u.SenderID)
	if !hasState {
		return nil // No states found
	}

	switch state {
	case SignUpAskFirstName:
		return signUpAskFirstNameState(u)
	case SignUpAskLastName:
		return signUpAskLastNameState(u)
	case SignUpAskPhoneNumber:
		return signUpAskPhoneNumber(u)
	default:
		return nil
	}
}

func signUpAskFirstNameState(u UpdateMessage) error {
	ok, err := CheckPersianText(u)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// Add user to map
	user := User{FirstName: u.Message.Message}
	UserMap.Set(u.SenderID, user)

	// Next state: Ask last name
	if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, messageAskLastName()...); err != nil {
		return err
	}
	StateMap.Set(u.SenderID, SignUpAskLastName)

	return nil
}

func signUpAskLastNameState(u UpdateMessage) error {
	ok, err := CheckPersianText(u)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// Update user
	user, _ := UserMap.Get(u.SenderID)
	user.LastName = u.Message.Message
	UserMap.Set(u.SenderID, user)

	// Next state: Ask phone number
	if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, messageAskPhone()...); err != nil {
		return err
	}
	StateMap.Set(u.SenderID, SignUpAskPhoneNumber)

	return nil
}

func signUpAskPhoneNumber(u UpdateMessage) error {
	ok, err := CheckPhoneText(u)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// Update user
	user, _ := UserMap.Get(u.SenderID)
	user.PhoneNumber = u.Message.Message
	UserMap.Set(u.SenderID, user)

	// Next state: Ask gender
	// TODO: send callback query button
	StateMap.Set(u.SenderID, SignUpAskGender)

	return nil
}

func handleCommands(u UpdateMessage) error {
	command := getCommandName(u.Message)
	if command == "" {
		return nil
	}

	StateMap.Set(u.SenderID, CommandState)

	switch command {
	case "start":
		return startCommand(u)
	case "signup":
		return signupCommand(u)
	default:
		return nil
	}
}

func startCommand(u UpdateMessage) error {
	if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, messageStart(u.SenderID)...); err != nil {
		return err
	}
	return nil
}
func signupCommand(u UpdateMessage) error {
	var user User
	if err := db.Model(&User{}).First(&user, u.SenderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Accepted. Ask for first name
			if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, messageAskFirstName()...); err != nil {
				return nil
			}
			StateMap.Set(u.SenderID, SignUpAskFirstName)
			return nil
		} else {
			return err
		}
	}

	// User is found. They can't sign up again
	_, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, messageYouAlreadySignedUp(user.FirstName)...)
	return err
}

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
