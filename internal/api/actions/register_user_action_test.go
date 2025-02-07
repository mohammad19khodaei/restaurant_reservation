package actions_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bxcodec/faker/v3"
	mockdb "github.com/mohammad19khodaei/restaurant_reservation/db/mock"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/application"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type registerRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestRegisterUserAction(t *testing.T) {
	testCases := []struct {
		name          string
		requestBody   registerRequestBody
		buildStubs    func(repository *mockdb.UserMockRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody registerRequestBody)
	}{
		{
			name:        "without request body",
			requestBody: registerRequestBody{},
			buildStubs: func(repository *mockdb.UserMockRepository) {
				repository.EXPECT().Register(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody registerRequestBody) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "without username",
			requestBody: registerRequestBody{
				Password: faker.Password(),
			},
			buildStubs: func(repository *mockdb.UserMockRepository) {
				repository.EXPECT().Register(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody registerRequestBody) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "without password",
			requestBody: registerRequestBody{
				Username: faker.Username(),
			},
			buildStubs: func(repository *mockdb.UserMockRepository) {
				repository.EXPECT().Register(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody registerRequestBody) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ok",
			requestBody: registerRequestBody{
				Username: faker.Username(),
				Password: faker.Password(),
			},
			buildStubs: func(repository *mockdb.UserMockRepository) {
				repository.EXPECT().Register(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody registerRequestBody) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mockdb.NewUserMockRepository(ctrl)
	app, err := application.New(c)
	require.NoError(t, err)

	app.SetUserRepository(repository)
	app.RegisterRoutes()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(repository)

			recorder := httptest.NewRecorder()
			jsonData, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)
			request := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonData))

			app.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder, tc.requestBody)
		})
	}
}
