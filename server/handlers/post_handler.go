package handlers

import (
	"das/models"
	"das/models/request"
	"das/security"
	"das/storage"
	"das/utils"
	"errors"
	"github.com/Jeffail/gabs/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"regexp"
	"strconv"
	"strings"
)

const MaxPostTitleLength = 300
const MinPostTitleLength = 10
const MaxPostContentLength = 100000
const MinPostContentLength = 100
const MaxPostHeadlineLength = 250
const MinPostHeadlineLength = 30
const MaxPostHashtagLength = 20
const MinPostHashtagLength = 5
const PostListLimit = 50

func getPost(id interface{}) *models.Post {
	var post models.Post
	result := storage.Db.Take(&post, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else if result.Error != nil {
		log.Error().Err(result.Error).Msg("An error occurred at #getPost while processing DB transaction")
		return nil
	} else {
		return &post
	}
}

func serializePost(post *models.Post, requester *models.User, withContent bool) *gabs.Container {
	response := post.Serialize(withContent)
	if requester != nil && security.IsLoggedIn(requester.Role) {
		result := struct {
			likeCheck int64
			viewCheck int64
		}{}

		x := storage.Db.Raw("select count(if(post_id = ? and action = 'like' and user_id = ?, 1, null)) as likeCheck, count(if(post_id = ? and action = 'view' and user_id = ?, 1, null)) as likeCheck from post_stats", post.ID, requester.ID, post.ID, requester.ID)
		_ = x.Row().Scan(&result.likeCheck, &result.viewCheck)
		_, _ = response.Set(result.likeCheck, "stats", "liked")
		_, _ = response.Set(result.viewCheck, "stats", "viewed")
	}
	_, _ = response.Array("attachments")
	for _, att := range getAttachments(post.ID) {
		_ = response.ArrayAppend(att.Serialize(), "attachments")
	}
	return response
}

func updateOrCreatePost(id uint32, req *request.PostUpdateModel) *models.Post {
	post := models.Post{
		Link:     utils.GenerateLinkFromTitle(req.Title),
		Title:    req.Title,
		Content:  req.Content,
		Headline: req.Headline,
		Hashtag:  req.Hashtag,
		Privacy:  req.Privacy,
	}
	if id > 0 {
		post.ID = id
	}
	tx := storage.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "link", "content", "privacy", "headline", "hashtag"}),
	}).Create(&post)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #updateOrCreatePost while processing DB transaction")
	}
	return &post
}

func removePost(id interface{}) bool {
	var post models.Post
	tx := storage.Db.Where("id = ?", id).Delete(&post)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #removePost while processing DB transaction")
		return false
	}
	return tx.RowsAffected > 0
}

func getPosts(req *request.PostListModel, requester *models.User) []models.Post {
	var posts []models.Post
	cmd := storage.Db.Limit(int(req.Limit))
	cmd = cmd.Where("privacy <= ?", requester.Role)
	if req.BelowId > 0 {
		cmd = cmd.Where("id < ?", req.BelowId)
	}
	if len(req.FilterHashtags) > 0 {
		cmd = cmd.Where("hashtag IN ?", req.FilterHashtags)
	}
	if req.SortBy == "view" {
		cmd = cmd.Order("view_count DESC, id DESC")
		if req.LowerThan > 0 {
			cmd = cmd.Where("view_count < ?", req.LowerThan)
		}
	} else if req.SortBy == "like" {
		cmd = cmd.Order("like_count DESC, id DESC")
		if req.LowerThan > 0 {
			cmd = cmd.Where("like_count < ?", req.LowerThan)
		}
	} else if req.SortBy == "date" {
		cmd = cmd.Order("update_date DESC, id DESC")
		if req.LowerThan > 0 {
			cmd = cmd.Where("update_date < ?", req.LowerThan)
		}
	} else {
		cmd = cmd.Order("id DESC")
	}
	cmd = cmd.Find(&posts)
	if cmd.Error != nil {
		log.Error().Err(cmd.Error).Msg("An error occurred at #getPosts while processing DB transaction")
	}
	return posts
}

func getPostStat(postId uint32, action string) int64 {
	result := struct {
		v int64
	}{}
	err := storage.Db.Model(&models.Post{}).Where("id = ?", postId).Select(action + "_count").Row().Scan(&result.v)
	if err != nil {
		log.Error().Err(err).Msg("An error occurred at #getPostStat while processing DB transaction")
		return -1
	}
	return result.v
}

func deletePostStat(postId uint32, userId uint16, action string) bool {
	tx := storage.Db.Where("post_id = ? and user_id = ? and action = ?", postId, userId, action).Delete(&models.PostStat{})
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #deletePostStat while processing DB transaction (1)")
		return false
	}
	if tx.RowsAffected == 0 {
		return true
	}
	tx = storage.Db.Model(&models.Post{}).Where("id = ?", postId).UpdateColumn(action+"_count", gorm.Expr(action+"_count - "+strconv.FormatInt(tx.RowsAffected, 10)))
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #deletePostStat while processing DB transaction (2)")
		return false
	}
	return tx.RowsAffected > 0
}

// setPostStat 0=error; 1=existed; 2=increased
func setPostStat(postId uint32, userId uint16, action string) uint {
	tx := storage.Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "post_id"}, {Name: "user_id"}, {Name: "action"}},
		DoNothing: true,
	}).Create(models.PostStat{
		PostId: postId,
		UserId: userId,
		Action: action,
	})
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #setPostStat while processing DB transaction")
		return 0
	}
	if tx.RowsAffected == 0 {
		return 1
	}
	tx = storage.Db.Model(&models.Post{}).Where("id = ?", postId).UpdateColumn(action+"_count", gorm.Expr(action+"_count + 1"))
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #setPostStat while processing DB transaction (2)")
		return 0
	}
	if tx.RowsAffected > 0 {
		return 2
	}
	return 0
}

func PostGetRouteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" || !utils.ValidateNonNegativeInteger(id) {
		return ReturnError(c, utils.ErrUnknownPost)
	}
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}

	post := getPost(id)
	if post == nil {
		return ReturnError(c, utils.ErrUnknownPost)
	}
	if post.Privacy > requester.Role {
		return ReturnError(c, utils.ErrNoPermission)
	}
	return ReturnJSON(c, serializePost(post, requester, true))
}

func PostUpdateRouteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id != "" && !utils.ValidateNonNegativeInteger(id) {
		return ReturnError(c, utils.ErrUnknownPost)
	}
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}
	if security.GetRoleGroup(requester.Role) < security.RoleGroupGlobalManager {
		return ReturnError(c, utils.ErrNoPermission)
	}
	postId := uint32(0)
	if id != "" {
		post := getPost(id)
		if post == nil {
			return ReturnError(c, utils.ErrUnknownPost)
		}
		if post.Privacy > requester.Role {
			return ReturnError(c, utils.ErrNoPermission)
		}
		postId = post.ID
	}

	req := &request.PostUpdateModel{}
	if err2 := c.BodyParser(&req); err2 != nil {
		log.Error().Err(err2).Msg("There was an error occurred while parsing body at #PostUpdateRouteHandler")
		return ReturnError(c, utils.ErrInvalidRequestBody)
	}
	req.Title = security.SafeHTMLPolicy.Sanitize(strings.TrimSpace(req.Title))
	req.Headline = security.SafeHTMLPolicy.Sanitize(strings.TrimSpace(req.Headline))
	req.Content = security.SafeHTMLPolicy.Sanitize(strings.TrimSpace(req.Content))

	if len(req.Title) < MinPostTitleLength {
		return ReturnError(c, utils.ErrPostTitleTooShort)
	} else if len(req.Title) > MaxPostTitleLength {
		return ReturnError(c, utils.ErrPostTitleTooLong)
	}

	if req.Headline == "" {
		req.Headline = utils.LimitString(bluemonday.StrictPolicy().Sanitize(req.Content), MaxPostHeadlineLength)
	} else if len(req.Headline) < MinPostHeadlineLength {
		return ReturnError(c, utils.ErrPostHeadlineTooShort)
	} else if len(req.Headline) > MaxPostHeadlineLength {
		return ReturnError(c, utils.ErrPostHeadlineTooLong)
	}

	if len(req.Content) < MinPostContentLength {
		return ReturnError(c, utils.ErrPostContentTooShort)
	} else if len(req.Content) > MaxPostContentLength {
		return ReturnError(c, utils.ErrPostContentTooLong)
	}

	req.Hashtag = security.SafeHTMLPolicy.Sanitize(strings.TrimSpace(req.Hashtag))
	if len(req.Hashtag) < MinPostHashtagLength {
		return ReturnError(c, utils.ErrPostHashtagTooShort)
	} else if len(req.Hashtag) > MaxPostHashtagLength {
		return ReturnError(c, utils.ErrPostHashtagTooLong)
	}
	matched, err3 := regexp.MatchString("^[a-zA-Z\\d-_]+$", req.Hashtag)
	if !matched || err3 != nil {
		return ReturnError(c, utils.ErrInvalidPostHashtag)
	}

	post := updateOrCreatePost(postId, req)
	if post == nil {
		if postId == 0 {
			return ReturnError(c, utils.ErrPostCreateFailed)
		} else {
			return ReturnError(c, utils.ErrPostUpdateFailed)
		}
	}
	response := gabs.New()
	_, _ = response.Set(post.ID, "id")
	return ReturnJSON(c, response)
}

func PostDeleteRouteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" || !utils.ValidateNonNegativeInteger(id) {
		return ReturnError(c, utils.ErrUnknownPost)
	}
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}
	if security.GetRoleGroup(requester.Role) < security.RoleGroupGlobalManager {
		return ReturnError(c, utils.ErrNoPermission)
	}

	post := getPost(id)
	if post == nil {
		return ReturnError(c, utils.ErrUnknownPost)
	}
	if post.Privacy > requester.Role {
		return ReturnError(c, utils.ErrNoPermission)
	}
	if removePost(post.ID) {
		return ReturnEmpty(c)
	} else {
		return ReturnError(c, utils.ErrPostDeleteFailed)
	}
}

func PostListRouteHandler(c *fiber.Ctx) error {
	req := request.PostListModel{}
	if err := c.QueryParser(&req); err != nil {
		log.Error().Err(err).Msg("There was an error occurred while parsing body at #PostListRouteHandler")
		return ReturnError(c, utils.ErrInvalidRequestQuery)
	}
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}

	req.Limit = utils.ClampUint8(req.Limit, 0, PostListLimit)
	hashtagFilter := make([]string, 0)
	for _, v := range req.FilterHashtags {
		v = strings.TrimSpace(v)
		matched, err := regexp.MatchString("^[a-zA-Z\\d-_]+$", v)
		if matched && err == nil {
			hashtagFilter = append(hashtagFilter, v)
		}
	}
	req.FilterHashtags = hashtagFilter

	posts := gabs.New()
	_, _ = posts.Array("posts")
	for _, post := range getPosts(&req, requester) {
		_ = posts.ArrayAppend(serializePost(&post, requester, false), "posts")
	}
	return ReturnJSON(c, posts)
}

func PostStatUpdateRouteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" || !utils.ValidateNonNegativeInteger(id) {
		return ReturnError(c, utils.ErrUnknownPost)
	}
	json, err2 := gabs.ParseJSON(c.Body())
	if err2 != nil {
		return ReturnError(c, utils.ErrInvalidRequestBody)
	}
	requester, err := GetRequester(c)
	if err != "" {
		return ReturnError(c, err)
	}
	if !security.IsLoggedIn(requester.Role) {
		return ReturnError(c, utils.ErrNoPermission)
	}

	post := getPost(id)
	if post == nil {
		return ReturnError(c, utils.ErrUnknownPost)
	}
	if post.Privacy > requester.Role {
		return ReturnError(c, utils.ErrNoPermission)
	}
	response := gabs.New()
	if json.Exists("view") && json.Path("view").Data().(bool) {
		q := setPostStat(post.ID, requester.ID, "view")
		if q == 0 {
			return ReturnError(c, utils.ErrPostStatUpdateFailed)
		}
		v := getPostStat(post.ID, "view")
		if v == -1 {
			return ReturnError(c, utils.ErrPostStatUpdateFailed)
		}
		_, _ = response.Set(v, "views")
	}
	if json.Exists("like") && json.Path("like").Data().(bool) {
		q := setPostStat(post.ID, requester.ID, "like")
		if q == 0 {
			return ReturnError(c, utils.ErrPostStatUpdateFailed)
		} else if q == 1 { // existed: like: true -> false
			ok := deletePostStat(post.ID, requester.ID, "like")
			if !ok {
				return ReturnError(c, utils.ErrPostStatUpdateFailed)
			}
		}
		v := getPostStat(post.ID, "like")
		if v == -1 {
			return ReturnError(c, utils.ErrPostStatUpdateFailed)
		}
		_, _ = response.Set(v, "likes")
	}
	return ReturnJSON(c, response)
}

func PostHashtagListRouteHandler(c *fiber.Ctx) error {
	var hashtags []struct {
		Hashtag string
	}
	tx := storage.Db.Model(&models.Post{}).Distinct("hashtag").Find(&hashtags)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred at #PostHashtagListRouteHandler while processing DB transaction")
		return ReturnError(c, utils.ErrPostHashtagListFailed)
	}
	res := gabs.New()
	_, _ = res.Array("hashtags")
	for _, t := range hashtags {
		_ = res.ArrayAppend(t.Hashtag, "hashtags")
	}
	return ReturnJSON(c, res)
}
