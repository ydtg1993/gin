package model

type Label struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
	Sort uint   `gorm:"not null" json:"sort"`
}

func (Label) TableName() string {
	return "label"
}
