package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
)

type createTicketRequest struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description" binding:"required"`
	Status       string `json:"status" binding:"required,oneof=inprogress closed open"`
	UserAssigned string `json:"user_assigned" binding:"required"`
	CategoryID   int64  `json:"category_id" binding:"required"`
}

func (server *Server) createTicket(c *gin.Context) {
	var req createTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTicketParams{
		Title:        req.Title,
		Description:  req.Description,
		Status:       req.Status,
		UserAssigned: req.UserAssigned,
		CategoryID:   req.CategoryID,
	}
	ticket, err := server.store.CreateTicket(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, ticket)

}

type deleteTicketRequest struct {
	TicketID int64 `uri:"ticket_id" binding:"required,min=1"`
}

func (server *Server) deleteTicket(c *gin.Context) {
	var req deleteTicketRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.store.GetTicket(c, req.TicketID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err = server.store.DeleteTicket(c, req.TicketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusNoContent, true)

}

type listTicketRequest struct {
	UserAssigned string `form:"user_assigned"`
	PageID       int32  `form:"page_id" binding:"required,min=1"`
	PageSize     int32  `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listTicket(c *gin.Context) {
	var req listTicketRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListTicketsParams{
		UserAssigned: req.UserAssigned,
		Limit:        req.PageSize,
		Offset:       (req.PageID - 1) * req.PageSize,
	}

	tickets, err := server.store.ListTickets(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, tickets)

}

type updateTicketRequest struct {
	Ticket_ID   int64  `json:"ticket_id" binding:"required,min=1"`
	Status      string `json:"status" binding:"required,oneof=inprogress closed open"`
	Assigned_to string `json:"assigned_to"`
}

func (server *Server) updateTicket(c *gin.Context) {
	var req updateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateTicketParams{
		TicketID:  req.Ticket_ID,
		UpdatedAt: time.Now().Round(time.Second),
		Status:    req.Status,
	}
	if req.Assigned_to != "" {
		arg.AssignedTo = sql.NullString{String: req.Assigned_to, Valid: true}
	} else {
		arg.AssignedTo = sql.NullString{String: "", Valid: false}
	}

	tickets, err := server.store.UpdateTicket(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, tickets)

}
