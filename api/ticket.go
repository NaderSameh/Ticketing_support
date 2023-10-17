package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
	"github.com/naderSameh/ticketing_support/token"
	"golang.org/x/exp/slices"
)

type createTicketRequest struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description" binding:"required"`
	Status       string `json:"status" binding:"required,oneof=inprogress closed open"`
	UserAssigned string `json:"user_assigned" binding:"required"`
	CategoryID   int64  `json:"category_id" binding:"required"`
}

// Createicket godoc
//
//	@Summary		Create ticket
//	@Description	Create a support ticket for an end user
//	@Tags			Tickets
//
//	@Accept			json
//	@Produce		json
//	@Param			arg	body		createTicketRequest	true	"Create Ticket body"
//
//	@Success		200	{object}	db.Ticket
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/tickets [post]
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

// DeleteTicket godoc
//
//	@Summary		Delete ticket
//	@Description	Delete ticket by a ticket ID
//	@Tags			Tickets
//
//
//	@Produce		plain
//	@Param			ticket_id	path		string	true	"Ticket ID"
//
//	@Success		200			true		bool
//	@Failure		400			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Router			/tickets/{ticket_id} [delete]
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

// ListTickets godoc
//
//	@Summary		List tickets
//	@Description	List all tickets for a specific user
//	@Tags			Tickets
//
//
//	@Produce		json
//	@Param			user_assigned	query		string	true	"Ticket owner"
//	@Param			page_id			query		int		true	"Page ID"
//	@Param			page_size		query		int		true	"Page Size"
//
//	@Success		200				{array}		db.Ticket
//	@Failure		400				{object}	error
//	@Failure		500				{object}	error
//	@Router			/tickets [get]
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

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if slices.Contains(authPayload.Permissions, "tickets.GET") {
		arg2 := db.ListAllTicketsParams{
			Limit:  req.PageSize,
			Offset: (req.PageID - 1) * req.PageSize,
		}
		tickets, err := server.store.ListAllTickets(c, arg2)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, tickets)
		return
	}

	tickets, err := server.store.ListTickets(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, tickets)

}

type updateTicketRequestJSON struct {
	Status      string `json:"status" binding:"required,oneof=inprogress closed open"`
	Assigned_to string `json:"assigned_to"`
}
type updateTicketRequestURI struct {
	TicketID int64 `uri:"ticket_id" binding:"required,min=1"`
}

// UpdateTicket godoc
//
//	@Summary		Update ticket
//	@Description	Update ticket by a ticket ID
//	@Tags			Tickets
//
//
//	@Produce		json
//
//	@Accept			json
//
//	@Param			arg			body		updateTicketRequestJSON	true	"Update ticket body"
//	@Param			ticket_id	path		int						true	"ticket ID for update"
//
//	@Success		200			{object}	db.Ticket
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Router			/tickets/{ticket_id} [put]
func (server *Server) updateTicket(c *gin.Context) {
	var req updateTicketRequestJSON
	var reqURI updateTicketRequestURI
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if !slices.Contains(authPayload.Permissions, "tickets.PUT") {
		err := errors.New("Only admins update tickets")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	arg := db.UpdateTicketParams{
		TicketID:  reqURI.TicketID,
		UpdatedAt: time.Now().Round(time.Second),
		Status:    req.Status,
	}
	if req.Assigned_to != "" {
		arg.AssignedTo = sql.NullString{String: req.Assigned_to, Valid: true}
	} else {
		arg.AssignedTo = sql.NullString{String: "", Valid: false}
	}

	_, err := server.store.GetTicketForUpdate(c, arg.TicketID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	tickets, err := server.store.UpdateTicket(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, tickets)

}
