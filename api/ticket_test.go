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
	mockwk "github.com/naderSameh/ticketing_support/worker/mock"
	"github.com/stretchr/testify/require"
)

const (
	JWTtokenOK                = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJmcmVzaCI6ZmFsc2UsImlhdCI6MTY5NzQ1MzYwNywianRpIjoiOWRkZWI4ZWEtNzI0YS00NGYyLWIxMDgtODEzZjFkY2RmNjIxIiwidHlwZSI6ImFjY2VzcyIsInN1YiI6Ik1pbmEzIiwibmJmIjoxNjk3NDUzNjA3LCJleHAiOjE3OTc0NjQ0MDcsInBlcm1pc3Npb25zIjpbImFwaV92Mi5kZXZpY2VzLmFsbC1yZXF1ZXN0cy5HRVQiLCJhcGlfdjIuZGV2aWNlcy5yZXF1ZXN0cy5QVVQiLCJ0aWNrZXRzLlBPU1QiLCJ0aWNrZXRzLkdFVCIsInRpY2tldHMuUFVUIiwidGlja2V0cy5EZWxldGUiXX0.S-QPmsHobV0mdQ5CyBlfV0nT2-nMMmveTJwM9y5VeTA"
	JWTtokenNoPermission      = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJmcmVzaCI6ZmFsc2UsImlhdCI6MTY5NzQ1MzYwNywianRpIjoiOWRkZWI4ZWEtNzI0YS00NGYyLWIxMDgtODEzZjFkY2RmNjIxIiwidHlwZSI6ImFjY2VzcyIsInN1YiI6Ik1pbmE0IiwibmJmIjoxNjk3NDUzNjA3LCJleHAiOjE3OTc0NjQ0MDcsInBlcm1pc3Npb25zIjpbImFwaV92Mi5kZXZpY2VzLmFsbC1yZXF1ZXN0cy5HRVQiLCJhcGlfdjIuZGV2aWNlcy5yZXF1ZXN0cy5QVVQiLCJ0aWNrZXRzLkRlbGV0ZSJdfQ.ha5EdM4ngi8UOvzO3gOqjatMDyoDPX797atupQ-L84I"
	JWTtokenNoPermissionMina3 = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJmcmVzaCI6ZmFsc2UsImlhdCI6MTY5NzQ1MzYwNywianRpIjoiOWRkZWI4ZWEtNzI0YS00NGYyLWIxMDgtODEzZjFkY2RmNjIxIiwidHlwZSI6ImFjY2VzcyIsInN1YiI6Ik1pbmEzIiwibmJmIjoxNjk3NDUzNjA3LCJleHAiOjE3OTc0NjQ0MDcsInBlcm1pc3Npb25zIjpbImFwaV92Mi5kZXZpY2VzLmFsbC1yZXF1ZXN0cy5HRVQiLCJhcGlfdjIuZGV2aWNlcy5yZXF1ZXN0cy5QVVQiLCJ0aWNrZXRzLkRlbGV0ZSJdfQ.owoMfnpWyWb-7Dabfv9n8OJdIy7VMxSPKaDcd67j3A8"
	JWTtokenExpiration        = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJmcmVzaCI6ZmFsc2UsImlhdCI6MTY5NzQ1MzYwNywianRpIjoiOWRkZWI4ZWEtNzI0YS00NGYyLWIxMDgtODEzZjFkY2RmNjIxIiwidHlwZSI6ImFjY2VzcyIsInN1YiI6Ik1pbmEzIiwibmJmIjoxNjk3NDUzNjA3LCJleHAiOjE2OTc0NjQ0MDcsInBlcm1pc3Npb25zIjpbImFwaV92Mi5kZXZpY2VzLmFsbC1yZXF1ZXN0cy5HRVQiLCJhcGlfdjIuZGV2aWNlcy5yZXF1ZXN0cy5QVVQiLCJ0aWNrZXRzLlBPU1QiLCJ0aWNrZXRzLkdFVCIsInRpY2tldHMuUFVUIiwidGlja2V0cy5EZWxldGUiXX0._eu-oeHopfpSfSa0HZJyl5oxuCw6Q_h8O8mlmn2ascc"
	JWTtokenInvalid           = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJmcmVzaCI6ZmFsc2UsImlhdCI6MTY5NzQ1MzYwNywianRpIjoiOWRkZWI4ZWEtNzI0YS00NGYyLWIxMDgtODEzZjFkY2RmNjIxIiwidHlwZSI6ImFjY2VzcyIsInN1YiI6Ik1pbmEzIiwibmJmIjoxNjk3NDUzNjA3LCJleHAiOjE2OTc0NjQ0MDcsInBlcm1pc3Npb25zIjpbImFwaV92Mi5kZXZpY2VzLmFsbC1yZXF1ZXN0cy5HRVQiLCJhcGlfdjIuZGV2aWNlcy5yZXF1ZXN0cy5QVVQiLCJ0aWNrZXRzLkdFVCIsInRpY2tldHMuRGVsZXRlIl19.dQsLRbYndv_vT3N89sYMVcv3Hs5xrgR5BolVj4O1D4A"
	JWTtokenInvalidAlg        = "eyJ0eXAiOiJKV1QiLCJhbGciOiJQUzI1NiJ9.eyJmcmVzaCI6ZmFsc2UsImlhdCI6MTY5NzQ1MzYwNywianRpIjoiOWRkZWI4ZWEtNzI0YS00NGYyLWIxMDgtODEzZjFkY2RmNjIxIiwidHlwZSI6ImFjY2VzcyIsInN1YiI6Ik1pbmEzIiwibmJmIjoxNjk3NDUzNjA3LCJleHAiOjE2OTc0NjQ0MDcsInBlcm1pc3Npb25zIjpbImFwaV92Mi5kZXZpY2VzLmFsbC1yZXF1ZXN0cy5HRVQiLCJhcGlfdjIuZGV2aWNlcy5yZXF1ZXN0cy5QVVQiLCJ0aWNrZXRzLkdFVCIsInRpY2tldHMuRGVsZXRlIl19.wnAfwYO-yXDR9YXd_TmbV-Xo9PO73G-6fccVAXAj9Lff7zW6183vLDNidoy44bmAzJwJdtQoX-GIlbP2fbD5vFpZxaRbztGh7CpoIUarT4x7IbhUTKzsJfY-YQqU_NDtitpeWTNHHTZJnl6nTJw0Qj1iHDz5p3OxjKxohpkWJ7fKSDvwqgKLxpLhJN3cvt7J6GePfILaPvHv80pw1zW5dhB3MmvD-H8smrmy2iawh9j-m6NtxaTXmBVIIdDBUiatz_vglDbuFZkOHiCtb73_sYECA9ng8YOCNCFYU7sxlLPsk55STSBbFIp3xeUJNC3HiYuc_mouWvPYwPe6XoyA4w"
	JWTtokenInvalidPayload    = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJmcmVzaCI6ZmFsc2UsImlhdCI6MTY3NjI4NzI2OSwianRpIjoiYzcyZGY4MzctZTEwZi00NTFkLTkzM2QtNjdkODdiNjk2ZTJlIiwidHlwZSI6ImFjY2VzcyIsInN1YiI6Imdlb3JnZWhlc2htZXQiLCJuYmYiOjE2NzYyODcyNjksInJvbGVzIjpbeyJpZCI6MSwibmFtZSI6ImN5Y29sbGVjdG9yX2FjdGlvbnRha2VyX2xvY2siLCJkZXNjcmlwdGlvbiI6ImNhbiB2aWV3IGN5Y29sbGVjdG9yIGFuZCBsb2NrL3VubG9jayBkZXZpY2UifSx7ImlkIjoyLCJuYW1lIjoiY3l0YWdzX2FjdGlvbl90YWtlciIsImRlc2NyaXB0aW9uIjoiY2FuIGRvIGFsbCBjeXRhZyBhY3Rpb25zIn0seyJpZCI6NCwibmFtZSI6InVzZXJfZWRpdG9yIiwiZGVzY3JpcHRpb24iOiJlZGl0IGFuZCBwb3N0IHVzZXIifV0sInBlcm1pc3Npb25zIjpbeyJzY29wZSI6ImxvY2tjeWNvbGxlY3Rvci5QT1NUIiwiaWQiOjg5ODQsImRlc2NyaXB0aW9uIjoiUE9TVCAgL2N5Y29sbGVjdG9yL2xvY2sifSx7InNjb3BlIjoiYXBpX3YyLmFjdGlvbnMubG9jay5QT1NUIiwiaWQiOjkwNTcsImRlc2NyaXB0aW9uIjoiUE9TVCAgL2FwaS12Mi9hY3Rpb25zL2N5bG9jay9sb2NrLWFjdGlvbiJ9LHsic2NvcGUiOiJjeWNvbGxlY3RvcnMtYXNzaWduLlBVVCIsImlkIjo5MDA1LCJkZXNjcmlwdGlvbiI6IlBVVCAgL2N5Y29sbGVjdG9ycy1hc3NpZ24vPGludDpjeWNvbGxlY3Rvcl9pZD4ifV19.UvDC0rjxD9hJifj0ejWWlVB4uDcVZBeQZAntqoxSvM4"
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
		buildStubs    func(store *mockdb.MockStore, worker *mockwk.MockTaskDistributor)
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
			buildStubs: func(store *mockdb.MockStore, workermock *mockwk.MockTaskDistributor) {
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

				workermock.EXPECT().NewEmailDeliveryTask(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Missing data",
			body: gin.H{
				"title":         ticket.Title,
				"description":   ticket.Description,
				"status":        ticket.Status,
				"user_assigned": ticket.UserAssigned,
				// "category_id":   ticket.CategoryID,
			},
			buildStubs: func(store *mockdb.MockStore, workermock *mockwk.MockTaskDistributor) {
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
				workermock.EXPECT().NewEmailDeliveryTask(gomock.Any(), gomock.Any()).Times(0)
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
			buildStubs: func(store *mockdb.MockStore, workermock *mockwk.MockTaskDistributor) {
				store.EXPECT().
					CreateTicket(gomock.Any(), gomock.Any()).
					Times(0)
				workermock.EXPECT().NewEmailDeliveryTask(gomock.Any(), gomock.Any()).Times(0)
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
			buildStubs: func(store *mockdb.MockStore, workermock *mockwk.MockTaskDistributor) {
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
				workermock.EXPECT().NewEmailDeliveryTask(gomock.Any(), gomock.Any()).Times(0)
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

			ctrl2 := gomock.NewController(t)
			defer ctrl2.Finish()

			store := mockdb.NewMockStore(ctrl)
			taskDistributor := mockwk.NewMockTaskDistributor(ctrl2)
			tc.buildStubs(store, taskDistributor)

			server := newTestServer(t, store, taskDistributor)
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

			server := newTestServer(t, store, nil)
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

func TestListAllTicket(t *testing.T) {
	n := 10
	tickets := make([]db.Ticket, n)
	for i := 0; i < n; i++ {
		tickets[i] = createRandomTicket("Mina3")
	}
	type Query struct {
		user_assigned string
		assigned_to   string
		category_id   int64
		page_id       int32
		page_size     int32
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK admin",
			query: Query{
				user_assigned: tickets[0].UserAssigned,
				page_id:       1,
				page_size:     10,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAllTicketsParams{
					UserAssigned: sql.NullString{
						String: tickets[0].UserAssigned,
						Valid:  true,
					},
					Limit:  10,
					Offset: 0,
				}

				store.EXPECT().
					ListAllTickets(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(tickets, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OK admin filter by assigned to",
			query: Query{
				assigned_to: "engineer",
				page_id:     1,
				page_size:   10,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAllTicketsParams{
					AssignedTo: sql.NullString{
						String: "engineer",
						Valid:  true,
					},
					UserAssigned: sql.NullString{
						String: "",
						Valid:  false,
					},
					CategoryID: sql.NullInt64{
						Int64: 0,
						Valid: false,
					},
					Limit:  10,
					Offset: 0,
				}

				store.EXPECT().
					ListAllTickets(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(tickets, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		}, {
			name: "OK admin filter by category",
			query: Query{
				category_id: tickets[0].CategoryID,
				page_id:     1,
				page_size:   10,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAllTicketsParams{
					AssignedTo: sql.NullString{
						String: "",
						Valid:  false,
					},
					UserAssigned: sql.NullString{
						String: "",
						Valid:  false,
					},
					CategoryID: sql.NullInt64{
						Int64: tickets[0].CategoryID,
						Valid: true,
					},
					Limit:  10,
					Offset: 0,
				}

				store.EXPECT().
					ListAllTickets(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(tickets, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OK user",
			query: Query{
				user_assigned: tickets[0].UserAssigned,
				page_id:       1,
				page_size:     10,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenNoPermissionMina3, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAllTicketsParams{
					UserAssigned: sql.NullString{
						String: tickets[0].UserAssigned,
						Valid:  true,
					},
					Limit:  10,
					Offset: 0,
				}

				store.EXPECT().
					ListAllTickets(gomock.Any(), gomock.Eq(arg)).
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
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAllTicketsParams{
					UserAssigned: sql.NullString{
						String: tickets[0].UserAssigned,
						Valid:  true,
					},
					Limit:  10,
					Offset: 0,
				}

				store.EXPECT().
					ListAllTickets(gomock.Any(), gomock.Eq(arg)).
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
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAllTickets(gomock.Any(), gomock.Any()).
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
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAllTickets(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "internal server error admin",
			query: Query{
				user_assigned: tickets[0].UserAssigned,
				page_id:       1,
				page_size:     10,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAllTickets(gomock.Any(), gomock.Any()).
					Times(1).Return(tickets, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "internal server error user",
			query: Query{
				user_assigned: tickets[0].UserAssigned,
				page_id:       1,
				page_size:     10,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenNoPermission, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAllTickets(gomock.Any(), gomock.Any()).
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

			server := newTestServer(t, store, nil)
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
			q.Add("assigned_to", fmt.Sprintf("%s", tc.query.assigned_to))
			q.Add("category_id", fmt.Sprintf("%d", tc.query.category_id))
			request.URL.RawQuery = q.Encode()
			tc.setupAuth(t, request)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
func TestGetTicket(t *testing.T) {

	ticket := createRandomTicket("Mina3")

	testCases := []struct {
		name          string
		ticket_id     int64
		setupAuth     func(t *testing.T, request *http.Request)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK admin",
			ticket_id: ticket.TicketID,
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetTicket(gomock.Any(), ticket.TicketID).
					Times(1).
					Return(ticket, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "OK user",
			ticket_id: ticket.TicketID,
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenNoPermissionMina3, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetTicket(gomock.Any(), ticket.TicketID).
					Times(1).
					Return(ticket, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		}, {
			name: "Missing ticket ID",
			// ticket_id: ticket.TicketID,
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenNoPermissionMina3, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetTicket(gomock.Any(), ticket.TicketID).
					Times(0).
					Return(ticket, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Unauthorized request",
			ticket_id: ticket.TicketID,
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenNoPermission, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetTicket(gomock.Any(), ticket.TicketID).
					Times(1).
					Return(ticket, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name:      "internal server error admin",
			ticket_id: ticket.TicketID,
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetTicket(gomock.Any(), ticket.TicketID).
					Times(1).
					Return(ticket, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "internal server error user",
			ticket_id: ticket.TicketID,
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenNoPermissionMina3, authorizationTypeBearer)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetTicket(gomock.Any(), ticket.TicketID).
					Times(1).
					Return(ticket, sql.ErrConnDone)
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

			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON

			url := fmt.Sprintf("/tickets/%d", tc.ticket_id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request)
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
		setupAuth     func(t *testing.T, request *http.Request)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			TicketID: ticket.TicketID,
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
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
		},
		{
			name:     "No authorization",
			TicketID: ticket.TicketID,
			setupAuth: func(t *testing.T, request *http.Request) {
				// addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			body: gin.H{
				"ticket_id":   ticket.TicketID,
				"status":      "closed",
				"assigned_to": "someone",
			},
			buildStubs: func(store *mockdb.MockStore) {

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "Unauthorized no permission",
			TicketID: ticket.TicketID,
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenNoPermission, authorizationTypeBearer)
			},
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
					Times(0).
					Return(ticket, nil)

				store.EXPECT().
					UpdateTicket(gomock.Any(), gomock.Eq(arg)).
					Times(0).
					Return(ticket, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:     "Invalid URI ticket ID",
			TicketID: -ticket.TicketID,
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			body: gin.H{
				"ticket_id":   -ticket.TicketID,
				"status":      "closed",
				"assigned_to": "someone",
			},
			buildStubs: func(store *mockdb.MockStore) {

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "missing data",
			TicketID: ticket.TicketID,
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
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
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
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
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
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
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
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
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
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

			server := newTestServer(t, store, nil)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/tickets/%d", tc.TicketID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request)

			server.router.ServeHTTP(recorder, request)
			fmt.Printf("error %s", recorder.Body.String())
			tc.checkResponse(recorder)
		})
	}
}
