package main

import (
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/helpers"
	"os"
)

func main() {

	addr := helpers.GoDotEnvVariable("APP_PORT")
	if addr == "" {
		logger.Logger.Log("Environment variable APP_PORT is not set.")
		os.Exit(1)
	}
	server, err := NewAPIServer(addr)
	if err != nil {
		logger.Logger.Error("Unable to start server: ", err)
		panic(err)
	}
	server.Run()
	// Prevents app from crashes
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Log("Recovered from panic:", r)
			// Handle the panic gracefully or exit the application as needed
		}
	}()
}
