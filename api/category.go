package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
)

type createCategorytRequest struct {
	Name string `json:"name" binding:"required,min=1"`
}

func (server *Server) createCategory(c *gin.Context) {
	var req createCategorytRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ticket, err := server.store.CreateCategory(c, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, ticket)

}

type listCategorytRequest struct {
	Page_id   int32 `form:"page_id" binding:"required,min=1"`
	Page_size int32 `form:"page_size" binding:"required,min=1"`
}

func (server *Server) listCategories(c *gin.Context) {
	var req listCategorytRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListCategoriesParams{
		Limit:  req.Page_size,
		Offset: req.Page_size * (req.Page_id - 1),
	}

	ticket, err := server.store.ListCategories(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, ticket)

}
