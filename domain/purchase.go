package domain

import "time"

type Purchase struct {
	GUID      string    `json:"guid" gorm:"column:guid;primaryKey"`
	User      User      `json:"-" gorm:"foreignKey:UserGUID;references:GUID"`
	UserGUID  string    `json:"user_guid" gorm:"column:user_guid;not null"`
	MerchGUID string    `json:"merch_guid" gorm:"column:merch_guid;not null"`
	Merch     Merch     `json:"-" gorm:"foreignKey:MerchGUID;references:GUID"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}
