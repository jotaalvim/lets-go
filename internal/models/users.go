package models

import (
    "errors"
    "strings"
	"database/sql"
	"time"

    "golang.org/x/crypto/bcrypt"
    "github.com/go-sql-driver/mysql"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name string, email string, password string) error {

    hashed_password , err := bcrypt.GenerateFromPassword( []byte (password), 12)
    if err != nil { return err }


	stmt := `INSERT INTO users (name, email, hashed_password, created)
             VALUES (?,?,?, UTC_TIMESTAMP())`


    _, err = m.DB.Exec(stmt, name, email, string(hashed_password))


    if err != nil { 
        var mySQLError *mysql.MySQLError
        // check if the error has the mysql type
        // if it does the error will be asign to that variable
        if errors.As( err, &mySQLError) {
            if mySQLError.Number == 1062 && strings.Contains( mySQLError.Message, "users_uc_email") {
                return ErrDuplicateEmail
            }
        }
        return err
    }

    return nil

}

func (m *UserModel) Authenticate(email string, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
