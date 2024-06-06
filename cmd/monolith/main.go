package main

import (
	"fmt"
	"log/slog"
	"staffinc/internal/handler"
)

func main() {

	server := handler.NewServer()

	slog.Info("Starting Staffinc")

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
