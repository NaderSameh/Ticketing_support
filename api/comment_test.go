package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
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

func createRandomComment() db.Comment {
	return db.Comment{
		CommentID:     1,
		CommentText:   util.GenerateRandomString(50),
		CreatedAt:     time.Now(),
		TicketID:      1,
		UserCommented: util.GenerateRandomString(10),
	}
}
func TestCreateComment(t *testing.T) {
	comment := createRandomComment()
	testCases := []struct {
		name          string
		ticket_id     int64
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			ticket_id: comment.TicketID,
			body: gin.H{
				"comment_text":   comment.CommentText,
				"user_commented": comment.UserCommented,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateCommentParams{
					CommentText:   comment.CommentText,
					TicketID:      comment.TicketID,
					UserCommented: comment.UserCommented,
				}

				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(comment, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name:      "Invalid ticket ID",
			ticket_id: -comment.TicketID,
			body: gin.H{
				"comment_text":   comment.CommentText,
				"user_commented": comment.UserCommented,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateCommentParams{
					CommentText:   comment.CommentText,
					TicketID:      comment.TicketID,
					UserCommented: comment.UserCommented,
				}

				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:      "Missing fields in body",
			ticket_id: comment.TicketID,
			body: gin.H{
				"comment_text": comment.CommentText,
				// "user_commented": comment.UserCommented,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateCommentParams{
					CommentText:   comment.CommentText,
					TicketID:      comment.TicketID,
					UserCommented: comment.UserCommented,
				}

				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:      "Internal server error",
			ticket_id: comment.TicketID,
			body: gin.H{
				"comment_text":   comment.CommentText,
				"user_commented": comment.UserCommented,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateCommentParams{
					CommentText:   comment.CommentText,
					TicketID:      comment.TicketID,
					UserCommented: comment.UserCommented,
				}

				store.EXPECT().
					CreateComment(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(comment, sql.ErrConnDone)
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
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/tickets/%d/comments", tc.ticket_id)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	comment := createRandomComment()

	testCases := []struct {
		name          string
		ticket_id     int64
		comment_id    int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			ticket_id:  comment.TicketID,
			comment_id: comment.CommentID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetCommentForUpdate(gomock.Any(), comment.CommentID).Times(1).Return(comment, nil)

				store.EXPECT().
					DeleteComment(gomock.Any(), comment.CommentID).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		}, {
			name:       "invalid ticket ID",
			ticket_id:  comment.TicketID,
			comment_id: -comment.CommentID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetCommentForUpdate(gomock.Any(), comment.CommentID).Times(0).Return(comment, nil)

				store.EXPECT().
					DeleteComment(gomock.Any(), comment.CommentID).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:       "invalid ticket ID",
			ticket_id:  comment.TicketID,
			comment_id: -comment.CommentID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetCommentForUpdate(gomock.Any(), comment.CommentID).Times(0).Return(comment, nil)

				store.EXPECT().
					DeleteComment(gomock.Any(), comment.CommentID).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:       "comment not found",
			ticket_id:  comment.TicketID,
			comment_id: comment.CommentID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetCommentForUpdate(gomock.Any(), comment.CommentID).Times(1).Return(comment, sql.ErrNoRows)

				store.EXPECT().
					DeleteComment(gomock.Any(), comment.CommentID).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		}, {
			name:       "internal server error",
			ticket_id:  comment.TicketID,
			comment_id: comment.CommentID,
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetCommentForUpdate(gomock.Any(), comment.CommentID).Times(1).Return(comment, nil)

				store.EXPECT().
					DeleteComment(gomock.Any(), comment.CommentID).Times(1).Return(sql.ErrConnDone)
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

			url := fmt.Sprintf("/tickets/%d/comments/%d", tc.ticket_id, tc.comment_id)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
func TestListComment(t *testing.T) {
	comments := make([]db.Comment, 10)
	for i := 0; i < 10; i++ {
		comments[i] = createRandomComment() //fixed ticket id = 1
	}
	ticket := createRandomTicket(util.GenerateRandomString(10))
	type Query struct {
		PageID   int32
		PageSize int32
	}
	testCases := []struct {
		name          string
		query         Query
		ticket_id     int64
		comment_id    int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			ticket_id:  comments[0].TicketID,
			comment_id: comments[0].CommentID,
			query: Query{
				PageID:   1,
				PageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.ListCommentsParams{
					TicketID: comments[0].TicketID,
					Limit:    10,
					Offset:   0,
				}
				store.EXPECT().
					GetTicket(gomock.Any(), gomock.Any()).Times(1).Return(ticket, nil)
				store.EXPECT().
					ListComments(gomock.Any(), arg).Times(1).Return(comments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		}, {
			name:       "Invalid page ID",
			ticket_id:  comments[0].TicketID,
			comment_id: comments[0].CommentID,
			query: Query{
				PageID:   0,
				PageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					ListComments(gomock.Any(), gomock.Any()).Times(0).Return(comments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:       "Invalid page size",
			ticket_id:  comments[0].TicketID,
			comment_id: comments[0].CommentID,
			query: Query{
				PageID:   1,
				PageSize: 20,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					ListComments(gomock.Any(), gomock.Any()).Times(0).Return(comments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:       "Invalid Ticket ID",
			ticket_id:  -1,
			comment_id: comments[0].CommentID,
			query: Query{
				PageID:   1,
				PageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTicket(gomock.Any(), gomock.Any()).Times(0).Return(ticket, nil)
				store.EXPECT().
					ListComments(gomock.Any(), gomock.Any()).Times(0).Return(comments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:       "Ticket not found",
			ticket_id:  comments[0].TicketID,
			comment_id: comments[0].CommentID,
			query: Query{
				PageID:   1,
				PageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTicket(gomock.Any(), gomock.Any()).Times(1).Return(ticket, sql.ErrNoRows)
				store.EXPECT().
					ListComments(gomock.Any(), gomock.Any()).Times(0).Return(comments, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		}, {
			name:       "Internal server error",
			ticket_id:  comments[0].TicketID,
			comment_id: comments[0].CommentID,
			query: Query{
				PageID:   1,
				PageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetTicket(gomock.Any(), gomock.Any()).Times(1).Return(ticket, nil)
				store.EXPECT().
					ListComments(gomock.Any(), gomock.Any()).Times(1).Return(comments, sql.ErrConnDone)
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

			url := fmt.Sprintf("/tickets/%d/comments", tc.ticket_id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateComment(t *testing.T) {
	ticket := createRandomTicket(util.GenerateRandomString(10))
	comment := createRandomComment()
	comment_text := util.GenerateRandomString(40)
	testCases := []struct {
		name          string
		ticket_id     int64
		comment_id    int64
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			ticket_id:  ticket.TicketID,
			comment_id: comment.CommentID,
			body: gin.H{
				"comment_text": comment_text,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.UpdateCommentParams{
					CommentID:   comment.CommentID,
					CommentText: comment_text,
				}

				arg2 := db.UpdateTicketParams{
					TicketID:  ticket.TicketID,
					UpdatedAt: time.Now().Round(time.Second),
				}

				store.EXPECT().
					GetTicketForUpdate(gomock.Any(), ticket.TicketID).Times(1).Return(ticket, nil)
				store.EXPECT().
					GetCommentForUpdate(gomock.Any(), comment.CommentID).Times(1).Return(comment, nil)

				store.EXPECT().
					UpdateComment(gomock.Any(), arg).Times(1).Return(comment, nil)
				store.EXPECT().
					UpdateTicket(gomock.Any(), arg2).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		}, {
			name:       "invalid ticket ID",
			ticket_id:  -ticket.TicketID,
			comment_id: comment.CommentID,
			body: gin.H{
				"comment_text": comment_text,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.UpdateCommentParams{
					CommentID:   comment.CommentID,
					CommentText: comment_text,
				}

				arg2 := db.UpdateTicketParams{
					TicketID:  ticket.TicketID,
					UpdatedAt: time.Now().Round(time.Second),
				}

				store.EXPECT().
					GetTicketForUpdate(gomock.Any(), ticket.TicketID).Times(0).Return(ticket, nil)
				store.EXPECT().
					GetCommentForUpdate(gomock.Any(), comment.CommentID).Times(0).Return(comment, nil)

				store.EXPECT().
					UpdateComment(gomock.Any(), arg).Times(0).Return(comment, nil)
				store.EXPECT().
					UpdateTicket(gomock.Any(), arg2).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:       "ticket not found",
			ticket_id:  ticket.TicketID,
			comment_id: comment.CommentID,
			body: gin.H{
				"comment_text": comment_text,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.UpdateCommentParams{
					CommentID:   comment.CommentID,
					CommentText: comment_text,
				}

				arg2 := db.UpdateTicketParams{
					TicketID:  ticket.TicketID,
					UpdatedAt: time.Now().Round(time.Second),
				}

				store.EXPECT().
					GetTicketForUpdate(gomock.Any(), ticket.TicketID).Times(1).Return(ticket, sql.ErrNoRows)
				store.EXPECT().
					GetCommentForUpdate(gomock.Any(), comment.CommentID).Times(0).Return(comment, nil)

				store.EXPECT().
					UpdateComment(gomock.Any(), arg).Times(0).Return(comment, nil)
				store.EXPECT().
					UpdateTicket(gomock.Any(), arg2).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		}, {
			name:       "comment not found",
			ticket_id:  ticket.TicketID,
			comment_id: comment.CommentID,
			body: gin.H{
				"comment_text": comment_text,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.UpdateCommentParams{
					CommentID:   comment.CommentID,
					CommentText: comment_text,
				}

				arg2 := db.UpdateTicketParams{
					TicketID:  ticket.TicketID,
					UpdatedAt: time.Now().Round(time.Second),
				}

				store.EXPECT().
					GetTicketForUpdate(gomock.Any(), ticket.TicketID).Times(1).Return(ticket, nil)
				store.EXPECT().
					GetCommentForUpdate(gomock.Any(), comment.CommentID).Times(1).Return(comment, sql.ErrNoRows)

				store.EXPECT().
					UpdateComment(gomock.Any(), arg).Times(0).Return(comment, nil)
				store.EXPECT().
					UpdateTicket(gomock.Any(), arg2).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		}, {
			name:       "comment text not found in body",
			ticket_id:  ticket.TicketID,
			comment_id: comment.CommentID,
			body:       gin.H{
				// "comment_text": comment_text,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.UpdateCommentParams{
					CommentID:   comment.CommentID,
					CommentText: comment_text,
				}

				arg2 := db.UpdateTicketParams{
					TicketID:  ticket.TicketID,
					UpdatedAt: time.Now().Round(time.Second),
				}

				store.EXPECT().
					GetTicketForUpdate(gomock.Any(), ticket.TicketID).Times(0).Return(ticket, nil)
				store.EXPECT().
					GetCommentForUpdate(gomock.Any(), comment.CommentID).Times(0).Return(comment, nil)

				store.EXPECT().
					UpdateComment(gomock.Any(), arg).Times(0).Return(comment, nil)
				store.EXPECT().
					UpdateTicket(gomock.Any(), arg2).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:       "internal server error",
			ticket_id:  ticket.TicketID,
			comment_id: comment.CommentID,
			body: gin.H{
				"comment_text": comment_text,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.UpdateCommentParams{
					CommentID:   comment.CommentID,
					CommentText: comment_text,
				}

				arg2 := db.UpdateTicketParams{
					TicketID:  ticket.TicketID,
					UpdatedAt: time.Now().Round(time.Second),
				}

				store.EXPECT().
					GetTicketForUpdate(gomock.Any(), ticket.TicketID).Times(1).Return(ticket, sql.ErrConnDone)
				store.EXPECT().
					GetCommentForUpdate(gomock.Any(), comment.CommentID).Times(1).Return(comment, sql.ErrConnDone)

				store.EXPECT().
					UpdateComment(gomock.Any(), arg).Times(1).Return(comment, sql.ErrConnDone)
				store.EXPECT().
					UpdateTicket(gomock.Any(), arg2).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		}, {
			name:       "internal server error",
			ticket_id:  ticket.TicketID,
			comment_id: comment.CommentID,
			body: gin.H{
				"comment_text": comment_text,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.UpdateCommentParams{
					CommentID:   comment.CommentID,
					CommentText: comment_text,
				}

				arg2 := db.UpdateTicketParams{
					TicketID:  ticket.TicketID,
					UpdatedAt: time.Now().Round(time.Second),
				}

				store.EXPECT().
					GetTicketForUpdate(gomock.Any(), ticket.TicketID).Times(1).Return(ticket, sql.ErrConnDone)
				store.EXPECT().
					GetCommentForUpdate(gomock.Any(), comment.CommentID).Times(1).Return(comment, sql.ErrConnDone)

				store.EXPECT().
					UpdateComment(gomock.Any(), arg).Times(1).Return(comment, nil)
				store.EXPECT().
					UpdateTicket(gomock.Any(), arg2).Times(1).Return(ticket, sql.ErrConnDone)
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
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/tickets/%d/comments/%d", tc.ticket_id, tc.comment_id)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
