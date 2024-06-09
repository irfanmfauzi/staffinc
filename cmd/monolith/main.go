package main

import (
	"fmt"
	"log/slog"
	"staffinc/internal/handler"
)

// @title          Swagger Example API
// @version        1.0
// @description    This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     petstore.swagger.io
// @BasePath /api
func main() {

	server := handler.NewServer()

	slog.Info("Starting Staffinc")

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
