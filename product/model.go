package product

type Product struct {
	ID          string  `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string  `json:"name" gorm:"type:varchar(128);not null"`
	Description string  `json:"description" gorm:"type:text"`
	Price       float64 `json:"price" gorm:"type:decimal(10,2);not null"`
}
