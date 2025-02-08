package actions_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/mohammad19khodaei/restaurant_reservation/db/mock"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/api/middlewares"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/application"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/domains/reservation"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/services/token"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCancelAction(t *testing.T) {
	userID := 1
	testCases := []struct {
		name          string
		requestBody   cancelRequest
		setAuthHeader func(t *testing.T, manager token.Manager, req *http.Request)
		buildStubs    func(repository *mockdb.ReservationMockRepository, requestBody cancelRequest)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "reservation not found",
			requestBody: cancelRequest{1},
			setAuthHeader: func(t *testing.T, manager token.Manager, req *http.Request) {
				token, err := manager.GenerateToken(userID, c.App.TokenDuration)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, token))
			},
			buildStubs: func(repository *mockdb.ReservationMockRepository, requestBody cancelRequest) {
				repository.EXPECT().
					CancelReservation(gomock.Any(), requestBody.ID).
					Return(reservation.ErrReservationNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:        "ok",
			requestBody: cancelRequest{1},
			setAuthHeader: func(t *testing.T, manager token.Manager, req *http.Request) {
				token, err := manager.GenerateToken(userID, c.App.TokenDuration)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, token))
			},
			buildStubs: func(repository *mockdb.ReservationMockRepository, requestBody cancelRequest) {
				repository.EXPECT().
					CancelReservation(gomock.Any(), requestBody.ID).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	tokenManager, err := token.NewJWTManger(c.App.SecretKey)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mockdb.NewReservationMockRepository(ctrl)
	app, err := application.New(c)
	require.NoError(t, err)
	app.SetReservationRepository(repository)
	app.RegisterRoutes()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(repository, tc.requestBody)

			recorder := httptest.NewRecorder()
			jsonData, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)
			request := httptest.NewRequest("POST", "/cancel", bytes.NewReader(jsonData))
			tc.setAuthHeader(t, tokenManager, request)

			app.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

type cancelRequest struct {
	ID int `json:"id"`
}
