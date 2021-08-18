package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"kysuit.net/go-api/database"
	"kysuit.net/go-api/middlewares"
	"kysuit.net/go-api/routes"
)

func setupRouter(conn pgx.Conn) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.DbMiddleWare(conn))

	usersGroup := r.Group("/api/v1")
	{
		usersGroup.POST("/user/create", routes.UsersRegister)
		usersGroup.POST("/user/login", routes.UsersLogin)
	}

	accountsGroup := r.Group("/api/v1/account")
	{
		accountsGroup.GET("/", middlewares.AuthMiddleWare(), routes.AccountsIndex)
		accountsGroup.GET("/by_current_user", middlewares.AuthMiddleWare(), routes.AccountsByCurrentUser)
		accountsGroup.POST("/", middlewares.AuthMiddleWare(), routes.AccountsCreate)
		accountsGroup.PUT("/", middlewares.AuthMiddleWare(), routes.AccountsUpdate)
	}

	transactionGroup := r.Group("/api/v1/transaction")
	{
		transactionGroup.GET("/", middlewares.AuthMiddleWare(), routes.TransactionIndex)
		transactionGroup.POST("/", middlewares.AuthMiddleWare(), routes.TransactionCreate)
		transactionGroup.GET("/:id", middlewares.AuthMiddleWare(), routes.TransactionRead)
		transactionGroup.PUT("/", middlewares.AuthMiddleWare(), routes.TransactionUpdate)
		transactionGroup.DELETE("/:id", middlewares.AuthMiddleWare(), routes.TransactionDelete)
		transactionGroup.GET("/search", middlewares.AuthMiddleWare(), routes.TransactionSearch)
	}

	return r
}

func main() {

	conn, err := database.ConnectDB()

	if err != nil {
		fmt.Printf("Make db connection error: %v", err)
		return
	}

	r := setupRouter(*conn)
	r.Run(":5000")
}
