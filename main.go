package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"project/application"
)

func main() {
	app := application.New()

	//allow graceful shutdown upon signal to terminate (ctrl + C) - we want to finish database operations and so on BEFORE shut down
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	error := app.Start(ctx)
	if error != nil {
		fmt.Println("Failed to start app", error)
	}

}

// https://www.youtube.com/watch?v=qCv-q37qjZU
