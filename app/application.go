package app

import (
	"github.com/JingdaMai/bookstore_users-api/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	logger.Info("about to start application...")
	_ = router.Run(":8081")
}
