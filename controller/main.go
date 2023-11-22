package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"math"
	"net/http"
	"regexp"
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
	labels := model.GetFormattedLabelList(0, 0)
	data := gin.H{
		"vlist":      template.HTML(video_temp.String()),
		"label_list": template.HTML(labels),
	}
	c.HTML(http.StatusOK, "home.html", data)
}

func Tag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.Query("page"))
	sort := c.Query("sort")
	limit := 2
	var vlass []model.VLAss
	if page < 1 {
		page = 1
	}
	if id < 0 {
		c.Redirect(http.StatusNotFound, "/404.html")
		return
	}
	var Tag model.Label
	result := core.Mysql.Where("id = ?", id).First(&Tag)
	if result.Error != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/404.html")
		return
	}
	query := core.Mysql.
		Preload("Video").
		Joins("JOIN video ON video.id = video_label_ass.video_id").
		Where("video_label_ass.label_id = ?", id)
	if sort == "time" {
		query = query.Order("video.created_at desc")
	} else if sort == "like" {
		query = query.Order("video.like desc")
	} else if sort == "title" {
		query = query.Order("video.title desc")
	} else {
		query = query.Order("video.like desc").Order("video.created_at desc")
	}

	result = query.
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&vlass)
	if result.Error != nil || len(vlass) == 0 {
		c.Redirect(http.StatusTemporaryRedirect, "/404.html")
		return
	}
	var count int64
	core.Mysql.Model(model.VLAss{}).Where("label_id = ?", id).Count(&count)
	page_count := math.Ceil(float64(count) / float64(limit))
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
	labels := model.GetFormattedLabelList(0, uint(id))
	data := gin.H{
		"Title":      Tag.Name,
		"label_list": template.HTML(labels),
		"vlist":      template.HTML(video_temp.String()),
		"url":        core.Config.GetString("app.host") + "tag/" + strconv.Itoa(id) + "/index.html?page=:page",
		"page":       page,
		"page_count": page_count,
	}
	c.HTML(http.StatusOK, "tag.html", data)
}

func Search(c *gin.Context) {
	keywords := c.Query("keywords")
	re := regexp.MustCompile(`[!@#$%^&*()_+{}:“<>?,.\/;'\[\]\\|` + "`" + `~"'【】，。、，]+`)
	keywords = re.ReplaceAllString(keywords, "")
	if len(keywords) > 64 {
		keywords = keywords[:64]
	}

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

	labels := model.GetFormattedLabelList(video.ID, 0)

	data := gin.H{
		"Title":      video.Title,
		"url":        "https://xgct-video.vzcdn.net/4244a3d1-227f-467c-a5d9-d4209ea7e270/1280x720/video.m3u8",
		"vlist":      template.HTML(video_temp.String()),
		"label_list": template.HTML(labels),
	}
	c.HTML(http.StatusOK, "video.html", data)
}
