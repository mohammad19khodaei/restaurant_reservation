package actions_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	mockdb "github.com/mohammad19khodaei/restaurant_reservation/db/mock"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/api/actions"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/application"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/domains/user"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestLoginAction(t *testing.T) {
	password := faker.Password()
	u := createRandomUser(password)

	testCases := []struct {
		name          string
		requestBody   loginRequestBody
		buildStubs    func(repository *mockdb.UserMockRepository, requestBody loginRequestBody)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody loginRequestBody)
	}{
		{
			name: "username not found",
			requestBody: loginRequestBody{
				Username: faker.Username(),
				Password: faker.Password(),
			},
			buildStubs: func(repository *mockdb.UserMockRepository, requestBody loginRequestBody) {
				repository.EXPECT().FindByUsername(gomock.Any(), gomock.Eq(requestBody.Username)).
					Times(1).
					Return(nil, user.ErrUserNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody loginRequestBody) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "invalid password",
			requestBody: loginRequestBody{
				Username: u.Username,
				Password: faker.Password(),
			},
			buildStubs: func(repository *mockdb.UserMockRepository, requestBody loginRequestBody) {
				repository.EXPECT().FindByUsername(gomock.Any(), gomock.Eq(requestBody.Username)).
					Times(1).
					Return(&u, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody loginRequestBody) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ok",
			requestBody: loginRequestBody{
				Username: u.Username,
				Password: password,
			},
			buildStubs: func(repository *mockdb.UserMockRepository, requestBody loginRequestBody) {
				repository.EXPECT().FindByUsername(gomock.Any(), gomock.Eq(requestBody.Username)).
					Times(1).
					Return(&u, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody loginRequestBody) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var resp actions.LoginResponse
				err := json.NewDecoder(recorder.Body).Decode(&resp)
				require.NoError(t, err)

				require.NotEmpty(t, resp.AccessToken)
				require.Equal(t, u.Username, resp.User.Username)
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
			tc.buildStubs(repository, tc.requestBody)
			recorder := httptest.NewRecorder()
			requestBody := tc.requestBody

			jsonData, err := json.Marshal(requestBody)
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(jsonData))
			app.Router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder, requestBody)
		})
	}
}

func createRandomUser(password string) user.User {
	hashedPassword, _ := utils.HashPassword(password)
	return user.User{
		Username:  faker.Username(),
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}
}
