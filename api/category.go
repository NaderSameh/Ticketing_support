package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naderSameh/ticketing_support/cache"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
	"github.com/rs/zerolog/log"
)

type createCategorytRequest struct {
	Name string `json:"name" binding:"required,min=1"`
}

// CreateCategory godoc
//
//	@Summary		Create new category
//	@Description	Create a new category specifying its name
//	@Tags			Categories
//	@Produce		json
//	@Accept			json
//	@Param			arg	body		createCategorytRequest	true	"Create category body"
//
//	@Success		200	{object}	db.Category
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/categories [post]
func (server *Server) createCategory(c *gin.Context) {
	var req createCategorytRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	category, err := server.store.CreateCategory(c, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, category)

}

type listCategorytRequest struct {
	Page_id   int32 `form:"page_id" binding:"required,min=1"`
	Page_size int32 `form:"page_size" binding:"required,min=1"`
}

// ListCategories godoc
//
//	@Summary		Get all categories
//	@Description	get all categories names, pagination options available
//	@Tags			Categories
//	@Produce		json
//
//	@Param			page_id		query		int	true	"Page ID"
//	@Param			page_size	query		int	true	"Page Size"
//
//	@Success		200			{array}		db.Category
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Router			/categories [get]
func (server *Server) listCategories(c *gin.Context) {
	var req listCategorytRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	url := c.FullPath() + c.Query("page_id") + c.Query("page_size")
	redisClient := cache.NewCacheClient().RedisClient

	arg := db.ListCategoriesParams{
		Limit:  req.Page_size,
		Offset: req.Page_size * (req.Page_id - 1),
	}
	ctx := context.Background()
	res, err := redisClient.Get(ctx, url).Result()

	if err == nil {
		log.Info().Str("Source", "cache").Msg("Cached resource")
		resStruct := []db.Category{}
		_ = json.Unmarshal([]byte(res), &resStruct)
		c.JSON(http.StatusOK, resStruct)
		return
	}
	categories, err := server.store.ListCategories(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// put it to the cache
	content, _ := json.Marshal(categories)
	err = redisClient.Set(c, url, content, time.Minute*5).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, categories)

}
