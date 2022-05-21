package persistence

import (
	"errors"
	"fmt"
	"log"

	"github.com/atcamposs/url-shortener/domain/user/entity"
	"github.com/atcamposs/url-shortener/infrastructure"
)

type PostgresRepository struct{}

func New() UserRepositoryInterface {
	return &PostgresRepository{}
}

func (r *PostgresRepository) Create(u entity.User) bool {
	stmt, err := infrastructure.DefaultConnection.Prepare(`INSERT INTO users (uuid,email,username,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)`)
	checkErrors(err)

	result, err := stmt.Exec(u.UUID.String(), u.Email, u.Username, u.Password, u.CreatedAt, u.UpdatedAt)
	checkErrors(err)

	affect, err := result.RowsAffected()
	checkErrors(err)

	fmt.Println(affect)
	return true
}

func (r *PostgresRepository) CheckEmailExists(email string) bool {
	stmt, err := infrastructure.DefaultConnection.Prepare(`SELECT COUNT(*) FROM users WHERE email=$1`)
	checkErrors(err)

	log.Printf(`SELECT COUNT(*) FROM users WHERE email='%s'`, email)

	count := 0
	err = stmt.QueryRow(email).Scan(&count)
	checkErrors(err)

	return count > 0
}

func (r *PostgresRepository) CheckUsernameExists(username string) bool {
	stmt, err := infrastructure.DefaultConnection.Prepare(`SELECT COUNT(*) FROM users WHERE username=$1`)
	checkErrors(err)

	log.Printf(`SELECT COUNT(*) FROM users WHERE username='%s'`, username)

	count := 0
	err = stmt.QueryRow(username).Scan(&count)
	checkErrors(err)

	return count > 0
}

func (r *PostgresRepository) RetrieveUserByEmail(email string) (entity.User, error) {
	stmt, err := infrastructure.DefaultConnection.Prepare(`SELECT * FROM users WHERE email=$1`)
	checkErrors(err)

	log.Printf(`SELECT * FROM users WHERE email='%s'`, email)

	var user entity.User
	err = stmt.QueryRow(email).Scan(&user.UUID, &user.Email, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err == nil {
		return user, nil
	}

	return user, errors.New("no rows return")
}

func checkErrors(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}
