package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"kysuit.net/go-api/models"
)

func AccountsIndex(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	accounts, err := models.FetchAllAccounts(&conn)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}

func AccountsCreate(c *gin.Context) {
	userId := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	account := models.Account{}
	c.ShouldBindJSON(&account)
	err := account.Create(&conn, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func AccountsByCurrentUser(c *gin.Context) {
	userId := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	accounts, err := models.GetAccountsByUser(userId, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

func AccountsUpdate(c *gin.Context) {
	userId := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	accountSent := models.Account{}
	err := c.ShouldBindJSON(&accountSent)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid form sent",
		})
		return
	}

	accountBeginUpdated, err := models.FindAccountById(accountSent.ID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if accountBeginUpdated.UserAccountId.String() != userId {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to update this account",
		})
		return
	}

	// Display uuid in response
	accountSent.UserAccountId = accountBeginUpdated.UserAccountId
	accountSent.CurrencyId = accountBeginUpdated.CurrencyId

	err = accountSent.Update(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": accountSent})
}
