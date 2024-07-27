package controllers

import (
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"go-svc/config"
	"go-svc/src/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Define database client
var db *gorm.DB = config.ConnectDB()

// request object
type postRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// response object
type postResponse struct {
	postRequest
	ID uint `json:"id"`
}

// Getting all post datas
func GetAllPosts(context *gin.Context) {
	var posts []models.Post

	// check posts is exist
	errPost := db.Find(&posts).Where("is_active = true")
	if errPost.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error getting post data"})
		return
	}

	// send response
	context.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "Success",
		"data":    constructDataResponse(posts),
	})
}

// Getting post by id
func GetPostById(context *gin.Context) {
	var post models.Post
	var posts []models.Post

	// check post is exist
	errPost := db.Where("is_active = true").Where("id = ?", context.Param("postId")).Find(&post)
	if errPost.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error getting post data"})
		return
	}
	posts = append(posts, post)

	// send response
	context.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "Success",
		"data":    constructDataResponse(posts),
	})
}

// Create post
func CreatePost(context *gin.Context) {

	var request postRequest
	var posts []models.Post
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post = models.Post{
		Title:     request.Title,
		Content:   request.Content,
		IsActive:  true,
		CreatedBy: getFunctionName(CreatePost),
		CreatedAt: time.Now(),
	}
	db.Create(&post)

	if request.Tags != nil {
		var postTags []models.PostsTag

		for _, element := range request.Tags {
			var postTag models.PostsTag
			postTag.PostId = post.ID

			// check tag is exist
			var tag models.Tag
			errTag := db.Where("is_active = true").Where("label = ?", element).Find(&tag)
			if errTag.Error != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": "Error getting tag data"})
				return
			}
			if tag.ID == 0 {
				tag.Label = element
				tag.IsActive = true
				tag.CreatedAt = time.Now()
				tag.CreatedBy = post.CreatedBy
				db.Create(&tag)
			}
			postTag.TagId = postTag.ID
			postTag.TagId = tag.ID
			postTag.IsActive = true
			postTag.CreatedAt = time.Now()
			postTag.CreatedBy = post.CreatedBy
			postTags = append(postTags, postTag)
		}
		db.CreateInBatches(postTags, len(postTags))
	}
	posts = append(posts, post)

	// send response
	context.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "Success",
		"data":    constructDataResponse(posts),
	})
}

// Update post
func UpdatePost(context *gin.Context) {

	var request postRequest
	var post models.Post
	var posts []models.Post

	time := time.Now()
	updatedBy := getFunctionName(UpdatePost)

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check post is exist
	errPost := db.Where("id = ?", context.Param("postId")).Find(&post)
	if errPost.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error getting post data"})
		return
	}
	post.Content = request.Content
	post.Title = request.Title
	post.IsActive = true
	post.LastUpdatedAt = &time
	post.LastUpdatedBy = &updatedBy
	db.Save(&post)

	if request.Tags != nil {
		var postTags []models.PostsTag

		for _, element := range request.Tags {
			var postTag models.PostsTag
			postTag.PostId = post.ID

			// check tag is exist
			var tag models.Tag
			errTag := db.Where("is_active = true").Where("label = ?", element).Find(&tag)
			if errTag.Error != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": "Error getting tag data"})
				return
			}

			if tag.ID == 0 {
				tag.Label = element
				tag.IsActive = true
				tag.CreatedAt = time
				tag.CreatedBy = post.CreatedBy
				db.Create(&tag)

				postTag.TagId = postTag.ID
				postTag.TagId = tag.ID
				postTag.IsActive = true
				postTag.CreatedAt = time
				postTag.CreatedBy = post.CreatedBy
				postTags = append(postTags, postTag)
				continue
			}

			// check post tag is exist
			errPostTag := db.Where("is_active = true").Where("post_id = ?", post.ID).Where("tag_id = ?", tag.ID).Find(&postTag)
			if errPostTag.Error != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": "Error getting post tag data"})
				return
			}

			if postTag.ID != 0 {
				continue
			}
			postTag.TagId = postTag.ID
			postTag.TagId = tag.ID
			postTag.IsActive = true
			postTag.CreatedAt = time
			postTag.CreatedBy = post.CreatedBy
			postTags = append(postTags, postTag)

		}
		db.CreateInBatches(postTags, len(postTags))
	}

	// delete post tag
	var postTags []models.PostsTag
	errpostTags := db.Where("is_active = true").Where("post_id = ?", post.ID).Find(&postTags)
	if errpostTags.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error getting postTag data"})
		return
	}
	postTagStr := getTagLabelByPostId(post.ID)

	for idx, postTag := range postTagStr {
		isDelete := true
		for _, tag := range request.Tags {
			if postTag == tag {
				isDelete = false
				continue
			}
		}
		if isDelete {
			postTags[idx].IsActive = false
			postTags[idx].LastUpdatedAt = &time
			postTags[idx].LastUpdatedBy = &post.CreatedBy
			db.Save(postTags[idx])
		}
	}

	posts = append(posts, post)
	// Creating http response
	context.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "Success",
		"data":    constructDataResponse(posts),
	})
}

// Delete post
func DeletePost(context *gin.Context) {
	var post models.Post

	time := time.Now()
	updatedBy := getFunctionName(DeletePost)

	// check post is exist
	errPost := db.Where("is_active = true").Where("id = ?", context.Param("postId")).Find(&post)
	if errPost.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error getting delete data"})
		return
	}
	post.IsActive = false
	post.LastUpdatedAt = &time
	post.LastUpdatedBy = &updatedBy

	db.Save(&post)

	// send response
	context.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "Delete Success for id :" + strconv.Itoa(int(post.ID)),
	})
}

// private function
func constructDataResponse(posts []models.Post) []postResponse {
	var result []postResponse

	if posts == nil {
		return result
	}

	for _, element := range posts {
		var data postResponse
		data.ID = element.ID
		data.Title = element.Title
		data.Content = element.Content
		data.Tags = getTagLabelByPostId(element.ID)
		result = append(result, data)
	}

	return result
}

// private function
func getTagLabelByPostId(post_id uint) []string {
	var tags []string

	query := "SELECT t2.label FROM posts_tags t1 JOIN tags t2 ON t1.tag_id = t2.id WHERE t1.is_active = true And t1.post_id = ?"

	// Raw SQL
	db.Raw(query, post_id).Scan(&tags)
	return tags
}

// private function
func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
