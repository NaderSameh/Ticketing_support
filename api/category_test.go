package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/naderSameh/ticketing_support/db/mock"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
	"github.com/stretchr/testify/require"
)

func createRandomCategoryWithName(name string) db.Category {
	return db.Category{
		CategoryID: rand.Int63n(100),
		Name:       name,
	}
}

func TestListCategories(t *testing.T) {

	categories := make([]db.Category, 10)
	for i := 0; i < 10; i++ {
		categories[i] = createRandomCategoryWithName("VIP")
	}

	type Query struct {
		page_id   int32
		page_size int32
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
				page_id:   1,
				page_size: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.ListCategoriesParams{
					Limit:  10,
					Offset: 0,
				}

				store.EXPECT().
					ListCategories(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(categories, nil)

			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "Invalid params",
			query: Query{
				page_id:   -1,
				page_size: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.ListCategoriesParams{
					Limit:  10,
					Offset: 0,
				}

				store.EXPECT().
					ListCategories(gomock.Any(), gomock.Eq(arg)).
					Times(0)

			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		}, {
			name: "internal server error",
			query: Query{
				page_id:   1,
				page_size: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.ListCategoriesParams{
					Limit:  10,
					Offset: 0,
				}

				store.EXPECT().
					ListCategories(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(categories, sql.ErrConnDone)

			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		}, {
			name: "invalid page id",
			query: Query{
				page_id:   -1,
				page_size: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.ListCategoriesParams{
					Limit:  10,
					Offset: 0,
				}

				store.EXPECT().
					ListCategories(gomock.Any(), gomock.Eq(arg)).
					Times(0)

			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
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

			url := "/categories"

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.page_id))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.page_size))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestCreateCategory(t *testing.T) {
	categoryName := "Test Category"
	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request)
		buildstuds    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name": categoryName,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(categoryName)).Times(1).Return(db.Category{CategoryID: 1, Name: categoryName}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name: "Invalid token",
			body: gin.H{
				"name": categoryName,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenInvalid, authorizationTypeBearer)
			},
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(categoryName)).Times(0).Return(db.Category{CategoryID: 1, Name: categoryName}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.True(t, strings.Contains(recorder.Body.String(), "invalid"))
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Unauthorized no permission",
			body: gin.H{
				"name": categoryName,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenNoPermission, authorizationTypeBearer)
			},
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(categoryName)).Times(0).Return(db.Category{CategoryID: 1, Name: categoryName}, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.True(t, strings.Contains(recoder.Body.String(), "Only admins post new categories"))
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "Unauthorized expired token",
			body: gin.H{
				"name": categoryName,
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenExpiration, authorizationTypeBearer)
			},
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(categoryName)).Times(0).Return(db.Category{CategoryID: 1, Name: categoryName}, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.True(t, strings.Contains(recoder.Body.String(), "token is expired"))
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},

		{
			name: "Invalid param",
			body: gin.H{
				"name": 1,
			},
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().CreateCategory(gomock.Any(), gomock.Eq(categoryName)).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name: "Internal server error",
			body: gin.H{
				"name": categoryName,
			},
			buildstuds: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCategory(gomock.Any(), gomock.Eq(categoryName)).
					Times(1).
					Return(db.Category{CategoryID: 1, Name: categoryName}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request) {
				addAuthorization(t, request, JWTtokenOK, authorizationTypeBearer)
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
			tc.buildstuds(store)

			server := newTestServer(t, store)

			url := "/categories"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			tc.setupAuth(t, request)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}
