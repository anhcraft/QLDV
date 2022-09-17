package handlers

import (
	"das"
	"das/models"
	"das/utils"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func maxString(s string, max int) string {
	limit := len(s)
	if limit > max {
		limit = max
	}
	return s[:limit]
}

func editOrCreatePost(id int, title string, content string, privacy uint8, hashtag string) *models.Post {
	post := models.Post{
		Link:     utils.GenerateLinkFromTitle(title),
		Title:    title,
		Content:  content,
		Headline: maxString(bluemonday.StrictPolicy().Sanitize(content), 250),
		Hashtag:  hashtag,
		Date:     time.Now().UnixMilli(),
		Privacy:  privacy,
	}
	if id > 0 {
		post.ID = id
	}
	_ = main.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "link", "content", "privacy", "headline", "hashtag"}),
	}).Create(&post)
	return &post
}

func getPosts(filterHashtags []string, sortBy string, lowerThan uint, belowId int, limit int) []models.Post {
	var posts []models.Post
	cmd := main.db.Limit(limit)
	if len(filterHashtags) > 0 {
		cmd = cmd.Where("hashtag IN ?", filterHashtags)
	}
	if belowId > 0 {
		cmd = cmd.Where("id < ?", belowId)
	}
	if sortBy == "view" {
		cmd = cmd.Order("view_count DESC, id DESC")
		if lowerThan > 0 {
			cmd = cmd.Where("view_count < ?", lowerThan)
		}
	} else if sortBy == "like" {
		cmd = cmd.Order("like_count DESC, id DESC")
		if lowerThan > 0 {
			cmd = cmd.Where("like_count < ?", lowerThan)
		}
	} else {
		cmd = cmd.Order("id DESC")
	}
	_ = cmd.Find(&posts)
	return posts
}

func setPostStat(postId int, userId string, action string) bool {
	a := main.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "post_id"}, {Name: "user_id"}, {Name: "action"}},
		DoNothing: true,
	}).Create(models.PostStat{
		PostId: postId,
		UserId: userId,
		Action: action,
		Date:   time.Now().UnixMilli(),
	}).RowsAffected
	if action == "like" && a == 0 {
		var postStat models.PostStat
		main.db.Where("post_id = ? and user_id = ? and action = ?", postId, userId, action).Delete(&postStat)
		main.db.Model(&models.Post{}).Where("id = ?", postId).UpdateColumn("like_count", gorm.Expr("like_count - 1"))
		return true
	} else if a == 1 {
		if action == "like" {
			main.db.Model(&models.Post{}).Where("id = ?", postId).UpdateColumn("like_count", gorm.Expr("like_count + 1"))
			return true
		} else if action == "view" {
			main.db.Model(&models.Post{}).Where("id = ?", postId).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
			return true
		}
	}
	return false
}

func PostUpdateRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, emailOrError := main.getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(emailOrError, "error")
		return c.SendString(res.String())
	}
	user := getUserByEmail(emailOrError)
	if user == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}
	if !user.Admin {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	payload := struct {
		Id                int      `json:"id,omitempty"`
		Title             string   `json:"title,omitempty"`
		Content           string   `json:"content,omitempty"`
		Privacy           uint8    `json:"privacy,omitempty"`
		Hashtag           string   `json:"hashtag,omitempty"`
		RemoveAttachments []string `json:"remove_attachments,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}
	payload.Title = strings.TrimSpace(payload.Title)
	payload.Content = strings.TrimSpace(payload.Content)
	payload.Hashtag = strings.TrimSpace(payload.Hashtag)

	if len(payload.Title) < 5 {
		_, _ = res.Set("ERR_POST_TITLE_MIN", "error")
		return c.SendString(res.String())
	} else if len(payload.Title) > 300 {
		_, _ = res.Set("ERR_POST_TITLE_MAX", "error")
		return c.SendString(res.String())
	}

	if len(payload.Content) < 10 {
		_, _ = res.Set("ERR_POST_CONTENT_MIN", "error")
		return c.SendString(res.String())
	} else if len(payload.Content) > 100000 {
		_, _ = res.Set("ERR_POST_CONTENT_MAX", "error")
		return c.SendString(res.String())
	}

	if len(payload.RemoveAttachments) > 0 {
		if !removeAttachments(payload.Id, payload.RemoveAttachments) {
			_, _ = res.Set("ERR_ATTACHMENT_REMOVE_FAILED", "error")
			return c.SendString(res.String())
		}
	}

	matched, err := regexp.MatchString("^[a-zA-Z\\d-_]+$", payload.Hashtag)
	if !matched || err != nil || len(payload.Hashtag) < 5 || len(payload.Hashtag) > 30 {
		_, _ = res.Set("ERR_INVALID_HASHTAG", "error")
		return c.SendString(res.String())
	}

	payload.Title = main.ugcPolicy.Sanitize(payload.Title)
	payload.Content = main.ugcPolicy.Sanitize(payload.Content)

	p := editOrCreatePost(payload.Id, payload.Title, payload.Content, payload.Privacy, payload.Hashtag)
	_, _ = res.Set(true, "success")
	_, _ = res.Set(p.ID, "id")
	return c.SendString(res.String())
}

func PostStatUpdateRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, emailOrError := main.getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(emailOrError, "error")
		return c.SendString(res.String())
	}
	user := getUserByEmail(emailOrError)
	if user == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}

	id, err := strconv.Atoi(c.Get("id"))
	if err != nil {
		_, _ = res.Set("ERR_INVALID_POST_ID", "error")
		return c.SendString(res.String())
	}
	post := getPost(id)
	if post == nil {
		_, _ = res.Set("ERR_UNKNOWN_POST", "error")
		return c.SendString(res.String())
	}
	if (post.Privacy&2) == 2 && !(user.Mod || user.Admin) {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}
	if (post.Privacy&4) == 4 && !user.Admin {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	action := c.Get("action")
	if !(action == "like" || action == "view") {
		_, _ = res.Set("ERR_UNKNOWN_POST_ACTION", "error")
		return c.SendString(res.String())
	}
	_, _ = res.Set(setPostStat(id, emailOrError, action), "success")
	return c.SendString(res.String())
}

func PostListRouteHandler(c *fiber.Ctx) error {
	token := c.Get("token")
	success, emailOrError := main.getEmailFromToken(token, c.UserContext())
	var user *models.User = nil
	if success {
		user = getUserByEmail(emailOrError)
	}

	res := gabs.New()
	payload := struct {
		Limit          int      `json:"limit,omitempty"`
		FilterHashtags []string `json:"filter_hashtags,omitempty"`
		BelowId        int      `json:"below_id,omitempty"`
		SortBy         string   `json:"sort_by,omitempty"`
		LowerThan      uint     `json:"lower_than,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}

	if payload.Limit > 50 {
		payload.Limit = 50
	} else if payload.Limit < 1 {
		payload.Limit = 1
	}

	hashtagFilter := make([]string, 0)
	for i := range payload.FilterHashtags {
		v := strings.TrimSpace(payload.FilterHashtags[i])
		matched, err := regexp.MatchString("^[a-zA-Z\\d-_]+$", v)
		if matched && err == nil {
			hashtagFilter = append(hashtagFilter, v)
		}
	}

	_, _ = res.Array("posts")
	for _, post := range getPosts(hashtagFilter, payload.SortBy, payload.LowerThan, payload.BelowId, payload.Limit) {
		if (post.Privacy&1) == 1 && user == nil {
			continue
		}
		if (post.Privacy&2) == 2 && (user == nil || !(user.Mod || user.Admin)) {
			continue
		}
		if (post.Privacy&4) == 4 && (user == nil || !user.Admin) {
			continue
		}

		p := post.serialize()
		_, _ = p.Array("attachments")
		for _, att := range getAttachments(post.ID) {
			_ = p.ArrayAppend(att.serialize(), "attachments")
		}
		_ = res.ArrayAppend(p, "posts")
	}
	return c.SendString(res.String())
}

func PostRemoveRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, emailOrError := main.getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(emailOrError, "error")
		return c.SendString(res.String())
	}
	user := getUserByEmail(emailOrError)
	if user == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}
	if !user.Admin {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	id, err := strconv.Atoi(c.Get("id"))
	if err != nil {
		_, _ = res.Set("ERR_INVALID_POST_ID", "error")
		return c.SendString(res.String())
	}
	_, _ = res.Set(removePost(id), "success")
	return c.SendString(res.String())
}

func PostGetRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, emailOrError := main.getEmailFromToken(token, c.UserContext())
	if !success {
		// guest can view public posts
		emailOrError = "***"
	}
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		_, _ = res.Set("ERR_INVALID_POST_ID", "error")
		return c.SendString(res.String())
	}
	post := getPost(id)
	if post == nil {
		_, _ = res.Set("ERR_UNKNOWN_POST", "error")
		return c.SendString(res.String())
	}

	var user *models.User = nil
	if success {
		user = getUserByEmail(emailOrError)
	}
	if (post.Privacy&1) == 1 && user == nil {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}
	if (post.Privacy&2) == 2 && (user == nil || !(user.Mod || user.Admin)) {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}
	if (post.Privacy&4) == 4 && (user == nil || !user.Admin) {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	res = post.serialize()
	_, _ = res.Set(post.Content, "content")
	_, _ = res.Array("attachments")
	for _, att := range getAttachments(id) {
		_ = res.ArrayAppend(att.serialize(), "attachments")
	}

	result := struct {
		likeCheck int64
	}{}

	x := main.db.Raw("count(if(post_id = ? and action = 'like' and user_id = ?, 1, null)) as likeCheck", id, emailOrError)
	_ = x.Row().Scan(&result.likeCheck)
	_, _ = res.Set(result.likeCheck > 0, "liked")

	return c.SendString(res.String())
}

func PostHashtagListRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	var hashtags []struct {
		Hashtag string
	}
	main.db.Model(&models.Post{}).Distinct("hashtag").Find(&hashtags)
	_, _ = res.Array("hashtags")
	for _, t := range hashtags {
		_ = res.ArrayAppend(t.Hashtag, "hashtags")
	}
	return c.SendString(res.String())
}
