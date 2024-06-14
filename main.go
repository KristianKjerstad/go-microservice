package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("hello")
	server := &http.Server{
		Addr:    ":3000",
		Handler: http.HandlerFunc(basicHandler),
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to listen to server")
	}
}

func basicHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("hello world"))

}
