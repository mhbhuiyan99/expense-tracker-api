package models

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func getUserFilePath() string {
	dir := beego.AppConfig.DefaultString("datadir", "data")
	file := beego.AppConfig.DefaultString("userfile", "users.csv")
	return dir + "/" + file
}

// GetAllUsers reads all users from the CSV file
func GetAllUsers() ([]User, error) {
	path := getUserFilePath()

	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return []User{}, nil // no file means no users
		}
		return []User{}, err
	}
	defer f.Close()

	row, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	if len(row) <= 1 {
		return []User{}, nil // only header or empty file
	}

	var users []User
	for _, row := range row[1:] { // skip header
		user, err := rowToUse(row)
		if err != nil {
			logs.Warn("Skipping invalid user row: %v", err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUserByEmail finds a user by email
func GetUserByEmail(email string) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, nil
}

// GetUserByID finds a user by ID
func GetUserByID(id int) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, nil
}

// CreateUser adds a new user to the CSV file
func CreateUser(user *User) error {
	path := getUserFilePath()

	// ensure data/ directory exists
	dir := beego.AppConfig.DefaultString("datadir", "data")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// write header if file is new
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		if err := writer.Write([]string{"ID", "Name", "Email", "Password", "CreatedAt"}); err != nil {
			return err
		}
	}

	user.ID, err = GetNextUserID()
	if err != nil {
		return err
	}

	user.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := writer.Write(userToRow(user)); err != nil {
		return err
	}

	return writer.Error()
}

// GetNextUserID returns the next available user ID
func GetNextUserID() (int, error) {
	users, err := GetAllUsers()
	if err != nil {
		return 0, err
	}

	maxID := 0
	for _, u := range users {
		if u.ID > maxID {
			maxID = u.ID
		}
	}

	return maxID + 1, nil
}

// rowToUser converts a CSV row to a User struct
func rowToUse(row []string) (User, error) {
	if len(row) != 5 {
		return User{}, fmt.Errorf("invalid row length: %d", len(row))
	}

	id, err := strconv.Atoi(row[0])
	if err != nil {
		return User{}, fmt.Errorf("invalid user ID %q: %w", row[0], err)
	}

	return User{
		ID:        id,
		Name:      row[1],
		Email:     row[2],
		Password:  row[3],
		CreatedAt: row[4],
	}, nil
}

// userToRow converts a User struct to a CSV row
func userToRow(user *User) []string {
	return []string{
		strconv.Itoa(user.ID),
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
	}
}
