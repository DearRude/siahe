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
	AppID    int    //  Get from https://my.telegram.org/apps
	AppHash  string //  Get from https://my.telegram.org/apps
	BotToken string //  Get from https://t.me/BotFather

	// Optionals
	SqlitePath  string
	SessionPath string
	RateLimit   time.Duration
	RateBurst   int
}

func GenConfig() Config {
	log.Println("Read configurations.")
	fs := flag.NewFlagSet("fumTheatreBot", flag.ContinueOnError)
	var (
		appId    = fs.Int("appId", 0, "AppID in Telegram API")
		appHash  = fs.String("appHash", "", "AppHash in Telegram API")
		botToken = fs.String("botToken", "", "BotToken given by BotFather bot")

		sqlitePath  = fs.String("sqlitePath", "./assets/sqlite.db", "relative or absloute path of sqlite db")
		sessionPath = fs.String("sessionPath", "./assets/session.json", "relative or absloute path of session auth file")

		rateLimit = fs.Duration("rateLimit", time.Millisecond*100, "limit maximum rpc call rate")
		rateBurst = fs.Int("rateBurst", 3, "limit rpc call burst")
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

		SqlitePath:  *sqlitePath,
		SessionPath: *sessionPath,

		RateLimit: *rateLimit,
		RateBurst: *rateBurst,
	}
}