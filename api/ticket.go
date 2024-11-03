package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/naderSameh/ticketing_support/cache"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
	worker "github.com/naderSameh/ticketing_support/worker"
	"github.com/rs/zerolog/log"
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
	taskPayload := &worker.PayloadSendEmail{
		User:    req.UserAssigned,
		Content: req.Description,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	err = server.taskDistributor.NewEmailDeliveryTask(taskPayload, opts...)
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
	Assigned_to  string `form:"assigned_to"`
	Category_ID  int64  `form:"category_id"`
	PageID       int32  `form:"page_id" binding:"required,min=1"`
	PageSize     int32  `form:"page_size" binding:"required,min=5,max=10"`
	Is_Admin     *bool  `form:"is_admin" binding:"required"`
	User         string `form:"requester" binding:"required"`
}

// ListTickets godoc
//
//	@Summary		List tickets
//	@Description	List all tickets for a specific user, Admin can get all tickets and can add query param to filter by category ID, assigned engineer and ticket owner normal user only can only get all tickets assigned to him
//
//
//	@Tags			Tickets
//
//
//	@Produce		json
//	@Param			user_assigned	query		string	false	"Filter Ticket owner"
//	@Param			page_id			query		int		true	"Page ID"
//	@Param			page_size		query		int		true	"Page Size"
//	@Param			category_id		query		int		false	"Filter Category ID"
//	@Param			assigned_to		query		string	false	"Filter Assigned engineer"
//	@Param			is_admin		query		int		true	"Is admin"
//	@Param			requester		query		string	true	"User sending the request"
//
//	@Success		200				{array}		db.Ticket
//	@Failure		400				{object}	error
//	@Failure		401				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/tickets [get]
func (server *Server) listTicket(c *gin.Context) {
	var req listTicketRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if *req.Is_Admin {
		var userAssignedBool, assignedToBool, categoryIDBool bool
		if req.UserAssigned != "" {
			userAssignedBool = true
		} else if req.Assigned_to != "" {
			assignedToBool = true
		} else if req.Category_ID != 0 {
			categoryIDBool = true
		}
		arg2 := db.ListAllTicketsParams{
			UserAssigned: sql.NullString{
				String: req.UserAssigned,
				Valid:  userAssignedBool,
			},
			AssignedTo: sql.NullString{
				String: req.Assigned_to,
				Valid:  assignedToBool,
			},
			CategoryID: sql.NullInt64{
				Int64: req.Category_ID,
				Valid: categoryIDBool,
			},
			Limit:  req.PageSize,
			Offset: (req.PageID - 1) * req.PageSize,
		}
		tickets, err := server.store.ListAllTickets(c, arg2)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		url := c.FullPath()
		redisClient := cache.NewCacheClient().RedisClient
		// put it to the cache
		content, _ := json.Marshal(tickets)
		err = redisClient.Set(c, url, content, time.Minute*5).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, tickets)
		return
	}
	arg := db.ListAllTicketsParams{
		UserAssigned: sql.NullString{
			String: req.User,
			Valid:  true,
		},
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	tickets, err := server.store.ListAllTickets(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	url := c.FullPath()
	redisClient := cache.NewCacheClient().RedisClient
	// put it to the cache
	content, _ := json.Marshal(tickets)
	err = redisClient.Set(c, url, content, time.Minute*5).Err()
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
//	@Security		ApiKeyAuth
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

type getTicketRequest struct {
	Is_Admin *bool  `form:"is_admin" binding:"required"`
	User     string `form:"requester" binding:"required"`
}

type getTicketRequestURI struct {
	TicketID int64 `uri:"ticket_id" binding:"required,min=1"`
}

// GetTicket godoc
//
//	@Summary		Get ticket by ID
//	@Description	Admins get any ticket, normal user only get a ticket he owns
//	@Tags			Tickets
//
//
//	@Produce		json
//	@Param			ticket_id	path		string	true	"Ticket ID"
//	@Param			is_admin		query		int		true	"Is admin"
//	@Param			requester		query		string	true	"User sending the request"
//
//	@Success		200			{array}		db.Ticket
//	@Failure		401			{object}	error
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/tickets/{ticket_id} [get]
func (server *Server) getTicket(c *gin.Context) {
	var reqURI getTicketRequestURI
	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var req getTicketRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Info().Bool("is_admin", *req.Is_Admin)
	if *req.Is_Admin {
		ticket, err := server.store.GetTicket(c, reqURI.TicketID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, ticket)
		return
	}

	ticket, err := server.store.GetTicket(c, reqURI.TicketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if ticket.UserAssigned != req.User {
		c.JSON(http.StatusUnauthorized, errorResponse(errors.New("user doesn't own that ticket")))
		return
	}
	c.JSON(http.StatusOK, ticket)

}
