package main

import (
	"context"
	"log"
	"os"

	"github.com/jtbry/camvera/internal/vision"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs/")
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Println("Config File Not Found")
		} else {
			logger.Fatal("Config file found but error reading it")
		}
	}

	cv := vision.NewCvWorker(ctx, logger)

	go cv.Start()

	<-ctx.Done()
}
