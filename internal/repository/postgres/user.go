package postgres

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"shop/domain"
)

type Users struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *Users {
	return &Users{db: db}
}

func (r *Users) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	r.db.Where("username = ?", username).First(&user)
	if r.db.Error != nil {
		log.Errorf(r.db.Error.Error())
		return nil, r.db.Error
	}
	return &user, nil
}

func (r *Users) UpdateUser(user *domain.User) error {
	r.db.Save(&user)
	if r.db.Error != nil {
		log.Errorf(r.db.Error.Error())
		return r.db.Error
	}
	return nil
}
