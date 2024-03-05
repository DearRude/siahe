package main

import (
	"context"

	"golang.org/x/time/rate"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"go.uber.org/zap"

	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"

	db "github.com/DearRude/fumCommunityBot/database"
	"github.com/DearRude/fumCommunityBot/handlers"
)

func main() {
	c := GenConfig()
	ctx := context.Background()

	// Init zap logger
	logConf := zap.NewDevelopmentConfig()
	logConf.Level.SetLevel(zap.InfoLevel)
	logger, err := logConf.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer func() {
		_ = logger.Sync()
	}()
	sugar := logger.Sugar()

	// Init database
	database := db.DbConfig{Path: c.SqlitePath}
	if err := database.InitDatabase(); err != nil {
		sugar.Panicf("Error init database: %w", err)
	}

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

	// Run the bot
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

		// Init handlers
		api := tg.NewClient(client)
		handlers.InitHandlers(database.Db, api, message.NewSender(api), c.AdminPassword, c.VarificationChat)

		// Setting up handler for incoming message.
		dispatcher.OnNewMessage(handlers.HandleNewMessage)
		dispatcher.OnBotCallbackQuery(handlers.HandleCallbacks)

		select {}
	}); err != nil {
		sugar.Fatalf("Error running client: %w", err)
	}
}
