package model

type VLAss struct {
	ID      uint `gorm:"primaryKey" json:"id"`
	LabelId uint `gorm:"not null" json:"label_id"`
	VideoId uint `gorm:"not null" json:"video_id"`
	Video   Video
}

func (VLAss) TableName() string {
	return "video_label_ass"
}
