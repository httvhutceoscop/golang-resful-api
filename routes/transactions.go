package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"kysuit.net/go-api/models"
)

func TransactionIndex(c *gin.Context) {
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	transaction, err := models.FetchAllTransaction(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": transaction,
	})
}

func TransactionCreate(c *gin.Context) {
	userId := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	newTransaction := models.Transaction{}

	// Call BindJSON to bind the received JSON to newTransaction
	if err := c.ShouldBindJSON(&newTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := newTransaction.Create(&conn, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, newTransaction)
}

func TransactionRead(c *gin.Context) {
	// userId := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		fmt.Println("id is invalid uuid")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	transaction, err := models.FindTransactionById(id, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

func TransactionUpdate(c *gin.Context) {
	// userId := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	transactionSent := models.Transaction{}
	err := c.ShouldBindJSON(&transactionSent)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid form sent",
		})
		return
	}

	transactionBeginUpdated, err := models.FindTransactionById(transactionSent.ID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Display uuid in response
	transactionSent.AccountId = transactionBeginUpdated.AccountId
	transactionSent.TransactionTypeId = transactionBeginUpdated.TransactionTypeId

	err = transactionSent.Update(&conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactionSent})
}

func TransactionDelete(c *gin.Context) {
	// userId := c.GetString("user_id")
	db, _ := c.Get("db")
	conn := db.(pgx.Conn)

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		fmt.Println("id is invalid uuid")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	transaction, err := models.FindTransactionById(id, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = transaction.Delete(transaction.ID, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delete transaction successfully"})
}

func TransactionSearch(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	searchType := c.Query("searchType")

	if len(from) == 0 {
		from = "1790-01-01"
	}
	from += " 00:00:00"

	if len(to) == 0 {
		to = time.Now().Format("2006-01-02")
	}
	to += " 23:59:59"

	db, _ := c.Get("db")
	conn := db.(pgx.Conn)
	response, err := models.SearchTransaction(from, to, searchType, &conn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
