package main

import (
	"os"
	"sync"
	"timeTracker/config"
	api "timeTracker/internal"
	"timeTracker/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// @title Time Tracker
// @description Api Endpoints for time tracker
func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"err": err,
		}).Error("Can't load config from .env. Problem with .env, or the server is in production environment.")
		return
	}

	config := config.ApiEnvConfig{
		Port:        os.Getenv("PORT"),
		Env:         os.Getenv("ENV"),
		Host:        os.Getenv("HOST"),
		AuthService: os.Getenv("AUTH_URL"),
	}

	logger.Log.WithFields(logrus.Fields{
		"port": config.Port,
		"host": config.Host,
	}).Info("Loaded app config")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		server := api.AppServer{}
		defer func() {
			if r := recover(); r != nil {
				server.OnShutdown()
			}
		}()

		server.Run(config)
	}()
	wg.Wait()

}

// cSpell:ignore logrus godotenv
