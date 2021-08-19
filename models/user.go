package models

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	tokenSecret = []byte(os.Getenv("TOKEN_SECRET"))
)

type User struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	Password        string    `json:"password"`
	PasswordConfirm string    `json:"password_confirm"`
}

func (u *User) Register(conn *pgx.Conn) error {
	if len(u.Password) < 4 || len(u.PasswordConfirm) < 4 {
		return fmt.Errorf("password must be at least 4 characters long")
	}

	if u.Password != u.PasswordConfirm {
		return fmt.Errorf("password do not match")
	}

	if len(u.Email) < 4 {
		return fmt.Errorf("password must be at least 4 characters long")
	}

	u.Email = strings.ToLower(u.Email)
	row := conn.QueryRow(context.Background(), "SELECT id FROM auth_user WHERE email = $1", u.Email)
	userLookup := User{}
	err := row.Scan(&userLookup)
	if err != pgx.ErrNoRows {
		fmt.Println("found usr")
		fmt.Println(userLookup.Email)
		return fmt.Errorf("an user with that email already exists")
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("there was an error creating your account")
	}
	u.PasswordHash = string(pwdHash)

	now := time.Now()
	_, err = conn.Exec(context.Background(), "INSERT INTO auth_user (created_at, updated_at, email, password_hash) VALUES($1, $2, $3, $4)", now, now, u.Email, u.PasswordHash)

	return err
}

func (u *User) GetAuthToken() (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	authToken, err := token.SignedString(tokenSecret)
	return authToken, err
}

// IsAuthenticated checks to make sure password is correct and user is active
func (u *User) IsAuthenticated(conn *pgx.Conn) error {
	row := conn.QueryRow(context.Background(), "SELECT id, password_hash FROM auth_user WHERE email = $1", u.Email)
	err := row.Scan(&u.ID, &u.PasswordHash)
	if err == pgx.ErrNoRows {
		fmt.Println("user with email not found")
		return fmt.Errorf("invalid login credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return fmt.Errorf("invalid login credentials")
	}

	return nil
}

func IsTokenValid(tokenString string) (bool, string) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		fmt.Printf("Parsing: %v \n", t)
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("token signing method is not valid: %v", t.Header["alg"])
		}

		return tokenSecret, nil
	})

	if err != nil {
		fmt.Printf("Err %v \n", err)
		return false, ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		userId := claims["user_id"]
		return true, userId.(string)
	} else {
		fmt.Printf("The alg header %v \n", claims["alg"])
		fmt.Println(err)
		return false, "uuid.UUID{}"
	}
}
