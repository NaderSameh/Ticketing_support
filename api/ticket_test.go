package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/naderSameh/ticketing_support/db/mock"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
	"github.com/naderSameh/ticketing_support/util"
	"github.com/stretchr/testify/require"
)

func createRandomCategory() db.Category {
	return db.Category{
		CategoryID: rand.Int63n(100),
		Name:       util.GenerateRandomString(10),
	}
}

func createRandomTicket(user_assigned string) db.Ticket {
	category := createRandomCategory()
	return db.Ticket{
		TicketID:     rand.Int63n(100),
		Title:        util.GenerateRandomString(10),
		Description:  util.GenerateRandomString(40),
		Status:       "open",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CategoryID:   category.CategoryID,
		UserAssigned: user_assigned,
	}
}

func TestCreateTicket(t *testing.T) {

	ticket := createRandomTicket(util.GenerateRandomString(10))
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"title":         ticket.Title,
				"description":   ticket.Description,
				"status":        ticket.Status,
				"user_assigned": ticket.UserAssigned,
				"category_id":   ticket.CategoryID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTicketParams{
					Title:        ticket.Title,
					Description:  ticket.Description,
					Status:       ticket.Status,
					UserAssigned: ticket.UserAssigned,
					CategoryID:   ticket.CategoryID,
				}

				store.EXPECT().
					CreateTicket(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(ticket, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		}, {
			name: "Missing data",
			body: gin.H{
				"title":         ticket.Title,
				"description":   ticket.Description,
				"status":        ticket.Status,
				"user_assigned": ticket.UserAssigned,
				// "category_id":   ticket.CategoryID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTicketParams{
					Title:        ticket.Title,
					Description:  ticket.Description,
					Status:       ticket.Status,
					UserAssigned: ticket.UserAssigned,
					CategoryID:   ticket.CategoryID,
				}

				store.EXPECT().
					CreateTicket(gomock.Any(), gomock.Eq(arg)).
					Times(0).
					Return(ticket, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not valid status",
			body: gin.H{
				"title":         ticket.Title,
				"description":   ticket.Description,
				"status":        "finished",
				"user_assigned": ticket.UserAssigned,
				"category_id":   ticket.CategoryID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateTicket(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			body: gin.H{
				"title":         ticket.Title,
				"description":   ticket.Description,
				"status":        ticket.Status,
				"user_assigned": ticket.UserAssigned,
				"category_id":   ticket.CategoryID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTicketParams{
					Title:        ticket.Title,
					Description:  ticket.Description,
					Status:       ticket.Status,
					UserAssigned: ticket.UserAssigned,
					CategoryID:   ticket.CategoryID,
				}

				store.EXPECT().
					CreateTicket(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(ticket, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/tickets"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteTicket(t *testing.T) {

	ticket := createRandomTicket(util.GenerateRandomString(10))
	type Query struct {
		ticket_id int64
	}
	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				ticket_id: ticket.TicketID,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(1)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},

		{
			name: "Invalid request",
			query: Query{
				ticket_id: -ticket.TicketID,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "non existing ticket",
			query: Query{
				ticket_id: 1000000,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(1).Return(ticket, sql.ErrNoRows)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "internal server error",
			query: Query{
				ticket_id: ticket.TicketID,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(1).Return(ticket, sql.ErrConnDone)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Any()).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/tickets/%d", tc.query.ticket_id)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("ticket_id", fmt.Sprintf("%d", tc.query.ticket_id))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListTicket(t *testing.T) {
	n := 10
	tickets := make([]db.Ticket, n)
	for i := 0; i < n; i++ {
		tickets[i] = createRandomTicket("martin")
	}
	type Query struct {
		user_assigned string
		page_id       int32
		page_size     int32
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				user_assigned: tickets[0].UserAssigned,
				page_id:       1,
				page_size:     10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTicketsParams{
					UserAssigned: tickets[0].UserAssigned,
					Limit:        10,
					Offset:       0,
				}

				store.EXPECT().
					ListTickets(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(tickets, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Missing data",
			query: Query{
				user_assigned: tickets[0].UserAssigned,
				// "page_id":       1,
				// "page_size":     10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTicketsParams{
					UserAssigned: tickets[0].UserAssigned,
					Limit:        10,
					Offset:       0,
				}

				store.EXPECT().
					ListTickets(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "too big page size",
			query: Query{
				user_assigned: tickets[0].UserAssigned,
				page_id:       1,
				page_size:     100,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTickets(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name: "invalid page id",
			query: Query{
				user_assigned: tickets[0].UserAssigned,
				page_id:       -1,
				page_size:     10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTickets(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "internal server error",
			query: Query{
				user_assigned: tickets[0].UserAssigned,
				page_id:       1,
				page_size:     10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListTickets(gomock.Any(), gomock.Any()).
					Times(1).Return(tickets, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			// data, err := json.Marshal(tc.body)
			// require.NoError(t, err)

			url := "/tickets"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("user_assigned", fmt.Sprintf("%s", tc.query.user_assigned))
			q.Add("page_id", fmt.Sprintf("%d", tc.query.page_id))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.page_size))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
func TestUpdateTicket(t *testing.T) {

	ticket := createRandomTicket("martin")

	testCases := []struct {
		name          string
		body          gin.H
		TicketID      int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			TicketID: ticket.TicketID,
			body: gin.H{
				"ticket_id":   ticket.TicketID,
				"status":      "closed",
				"assigned_to": "someone",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateTicketParams{
					TicketID:   ticket.TicketID,
					UpdatedAt:  time.Now().Round(time.Second),
					Status:     "closed",
					AssignedTo: sql.NullString{String: "someone", Valid: true},
				}
				ticket.UpdatedAt = time.Now()
				ticket.Status = "closed"

				store.EXPECT().
					GetTicketForUpdate(gomock.Any(), ticket.TicketID).
					Times(1).
					Return(ticket, nil)

				store.EXPECT().
					UpdateTicket(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(ticket, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		}, {
			name:     "missing data",
			TicketID: ticket.TicketID,
			body: gin.H{
				"ticket_id": ticket.TicketID,
				// "status":      "closed",
				"assigned_to": "someone",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateTicket(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name: "Invalid status",
			body: gin.H{
				"ticket_id":   ticket.TicketID,
				"status":      "finished",
				"assigned_to": "someone",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateTicket(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "Internal server error",
			TicketID: ticket.TicketID,
			body: gin.H{
				"ticket_id":   ticket.TicketID,
				"status":      "closed",
				"assigned_to": "someone",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateTicketParams{
					TicketID:   ticket.TicketID,
					UpdatedAt:  time.Now().Round(time.Second),
					Status:     "closed",
					AssignedTo: sql.NullString{String: "someone", Valid: true},
				}
				ticket.UpdatedAt = time.Now()
				ticket.Status = "closed"

				store.EXPECT().
					GetTicketForUpdate(gomock.Any(), ticket.TicketID).
					Times(1).
					Return(ticket, nil)

				store.EXPECT().
					UpdateTicket(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(ticket, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name:     "Empty assigned to field",
			TicketID: ticket.TicketID,
			body: gin.H{
				"ticket_id": ticket.TicketID,
				"status":    "closed",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateTicketParams{
					TicketID:   ticket.TicketID,
					UpdatedAt:  time.Now().Round(time.Second),
					Status:     "closed",
					AssignedTo: sql.NullString{String: "", Valid: false},
				}
				ticket.UpdatedAt = time.Now()
				ticket.Status = "closed"

				store.EXPECT().
					GetTicketForUpdate(gomock.Any(), ticket.TicketID).
					Times(1).
					Return(ticket, nil)

				store.EXPECT().
					UpdateTicket(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(ticket, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		}, {
			name:     "Not an existing ticket",
			TicketID: ticket.TicketID,
			body: gin.H{
				"ticket_id": ticket.TicketID,
				"status":    "closed",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateTicketParams{
					TicketID:   ticket.TicketID,
					UpdatedAt:  time.Now().Round(time.Second),
					Status:     "closed",
					AssignedTo: sql.NullString{String: "", Valid: false},
				}
				ticket.UpdatedAt = time.Now()
				ticket.Status = "closed"

				store.EXPECT().
					GetTicketForUpdate(gomock.Any(), ticket.TicketID).
					Times(1).
					Return(ticket, sql.ErrNoRows)

				store.EXPECT().
					UpdateTicket(gomock.Any(), gomock.Eq(arg)).
					Times(0).
					Return(ticket, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/tickets/%d", tc.TicketID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
