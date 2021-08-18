package models

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

type Transaction struct {
	ID                uuid.UUID `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Amount            int64     `json:"amount"`
	Description       string    `json:"description"`
	AccountId         uuid.UUID `json:"account_id"`
	TransactionTypeId uuid.UUID `json:"transaction_type_id"`
}

func (t *Transaction) Create(conn *pgx.Conn, userId string) error {
	if t.Amount <= 0 {
		return fmt.Errorf("amount must be larger than 0")
	}
	now := time.Now()

	// TODO: need to validate if account's amount less than or equal zero

	_, err := conn.Exec(context.Background(), "INSERT INTO transaction (amount, description, account_id, transaction_type_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", t.Amount, t.Description, t.AccountId, t.TransactionTypeId, now, now)

	if err != nil {
		fmt.Println("Error in transaction.Create()")
		return fmt.Errorf(err.Error())
	}

	return nil
}

func (t *Transaction) Update(conn *pgx.Conn) error {
	if t.Amount <= 0 {
		return fmt.Errorf("amount must be larger than 0")
	}
	now := time.Now()

	// TODO: need to validate if account's amount less than or equal zero

	_, err := conn.Exec(context.Background(), "UPDATE transaction SET amount=$1, description=$2, account_id=$3, transaction_type_id=$4, updated_at=$5 WHERE id=$6", t.Amount, t.Description, t.AccountId, t.TransactionTypeId, now, t.ID)

	if err != nil {
		fmt.Printf("error updating transaction: (%v)", err)
		return fmt.Errorf("error updating transaction")
	}

	return nil
}

func FetchAllTransaction(conn *pgx.Conn) ([]Transaction, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, amount, description, account_id, transaction_type_id, created_at, updated_at FROM transaction")
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error getting transaction")
	}

	var arrData []Transaction
	for rows.Next() {
		obj := Transaction{}
		err = rows.Scan(&obj.ID, &obj.Amount, &obj.Description, &obj.AccountId, &obj.TransactionTypeId, &obj.CreatedAt, &obj.UpdatedAt)
		if err != nil {
			fmt.Print(err)
			continue
		}
		arrData = append(arrData, obj)
	}

	return arrData, nil
}

func FindTransactionById(id uuid.UUID, conn *pgx.Conn) (Transaction, error) {
	row := conn.QueryRow(context.Background(), "SELECT id, amount, description, account_id, transaction_type_id, created_at, updated_at FROM transaction WHERE id = $1", id)
	obj := Transaction{
		ID: id,
	}
	err := row.Scan(&obj.ID, &obj.Amount, &obj.Description, &obj.AccountId, &obj.TransactionTypeId, &obj.CreatedAt, &obj.UpdatedAt)
	if err != nil {
		fmt.Printf("the transaction doesn't exist: (%v)", err)
		return obj, fmt.Errorf("the transaction doesn't exist")
	}

	return obj, nil
}

func (t *Transaction) Delete(id uuid.UUID, conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM transaction WHERE id=$1", t.ID)

	if err != nil {
		fmt.Printf("error deleting transaction: (%v)", err)
		return fmt.Errorf("error deleting transaction")
	}

	return nil
}

func SearchTransaction(from string, to string, searchType string, conn *pgx.Conn) (map[string]interface{}, error) {

	// searchType: D - Income, W - expenses, B - Balance
	var sql string
	var rows pgx.Rows
	var err error

	if searchType == "B" {
		sql = "SELECT tbl1.ID, tbl1.amount, tbl1.description, tbl1.account_id, tbl1.transaction_type_id, tbl1.created_at, tbl2.code FROM transaction tbl1 LEFT JOIN transaction_type tbl2 ON tbl1.transaction_type_id = tbl2.id WHERE (tbl1.created_at BETWEEN $1 AND $2)"
		rows, err = conn.Query(context.Background(), sql, from, to)
	} else {
		sql = "SELECT tbl1.ID, tbl1.amount, tbl1.description, tbl1.account_id, tbl1.transaction_type_id, tbl1.created_at, tbl2.code FROM transaction tbl1 LEFT JOIN transaction_type tbl2 ON tbl1.transaction_type_id = tbl2.id WHERE tbl2.code = $1 AND (tbl1.created_at BETWEEN $2 AND $3)"
		rows, err = conn.Query(context.Background(), sql, searchType, from, to)
	}

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error searching transaction")
	}

	var totalIncome int64
	var totalExpense int64
	var arrData []Transaction
	var total int64
	var transactionTypeCode string
	for rows.Next() {
		obj := Transaction{}
		err = rows.Scan(&obj.ID, &obj.Amount, &obj.Description, &obj.AccountId, &obj.TransactionTypeId, &obj.CreatedAt, &transactionTypeCode)
		if err != nil {
			fmt.Print(err)
			continue
		}
		arrData = append(arrData, obj)

		// total += obj.Amount
		if transactionTypeCode == "D" {
			totalIncome += obj.Amount
		} else if transactionTypeCode == "W" {
			totalExpense += obj.Amount
		}
	}

	if searchType == "D" {
		total = totalIncome
	} else if searchType == "W" {
		total = totalExpense
	} else if searchType == "B" {
		total = totalIncome - totalExpense
	}

	result := map[string]interface{}{
		"transaction_type": searchType,
		"totalAmount":      total,
		"transaction":      arrData,
	}

	return result, nil
}
