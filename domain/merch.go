package domain

type Merch struct {
	Name  string  `json:"name" gorm:"column:name;not null;primaryKey"`
	Price float64 `json:"price" gorm:"column:price;type:decimal(20,8);not null"`
}
