package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
	"time"
)

func getPost(id string) *Post {
	var post Post
	result := db.Take(&post, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return &post
	}
}

func editOrCreatePost(id string, title string, content string, privacy uint8) *Post {
	if id == "" {
		hash := sha256.New()
		hash.Write([]byte(id + title + time.Now().String()))
		md := hash.Sum(nil)
		id = hex.EncodeToString(md)
	}
	post := Post{
		ID:      id,
		Title:   title,
		Content: content,
		Date:    time.Now().UnixMilli(),
		Privacy: privacy,
	}
	_ = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "content", "privacy"}),
	}).Create(&post)
	return &post
}

func removePost(id string) bool {
	var post Post
	db.Where("id = ?", id).Delete(&post)
	return true
}

func getPosts(limit int, older int64) []Post {
	var posts []Post
	_ = db.Where("date < ?", older).Order("date desc").Limit(limit).Find(&posts)
	return posts
}

func setPostStat(postId string, userId string, action string) bool {
	a := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "post_id"}, {Name: "user_id"}, {Name: "action"}},
		DoNothing: true,
	}).Create(PostStat{
		PostId: postId,
		UserId: userId,
		Action: action,
		Date:   time.Now().UnixMilli(),
	}).RowsAffected
	if action == "like" && a == 0 {
		var postStat PostStat
		db.Where("post_id = ? and user_id = ? and action = ?", postId, userId, action).Delete(&postStat)
		return true
	} else {
		return a == 1
	}
}

func postChangeRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(email, "error")
		return c.SendString(res.String())
	}
	user := getProfile(email)
	if user == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}
	if !user.Admin {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	payload := struct {
		Id                string   `json:"id,omitempty"`
		Title             string   `json:"title,omitempty"`
		Content           string   `json:"content,omitempty"`
		Privacy           uint8    `json:"privacy,omitempty"`
		RemoveAttachments []string `json:"remove_attachments,omitempty"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		_, _ = res.Set("ERR_PARSE_BODY: "+err.Error(), "error")
		return c.SendString(res.String())
	}
	payload.Title = strings.TrimSpace(payload.Title)
	payload.Content = strings.TrimSpace(payload.Content)

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

	payload.Content = ugcPolicy.Sanitize(payload.Content)

	p := editOrCreatePost(payload.Id, payload.Title, payload.Content, payload.Privacy)
	_, _ = res.Set(true, "success")
	_, _ = res.Set(p.ID, "id")
	return c.SendString(res.String())
}

func postStatUpdateRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(email, "error")
		return c.SendString(res.String())
	}
	user := getProfile(email)
	if user == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}
	id := c.Get("id")
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
	_, _ = res.Set(setPostStat(id, email, action), "success")
	return c.SendString(res.String())
}

func postListRouteHandler(c *fiber.Ctx) error {
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	var user *User = nil
	if success {
		user = getProfile(email)
	}
	res := gabs.New()
	limit, err1 := strconv.Atoi(c.Query("limit", ""))
	if err1 != nil || limit > 50 {
		limit = 50
	}
	older, err2 := strconv.ParseInt(c.Query("older", ""), 10, 64)
	if err2 != nil {
		older = time.Now().UnixMilli()
	}
	_, _ = res.Array("posts")
	for _, post := range getPosts(limit, older) {
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

		result := struct {
			like int64
			view int64
		}{}
		x := db.Raw("select count(if(post_id = ? and action = 'like', 1, null)) as 'like', count(if(post_id = ? and action = 'view', 1, null)) as view from post_stats", post.ID)
		_ = x.Row().Scan(&result.like, &result.view)
		_, _ = p.Set(result.like, "likes")
		_, _ = p.Set(result.view, "views")
		_ = res.ArrayAppend(p, "posts")
	}
	return c.SendString(res.String())
}

func postRemoveRouteHandler(c *fiber.Ctx) error {
	res := gabs.New()
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	if !success {
		_, _ = res.Set(email, "error")
		return c.SendString(res.String())
	}
	user := getProfile(email)
	if user == nil {
		_, _ = res.Set("ERR_UNKNOWN_USER", "error")
		return c.SendString(res.String())
	}
	if !user.Admin {
		_, _ = res.Set("ERR_NO_PERMISSION", "error")
		return c.SendString(res.String())
	}

	id := c.Get("id")
	_, _ = res.Set(removePost(id), "success")
	return c.SendString(res.String())
}

func postGetRouteHandler(c *fiber.Ctx) error {
	token := c.Get("token")
	success, email := getEmailFromToken(token, c.UserContext())
	if !success {
		email = "***"
	}
	res := gabs.New()
	id := c.Query("id", "")
	if id == "" {
		_, _ = res.Set("ERR_INVALID_POST_ID", "error")
		return c.SendString(res.String())
	}
	post := getPost(id)
	if post == nil {
		_, _ = res.Set("ERR_UNKNOWN_POST", "error")
		return c.SendString(res.String())
	}

	var user *User = nil
	if success {
		user = getProfile(email)
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
		like      int64
		likeCheck int64
		view      int64
	}{}

	x := db.Raw("select count(if(post_id = ? and action = 'like', 1, null)) as 'like', count(if(post_id = ? and action = 'like' and user_id = ?, 1, null)) as likeCheck, count(if(post_id = ? and action = 'view', 1, null)) as view from post_stats", id, id, email, id)
	_ = x.Row().Scan(&result.like, &result.likeCheck, &result.view)
	_, _ = res.Set(result.like, "likes")
	_, _ = res.Set(result.view, "views")
	_, _ = res.Set(result.likeCheck > 0, "liked")

	return c.SendString(res.String())
}
