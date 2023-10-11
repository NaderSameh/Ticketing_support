package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
)

type createCommentRequestJSON struct {
	CommentText   string `json:"comment_text" binding:"required,min=1"`
	UserCommented string `json:"user_commented" binding:"required,min=1"`
}

type createCommentRequestURI struct {
	TicketID int64 `uri:"ticket_id" binding:"required,min=1"`
}

func (server *Server) createComment(c *gin.Context) {
	var reqJSON createCommentRequestJSON
	var reqURI createCommentRequestURI
	if err := c.ShouldBindJSON(&reqJSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateCommentParams{
		CommentText:   reqJSON.CommentText,
		TicketID:      reqURI.TicketID,
		UserCommented: reqJSON.UserCommented,
	}
	comment, err := server.store.CreateComment(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, comment)

}

type deleteCommentRequest struct {
	CommentID int64 `uri:"comment_id" binding:"required,min=0"`
}

func (server *Server) deleteComment(c *gin.Context) {
	var req deleteCommentRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if _, err := server.store.GetCommentForUpdate(c, req.CommentID); err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	err := server.store.DeleteComment(c, req.CommentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, true)

}

type updateCommentRequestURI struct {
	CommentID int64 `uri:"comment_id" binding:"required,min=0"`
	TicketID  int64 `uri:"ticket_id" binding:"required,min=0"`
}

type updateCommentRequestJSON struct {
	CommentText string `json:"comment_text" binding:"required,min=0"`
}

func (server *Server) updateComment(c *gin.Context) {
	var reqJSON updateCommentRequestJSON
	var reqURI updateCommentRequestURI

	if err := c.ShouldBindJSON(&reqJSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if _, err := server.store.GetTicketForUpdate(c, reqURI.TicketID); err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if _, err := server.store.GetCommentForUpdate(c, reqURI.CommentID); err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	arg := db.UpdateCommentParams{
		CommentID:   reqURI.CommentID,
		CommentText: reqJSON.CommentText,
	}

	arg2 := db.UpdateTicketParams{
		TicketID:  reqURI.TicketID,
		UpdatedAt: time.Now().Round(time.Second),
	}

	comment, err := server.store.UpdateComment(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.UpdateTicket(c, arg2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, comment)

}

type listCommentsRequestURI struct {
	TicketID int64 `uri:"ticket_id" binding:"required,min=1"`
}

type listCommentsRequestQuery struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listComments(c *gin.Context) {
	var reqQ listCommentsRequestQuery
	var reqURI listCommentsRequestURI
	if err := c.ShouldBindQuery(&reqQ); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if _, err := server.store.GetTicket(c, reqURI.TicketID); err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	arg := db.ListCommentsParams{
		TicketID: reqURI.TicketID,
		Limit:    reqQ.PageSize,
		Offset:   (reqQ.PageID - 1) * reqQ.PageSize,
	}

	comments, err := server.store.ListComments(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, comments)

}
