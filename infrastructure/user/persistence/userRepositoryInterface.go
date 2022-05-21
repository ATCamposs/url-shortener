package persistence

import (
	"github.com/atcamposs/url-shortener/domain/user/entity"
)

type UserRepositoryInterface interface {
	Create(user entity.User) bool
	CheckEmailExists(email string) bool
	CheckUsernameExists(username string) bool
	RetrieveUserByEmail(email string) (entity.User, error)
}
