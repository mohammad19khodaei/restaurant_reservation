package actions_test

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mohammad19khodaei/restaurant_reservation/config"
)

var (
	c *config.Config
)

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig("../../../config", "config-testing")
	if err != nil {
		log.Fatal("could not load config", err)
	}
	c = cfg

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
