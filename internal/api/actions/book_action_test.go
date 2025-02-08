package actions_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/mohammad19khodaei/restaurant_reservation/db/mock"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/api/actions"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/api/middlewares"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/application"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/domains/reservation"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/services/token"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestBookAction(t *testing.T) {
	userID := 1
	seatPrice := 10
	testCases := []struct {
		name          string
		requestBody   bookRequest
		setAuthHeader func(t *testing.T, manager token.Manager, req *http.Request)
		buildStubs    func(repository *mockdb.ReservationMockRepository, requestBody bookRequest)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody bookRequest)
	}{
		{
			name: "booking a complete table",
			requestBody: bookRequest{
				SeatsCount: 4,
				Date:       time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			},
			setAuthHeader: func(t *testing.T, manager token.Manager, req *http.Request) {
				token, err := manager.GenerateToken(userID, c.App.TokenDuration)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, token))
			},
			buildStubs: func(repository *mockdb.ReservationMockRepository, requestBody bookRequest) {
				repository.EXPECT().
					BookTable(gomock.Any(), userID, requestBody.SeatsCount, gomock.Any()).
					Return(&reservation.Reservation{
						ID:         1,
						TableID:    1,
						SeatsCount: requestBody.SeatsCount,
						Price:      float64(seatPrice * (requestBody.SeatsCount - 1)),
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody bookRequest) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var resp actions.BookResponse
				err := json.NewDecoder(recorder.Body).Decode(&resp)
				require.NoError(t, err)

				require.NotEmpty(t, resp.ID)
				require.Equal(t, float64(seatPrice*(requestBody.SeatsCount-1)), resp.Price)
				require.Equal(t, requestBody.SeatsCount, resp.SeatsCount)
				require.NotEmpty(t, resp.TableID)
			},
		},
		{
			name: "booking a odd number of seats",
			requestBody: bookRequest{
				SeatsCount: 3,
				Date:       time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			},
			setAuthHeader: func(t *testing.T, manager token.Manager, req *http.Request) {
				token, err := manager.GenerateToken(userID, c.App.TokenDuration)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, token))
			},
			buildStubs: func(repository *mockdb.ReservationMockRepository, requestBody bookRequest) {
				repository.EXPECT().
					BookTable(gomock.Any(), userID, requestBody.SeatsCount+1, gomock.Any()).
					Return(&reservation.Reservation{
						ID:         1,
						TableID:    1,
						SeatsCount: requestBody.SeatsCount + 1,
						Price:      float64(seatPrice * (requestBody.SeatsCount + 1)),
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody bookRequest) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var resp actions.BookResponse
				err := json.NewDecoder(recorder.Body).Decode(&resp)
				require.NoError(t, err)

				require.NotEmpty(t, resp.ID)
				require.Equal(t, float64(seatPrice*(requestBody.SeatsCount+1)), resp.Price)
				require.Equal(t, requestBody.SeatsCount+1, resp.SeatsCount)
				require.NotEmpty(t, resp.TableID)
			},
		},
		{
			name: "no tables are available",
			requestBody: bookRequest{
				SeatsCount: 8,
				Date:       time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			},
			setAuthHeader: func(t *testing.T, manager token.Manager, req *http.Request) {
				token, err := manager.GenerateToken(userID, c.App.TokenDuration)
				require.NoError(t, err)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", middlewares.AuthorizationTypeBearer, token))
			},
			buildStubs: func(repository *mockdb.ReservationMockRepository, requestBody bookRequest) {
				repository.EXPECT().
					BookTable(gomock.Any(), userID, requestBody.SeatsCount, gomock.Any()).
					Return(nil, reservation.ErrNoTablesAreAvailable)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, requestBody bookRequest) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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
			request := httptest.NewRequest("POST", "/book", bytes.NewReader(jsonData))
			tc.setAuthHeader(t, tokenManager, request)

			app.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder, tc.requestBody)
		})
	}

}

type bookRequest struct {
	SeatsCount int    `json:"seats_count"`
	Date       string `json:"date"`
}
