package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/peterbourgon/ff/v3"
)

type Config struct {
	// Requireds
	AppID            int    //  Get from https://my.telegram.org/apps
	AppHash          string //  Get from https://my.telegram.org/apps
	BotToken         string //  Get from https://t.me/BotFather
	AdminPassword    string
	VarificationChat int
	BackupChat       int
	GotenbergURL     string

	// Optionals
	SqlitePath  string
	SessionPath string
	RateLimit   time.Duration
	RateBurst   int
}

func GenConfig() Config {
	log.Println("Read configurations.")
	fs := flag.NewFlagSet("siahe", flag.ContinueOnError)
	var (
		appId    = fs.Int("appId", 0, "AppID in Telegram API")
		appHash  = fs.String("appHash", "", "AppHash in Telegram API")
		botToken = fs.String("botToken", "", "BotToken given by BotFather bot")

		gotenbergUrl = fs.String("gotenbergUrl", "", "url of html to pdf service")

		sqlitePath  = fs.String("sqlitePath", "./assets/sqlite.db", "relative or absloute path of sqlite db")
		sessionPath = fs.String("sessionPath", "./assets/session.json", "relative or absloute path of session auth file")

		rateLimit = fs.Duration("rateLimit", time.Millisecond*100, "limit maximum rpc call rate")
		rateBurst = fs.Int("rateBurst", 3, "limit rpc call burst")

		adminPassword    = fs.String("adminPassword", "", "top-admin password")
		varificationChat = fs.Int("varificationChat", 0, "chatID where varification of paid tickets will be sent to")
		backupChat       = fs.Int("backupChat", 0, "chatID where backup file of db will be sent periodiacally")
	)

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fs.String("config", "", "config file")
	} else {
		fs.String("config", ".env", "config file")
	}

	err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVars(),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.EnvParser),
	)
	if err != nil {
		log.Fatalf("Unable to parse args. Error: %s", err)
	}

	return Config{
		AppID:    *appId,
		AppHash:  *appHash,
		BotToken: *botToken,

		GotenbergURL: *gotenbergUrl,

		SqlitePath:  *sqlitePath,
		SessionPath: *sessionPath,

		RateLimit: *rateLimit,
		RateBurst: *rateBurst,

		AdminPassword:    *adminPassword,
		VarificationChat: *varificationChat,
		BackupChat:       *backupChat,
	}
}
