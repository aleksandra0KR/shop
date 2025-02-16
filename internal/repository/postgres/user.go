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

func (r *Users) UpdateUser(tx *gorm.DB, user *domain.User) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}
	db.Save(&user)
	if db.Error != nil {
		log.Errorf(db.Error.Error())
		return db.Error
	}
	return nil
}
func (r *Users) CreateUser(user *domain.User) (*domain.User, error) {
	r.db.Create(&user)
	if r.db.Error != nil {
		log.Errorf(r.db.Error.Error())
		return nil, r.db.Error
	}
	return user, nil
}
