package main

import (
	"fmt"
	"staffinc/internal/handler"
)

func main() {

	server := handler.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
