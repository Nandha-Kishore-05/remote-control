

package main

import (
	"wakeonlan/config"
	turnoff "wakeonlan/controllers/TurnOffPc"
	turnon "wakeonlan/controllers/TurnOnPc"
	controllers "wakeonlan/controllers/safeExam"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	router := gin.Default()

	router.GET("/poweron", turnon.PowerOnHandler)
	router.POST("/shutdown", turnoff.ShutdownHandler)
	router.GET("/launch-seb", controllers.LaunchSEBHandler)
	router.GET("/exit-seb", controllers.ExitSEBHandler)
	router.Run(":8080")
}
