package models

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UsersModel interface {
	Register(db *sql.DB, inputRegister InputRegister) error
	Login(db *sql.DB, inputLogin InputLogin) (*User, error)
}

type InputRegister struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(db *sql.DB, inputRegister InputRegister) error {
	sqlStatement := `INSERT INTO users (email, password) VALUES (?, ?)`
	result, err := db.Exec(sqlStatement, inputRegister.Email, inputRegister.Password)
	if err != nil {
		log.Fatalf("could not insert row: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("could not get affected rows: %v", err)
		return err
	}
	// we can log how many rows were inserted
	fmt.Println("inserted", rowsAffected, "rows")
	return nil
}

type InputLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *sql.DB, inputLogin InputLogin) (*User, error) {
	var user User
	sqlStatement := `SELECT id,email from USERS where email = ? AND password = ?`
	err := db.QueryRow(sqlStatement, inputLogin.Email, inputLogin.Password).Scan(&user.ID, &user.Email)
	if err != nil {
		log.Fatalf("error, cause: %v", err)
		return nil, err
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user not found")
	}
	fmt.Printf("user found")
	return &user, nil
}

type InputCreateToken struct {
	UserId   string `json:"user_id"`
	Token    string `json:"token"`
	IsActive bool   `json:"is_active"`
}

func CreateToken(db *sql.DB, data InputCreateToken) error {
	sqlStatement := `INSERT INTO users_token (user_id, token, is_active) VALUES (?, ?, ?)`
	result, err := db.Exec(sqlStatement, data.UserId, data.Token, data.IsActive)
	if err != nil {
		log.Fatalf("could not insert row: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("could not get affected rows: %v", err)
		return err
	}
	// we can log how many rows were inserted
	fmt.Println("inserted", rowsAffected, "rows")
	return nil
}

type InputGetToken struct {
	Token string `json:"token"`
}

func GetToken(db *sql.DB, data InputGetToken) (*User, error) {
	var user User
	fmt.Println(data.Token)
	sqlStatement := `SELECT email,password from USERS LEFT JOIN users_token ON users_token.user_id = users.id where token = ? and users_token.is_active = true`
	err := db.QueryRow(sqlStatement, data.Token).Scan(&user.Email, &user.Password)
	if err != nil {
		log.Fatalf("error, cause: %v", err)
		return nil, err
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user not found")
	}
	fmt.Printf("user found")
	return &user, nil
}
