package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/f1monkey/search/internal/log"
	"github.com/spf13/viper"
)

func main() {
	localCfg := flag.String("config", "", "path to config file")
	flag.Parse()

	err := loadConfig(*localCfg)
	if err != nil {
		panic(err)
	}

	logger, err := log.New(viper.GetString("logger.level"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger.Info("init") // @todo initialize app
}
