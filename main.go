package main

import (
	"context"
	"fmt"
	"project/application"
)

func main() {
	app := application.New()

	error := app.Start(context.TODO())
	if error != nil {
		fmt.Println("Failed to start app", error)
	}
}
