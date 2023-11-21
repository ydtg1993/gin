package model

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"xo/core"
)

type Video struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SourceID  uint      `gorm:"not null" json:"source_id"`
	SourceURL string    `gorm:"not null;uniqueIndex;index:source_url_idx" json:"source_url"`
	Title     string    `gorm:"not null" json:"title"`
	Cover     string    `gorm:"not null" json:"cover"`
	BigCover  string    `gorm:"not null" json:"big_cover"`
	URL       string    `gorm:"not null" json:"url"`
	Like      int       `gorm:"not null;default:0" json:"like"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

func (Video) TableName() string {
	return "video"
}

func GetFormattedLabelList(vid, lid uint) string {
	var labelList strings.Builder

	var videoLabelAss []VLAss
	if vid > 0 {
		core.Mysql.Select("label_id").Where("video_id = ?", vid).Find(&videoLabelAss)
	} else {
		core.Mysql.Select("label_id").Find(&videoLabelAss)
	}

	var labelIds []uint
	for _, label := range videoLabelAss {
		labelIds = append(labelIds, label.LabelId)
	}
	var labels []Label
	core.Mysql.Where("id IN ?", labelIds).Order("sort desc").Find(&labels)

	for _, label := range labels {
		if label.ID == lid {
			labelList.WriteString(fmt.Sprintf(`<li><a class="green" href="%stag/%d/index.html">%s</a></li>`, core.Config.GetString("app.host"), label.ID, label.Name))
			continue
		}
		labelList.WriteString(fmt.Sprintf(`<li><a href="%stag/%d/index.html">%s</a></li>`, core.Config.GetString("app.host"), label.ID, label.Name))
	}
	return labelList.String()
}

func GetRandomVideos(limit int) (videos []Video) {
	var totalRecords int64
	core.Mysql.Model(&Video{}).Count(&totalRecords)

	// Generate a random offset
	rand.Seed(time.Now().UnixNano())
	randomOffset := rand.Intn(int(totalRecords - int64(limit) + 1))

	core.Mysql.Offset(randomOffset).Limit(limit).Find(&videos)
	return videos
}
