package domain

import "time"

type Purchase struct {
	GUID      string    `json:"guid" gorm:"column:guid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"column:user_id;not null;index:idx_user_merch"`
	User      User      `json:"-" gorm:"foreignKey:UserID;references:Username"`
	MerchName string    `json:"merch_name" gorm:"column:merch_name;not null;index:idx_user_merch"`
	Merch     Merch     `json:"-" gorm:"foreignKey:MerchName;references:name"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}
