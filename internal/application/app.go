package application

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohammad19khodaei/restaurant_reservation/config"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/domains/user"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/repositories"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/services/token"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Application struct
type Application struct {
	Config       *config.Config
	Router       *gin.Engine
	DB           *gorm.DB
	Repositories struct {
		UserRepository user.Repository
	}
	Services struct {
		TokenManger token.Manager
	}
}

// New creates a new Application
func New(config *config.Config) (*Application, error) {
	app := &Application{Config: config}
	err := app.registerDatabase()
	if err != nil {
		return nil, err
	}
	app.registerRepositories()
	app.registerServices()
	app.registerRouter()
	return app, nil
}

// Run the application
func (a *Application) Run(ctx context.Context) {
	srv := &http.Server{
		Addr:    a.Config.App.Address,
		Handler: a.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	shutdownCTX, cancel := context.WithTimeout(context.Background(), a.Config.App.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCTX); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}
}

// SetUserRepository sets the user repository for testing
func (a *Application) SetUserRepository(repository user.Repository) {
	a.Repositories.UserRepository = repository
}

func (a *Application) registerRouter() {
	router := gin.New()
	router.Use(gin.Logger())

	a.Router = router
}

func (a *Application) registerDatabase() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", a.Config.Database.Host, a.Config.Database.Username, a.Config.Database.Password, a.Config.Database.Name, a.Config.Database.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	a.DB = db
	return nil
}

func (a *Application) registerRepositories() {
	if a.Config.App.TestingMode {
		return
	}
	a.Repositories.UserRepository = repositories.NewGormUserRepository(a.DB)
}

func (a *Application) registerServices() {
	tokenManager, err := token.NewJWTManger(a.Config.App.SecretKey)
	if err != nil {
		log.Fatalf("could not create token manager: %v", err)
	}

	a.Services.TokenManger = tokenManager
}
