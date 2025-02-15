package domain

type Merch struct {
	GUID  string  `json:"guid" gorm:"column:guid;primaryKey"`
	Name  string  `json:"name" gorm:"column:name;not null;unique"`
	Price float64 `json:"price" gorm:"column:price;type:decimal(20,8);not null"`
}
