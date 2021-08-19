package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

type Account struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	Name           string    `json:"name"`
	InitialBalance int64     `json:"initial_balance"`
	Description    string    `json:"description"`
	CurrencyId     uuid.UUID `json:"currency_id"`
	UserAccountId  uuid.UUID `json:"auth_user_id"`
}

func (a *Account) Create(conn *pgx.Conn, userId string) error {
	a.Name = strings.Trim(a.Name, " ")
	if len(a.Name) < 1 {
		return fmt.Errorf("name must not be empty")
	}

	if a.InitialBalance < 0 {
		a.InitialBalance = 0
	}
	now := time.Now()

	// row := conn.QueryRow(context.Background(), "INSERT INTO account (name, initial_balance, description, currency_id, auth_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", a.Name, a.InitialBalance, a.Description, currency_id, userId, now, now)
	// err := row.Scan(&a.ID)

	_, err := conn.Exec(context.Background(), "INSERT INTO account (name, initial_balance, description, currency_id, auth_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", a.Name, a.InitialBalance, a.Description, a.CurrencyId, userId, now, now)

	if err != nil {
		// fmt.Println(err)
		return fmt.Errorf(err.Error())
	}

	return nil
}

func (a *Account) Update(conn *pgx.Conn) error {
	a.Name = strings.Trim(a.Name, " ")
	if len(a.Name) < 1 {
		return fmt.Errorf("name must not be empty")
	}
	if a.InitialBalance < 0 {
		a.InitialBalance = 0
	}
	now := time.Now()
	_, err := conn.Exec(context.Background(), "UPDATE account SET name=$1, initial_balance=$2, description=$3, updated_at=$4 WHERE id=$5", a.Name, a.InitialBalance, a.Description, now, a.ID)

	if err != nil {
		fmt.Printf("error updating account: (%v)", err)
		return fmt.Errorf("error updating account")
	}

	return nil
}

func FetchAllAccounts(conn *pgx.Conn) ([]Account, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, name, initial_balance, description, auth_user_id, currency_id FROM account")
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error getting accounts")
	}

	var accounts []Account
	for rows.Next() {
		account := Account{}
		err = rows.Scan(&account.ID, &account.Name, &account.InitialBalance, &account.Description, &account.UserAccountId, &account.CurrencyId)
		if err != nil {
			fmt.Print(err)
			continue
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func GetAccountsByUser(userId string, conn *pgx.Conn) ([]Account, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, name, initial_balance, description, auth_user_id, currency_id FROM account WHERE auth_user_id = $1", userId)
	if err != nil {
		fmt.Printf("error getting accounts %v", err)
		return nil, fmt.Errorf("there was an error getting the accounts")
	}
	var accounts []Account
	for rows.Next() {
		a := Account{}
		err = rows.Scan(&a.ID, &a.Name, &a.InitialBalance, &a.Description, &a.UserAccountId, &a.CurrencyId)
		if err != nil {
			fmt.Printf("error scaning accounts: %v", err)
			continue
		}
		accounts = append(accounts, a)
	}

	return accounts, nil
}

func FindAccountById(id uuid.UUID, conn *pgx.Conn) (Account, error) {
	row := conn.QueryRow(context.Background(), "SELECT id, name, initial_balance, description, auth_user_id, currency_id FROM account WHERE id = $1", id)
	account := Account{
		ID: id,
	}
	err := row.Scan(&account.ID, &account.Name, &account.InitialBalance, &account.Description, &account.UserAccountId, &account.CurrencyId)
	if err != nil {
		fmt.Printf("the account doesn't exist: (%v)", err)
		return account, fmt.Errorf("the account doesn't exist")
	}

	return account, nil
}
