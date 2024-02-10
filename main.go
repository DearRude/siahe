package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

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

type UserState uint

const (
	CommandState UserState = iota
	SignUpAskFirstName
	SignUpAskLastName
	SignUpAskPhoneNumber
)

var (
	userStates  = make(map[int64]UserState)
	statesMutex = sync.Mutex{}
	sender      *message.Sender
	db          *gorm.DB
)

func writeToState(userId int64, state UserState) {
	statesMutex.Lock()
	defer statesMutex.Unlock()
	userStates[userId] = state
}

func readFromState(userId int64) (UserState, bool) {
	statesMutex.Lock()
	defer statesMutex.Unlock()
	val, ok := userStates[userId]
	return val, ok
}

func deleteFromState(userId int64) {
	statesMutex.Lock()
	defer statesMutex.Unlock()
	delete(userStates, userId)
}

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

		select {}
	}); err != nil {
		sugar.Fatalf("Error running client: %w", err)
	}
}

func handleNewMessage(c context.Context, ent tg.Entities, u *tg.UpdateNewMessage) error {
	m, ok := u.Message.(*tg.Message)
	if !ok || m.Out {
		// Outgoing message, not interesting.
		return nil
	}

	// If new message is a command
	if command := getCommandName(m); command != "" {
		if err := handleCommands(c, ent, u); err != nil {
			return fmt.Errorf("Error handle command %s: %w", command, err)
		}
	} else { // If not command, handle by state
		if err := handleStates(c, ent, u); err != nil {
			return fmt.Errorf("Error handle by state: %w", err)
		}
	}

	// Sending reply.
	// formats := []message.StyledTextOption{
	// 	styling.Plain("plaintext"), styling.Plain("\n\n"),
	// 	styling.Mention("@durov"), styling.Plain("\n\n"),
	// 	styling.Hashtag("#hashtag"), styling.Plain("\n\n"),
	// 	styling.BotCommand("/command"), styling.Plain("\n\n"),
	// 	styling.URL("https://google.org"), styling.Plain("\n\n"),
	// 	styling.Email("example@example.org"), styling.Plain("\n\n"),
	// 	styling.Bold("bold"), styling.Plain("\n\n"),
	// 	styling.Italic("italic"), styling.Plain("\n\n"),
	// 	styling.Underline("underline"), styling.Plain("\n\n"),
	// 	styling.Strike("strike"), styling.Plain("\n\n"),
	// 	styling.Code("fmt.Println(`Hello, World!`)"), styling.Plain("\n\n"),
	// 	styling.Pre("fmt.Println(`Hello, World!`)", "Go"), styling.Plain("\n\n"),
	// 	styling.TextURL("clickme", "https://google.com"), styling.Plain("\n\n"),
	// 	styling.Phone("+71234567891"), styling.Plain("\n\n"),
	// 	styling.Cashtag("$CASHTAG"), styling.Plain("\n\n"),
	// 	styling.Blockquote("blockquote"), styling.Plain("\n\n"),
	// 	styling.BankCard("5550111111111111"), styling.Plain("\n\n"),
	// }

	// _, err := sender.Reply(ent, u).StyledText(c, formats...)
	// _, err := sender.Reply(ent, u).Text(c, m.Message)
	return nil
}

func handleStates(c context.Context, ent tg.Entities, u *tg.UpdateNewMessage) error {
	m, _ := u.Message.(*tg.Message)

	// Get sender user
	user := getSenderUser(m, ent)
	state, hasState := readFromState(user.ID)
	if !hasState {
		return nil // No states found
	}

	switch state {
	case SignUpAskFirstName:
		return signUpAskFirstNameState(c, ent, u)
	case SignUpAskLastName:
		return signUpAskLastNameState(c, ent, u)
	default:
		return nil
	}
}

func signUpAskFirstNameState(c context.Context, ent tg.Entities, u *tg.UpdateNewMessage) error {
	m, _ := u.Message.(*tg.Message)
	text := m.Message
	if text == "" {
		if _, err := sender.Reply(ent, u).StyledText(c, messageHasNoText()...); err != nil {
			return err
		}
	}
	if !IsStringPersian(text) {
		if _, err := sender.Reply(ent, u).StyledText(c, messageIsNotPersian()...); err != nil {
			return err
		}
	}
	return nil
}
func signUpAskLastNameState(c context.Context, ent tg.Entities, u *tg.UpdateNewMessage) error {
	return nil
}

func handleCommands(c context.Context, ent tg.Entities, u *tg.UpdateNewMessage) error {
	m, _ := u.Message.(*tg.Message)
	command := getCommandName(m)
	if command == "" {
		return nil
	}

	// Remove state of the user if a command is invoked
	writeToState(getSenderUser(m, ent).ID, CommandState)

	switch command {
	case "start":
		return startCommand(c, ent, u)
	case "signup":
		return signupCommand(c, ent, u)
	default:
		return nil
	}
}

func startCommand(c context.Context, ent tg.Entities, u *tg.UpdateNewMessage) error {
	if _, err := sender.Reply(ent, u).StyledText(c, messageStart()...); err != nil {
		return err
	}
	return nil
}
func signupCommand(c context.Context, ent tg.Entities, u *tg.UpdateNewMessage) error {
	m, _ := u.Message.(*tg.Message)
	telUser := getSenderUser(m, ent)

	var user User
	if err := db.Model(&User{}).First(&user, telUser.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Accepted. Ask for first name
			if _, err := sender.Reply(ent, u).StyledText(c, messageAskFirstName()...); err != nil {
				return nil
			}
			writeToState(telUser.ID, SignUpAskFirstName)
			return nil
		} else {
			return err
		}
	}

	// User is found. They can't sign up again
	_, err := sender.Reply(ent, u).StyledText(c, messageYouAlreadySignedUp(user.FirstName)...)
	return err
}

func getCommandName(m *tg.Message) string {
	if len(m.Message) <= 0 || m.Message[0] != '/' {
		return ""
	}
	return strings.Split(m.Message, " ")[0][1:]
}

func getSenderUser(m *tg.Message, ent tg.Entities) *tg.User {
	uId, ok := m.FromID.(*tg.PeerUser)
	if !ok {
		for _, user := range ent.Users {
			if user.Self && user.Bot {
				return nil
			}
			return user
		}
	}
	return ent.Users[uId.UserID]
}
