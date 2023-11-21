package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"xo/core"
	"xo/model"
)

func Home(c *gin.Context) {
	var videos []model.Video
	result := core.Mysql.Order("created_at desc").Limit(32).Find(&videos)
	if result.Error != nil {
		c.Redirect(http.StatusNotFound, "/404.html")
		return
	}
	var video_temp strings.Builder
	for _, video := range videos {
		temp := fmt.Sprintf(`
<article>
<a href="/video/%d/detail.html" class="vbox">
	<div class="cover"><img decoding="async" src="%s"/></div>
	<header class="trim_text">
		<span>%s</span>
	</header>
</a></article>
`, video.ID, "https://sta.anicdn.com/playerImg/8.jpg", video.Title)
		video_temp.WriteString(temp)
	}
	labels := model.GetFormattedLabelList(0)
	data := gin.H{
		"Title":      "Hello, Gin!",
		"vlist":      template.HTML(video_temp.String()),
		"label_list": template.HTML(labels),
	}
	c.HTML(http.StatusOK, "home.html", data)
}

func Tag(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.Query("p"))
	limit := 24
	var vlass []model.VLAss
	result := core.Mysql.
		Preload("Video").
		Joins("JOIN video ON video.id = video_label_ass.video_id").
		Where("video_label_ass.label_id = ?", id).
		Offset(page * limit).
		Limit(limit).
		Find(&vlass)
	if result.Error != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/404.html")
		return
	}
	var video_temp strings.Builder
	for _, vl := range vlass {
		temp := fmt.Sprintf(`
<article>
<a href="/video/%d/detail.html" class="vbox">
	<div class="cover"><img decoding="async" src="%s"/></div>
	<header class="trim_text">
		<span>%s</span>
	</header>
</a></article>
`, vl.Video.ID, "https://sta.anicdn.com/playerImg/8.jpg", vl.Video.Title)
		video_temp.WriteString(temp)
	}

	data := gin.H{
		"Title": "- tag",
		"vlist": template.HTML(video_temp.String()),
	}
	c.HTML(http.StatusOK, "tag.html", data)
}

func Search(c *gin.Context) {
	keywords := c.Query("keywords")
	var videos []model.Video
	result := core.Mysql.Where("title LIKE ?", "%"+keywords+"%").Limit(24).Find(&videos)
	if result.Error != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/404.html")
	}
	var video_temp strings.Builder
	result_total := len(videos)
	if result_total == 0 {
		videos = model.GetRandomVideos(4)
	}
	for _, video := range videos {
		temp := fmt.Sprintf(`
<article>
<a href="/video/%d/detail.html" class="vbox">
	<div class="cover"><img decoding="async" src="%s"/></div>
	<header class="trim_text">
		<span>%s</span>
	</header>
</a></article>
`, video.ID, "https://sta.anicdn.com/playerImg/8.jpg", video.Title)
		video_temp.WriteString(temp)
	}

	data := gin.H{
		"Title":    "- Search",
		"keywords": keywords,
		"total":    result_total,
		"vlist":    template.HTML(video_temp.String()),
	}
	c.HTML(http.StatusOK, "search.html", data)
}

func Video(c *gin.Context) {
	id := c.Param("id")
	var video model.Video
	result := core.Mysql.Where("id = ?", id).First(&video)
	if result.Error != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/404.html")
		return
	}
	var videos []model.Video
	result = core.Mysql.Order("created_at desc").Limit(12).Find(&videos)
	var video_temp strings.Builder
	if result.Error == nil {
		for _, video := range videos {
			temp := fmt.Sprintf(`
<article>
<a href="/video/%d/detail.html" class="vbox">
	<div class="cover"><img decoding="async" src="%s"/></div>
	<header class="trim_text">
		<span>%s</span>
	</header>
</a></article>
`, video.ID, "https://sta.anicdn.com/playerImg/8.jpg", video.Title)
			video_temp.WriteString(temp)
		}
	}

	labels := model.GetFormattedLabelList(video.ID)

	data := gin.H{
		"Title":      "Hello, Gin!",
		"url":        "https://xgct-video.vzcdn.net/4244a3d1-227f-467c-a5d9-d4209ea7e270/1280x720/video.m3u8",
		"vlist":      template.HTML(video_temp.String()),
		"label_list": template.HTML(labels),
	}
	c.HTML(http.StatusOK, "video.html", data)
}
