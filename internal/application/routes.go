package application

import (
	"github.com/mohammad19khodaei/restaurant_reservation/internal/api/actions"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/api/middlewares"
)

func (a *Application) RegisterRoutes() {
	a.Router.POST("users", actions.RegisterUserAction(a.Repositories.UserRepository))
	a.Router.POST("users/login", actions.LoginAction(a.Repositories.UserRepository, a.Services.TokenManger, a.Config.App.TokenDuration))

	authRoute := a.Router.Group("/").Use(middlewares.AuthMiddleware(a.Services.TokenManger))

	authRoute.POST("book", actions.BookAction(a.Repositories.ReservationRepository))
	authRoute.POST("cancel", actions.CancelAction(a.Repositories.ReservationRepository))
}
