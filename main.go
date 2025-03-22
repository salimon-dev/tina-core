package main

import (
	"fmt"
	"os"
	"salimon/nexus/auth"
	"salimon/nexus/db"
	"salimon/nexus/e2e"
	"salimon/nexus/entities"
	"salimon/nexus/entity"
	"salimon/nexus/invitations"
	"salimon/nexus/middlewares"
	"salimon/nexus/profile"
	"salimon/nexus/rest"
	"salimon/nexus/users"
	"salimon/nexus/websocket"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("no environment file, using session defaults")
	}
	db.SetupDatabase()
	e := echo.New()
	e.HideBanner = true
	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.DELETE, echo.PUT},
	}))

	// HTTP route
	e.GET("/", rest.HeartBeatHandler)

	// register
	e.POST("/auth/register", auth.RegisterHandler)

	// login
	e.POST("/auth/login", auth.LoginHandler)
	e.POST("/auth/rotate", auth.RotateHandler)
	e.POST("/auth/entity-token", auth.EntityTokenHandler, middlewares.UserAuthMiddleware)

	// user info
	e.GET("/profile", profile.GetHandler, middlewares.UserAuthMiddleware)

	// WebSocket route
	e.GET("/sck", websocket.WsHandler)

	// -- -- Admin APIs -- --
	// invitations
	e.GET("/invitations/search", invitations.SearchHandler, middlewares.UserAuthMiddleware, middlewares.AdminMiddleware)
	e.POST("/invitations/create", invitations.CreateHandler, middlewares.UserAuthMiddleware, middlewares.AdminMiddleware)
	e.POST("/invitations/delete/:id", invitations.DeleteHandler, middlewares.UserAuthMiddleware, middlewares.AdminMiddleware)
	e.POST("/invitations/update/:id", invitations.UpdateHandler, middlewares.UserAuthMiddleware, middlewares.AdminMiddleware)
	// users
	e.GET("/users/search", users.SearchHandler, middlewares.UserAuthMiddleware, middlewares.AdminMiddleware)
	e.POST("/users/create", users.CreateHandler, middlewares.UserAuthMiddleware, middlewares.AdminMiddleware)
	e.POST("/users/delete/:id", users.DeleteHandler, middlewares.UserAuthMiddleware, middlewares.AdminMiddleware)
	e.POST("/users/update/:id", users.UpdateHandler, middlewares.UserAuthMiddleware, middlewares.AdminMiddleware)
	// entities
	e.POST("/entities/create", entities.CreateHandler, middlewares.UserAuthMiddleware, middlewares.AdminMiddleware)
	e.POST("/entities/delete/:id", entities.DeleteHandler, middlewares.UserAuthMiddleware, middlewares.AdminMiddleware)
	e.POST("/entities/update/:id", entities.UpdateHandler, middlewares.UserAuthMiddleware, middlewares.AdminMiddleware)
	e.POST("/entities/token", entities.EntityTokenHandler, middlewares.UserAuthMiddleware, middlewares.AdminMiddleware)
	// search is hybrid endpoint for all permission
	e.GET("/entities/search", entities.SearchHandler, middlewares.UserAuthMiddleware)

	// -- -- Entity APIs -- --
	e.GET("/entity/user/:userId", entity.GetUserHandler, middlewares.EntityAuthMiddleware)

	// -- -- External APIs -- --
	// E2E control Endpoints
	e.POST("/e2e/interact", e2e.InteractHandler)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
