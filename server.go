package main

import (
	"net/http"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/auth"
	"github.com/byuoitav/tcl-control-microservice/handlers"
)

func main() {
	port := ":8023"

	router := common.NewRouter()

	// Log Endpoints
	router.GET("/log-level", log.GetLogLevel)
	router.PUT("/log-level/:level", log.SetLogLevel)

	// functionality endpoints
	write := router.Group("", auth.AuthorizeRequest("write-state", "room", auth.LookupResourceFromAddress))
	write.GET("/:address/power/on", handlers.PowerOn)
	write.GET("/:address/power/standby", handlers.Standby)
	write.GET("/:address/input/:port", handlers.SwitchInput)
	write.GET("/:address/volume/set/:value", handlers.SetVolume)
	write.GET("/:address/volume/up", handlers.VolumeUp)
	write.GET("/:address/volume/down", handlers.VolumeDown)
	write.GET("/:address/volume/mute", handlers.Mute)
	write.GET("/:address/volume/unmute", handlers.Unmute)
	write.GET("/:address/display/blank", handlers.BlankDisplay)
	write.GET("/:address/display/unblank", handlers.UnblankDisplay)

	// status endpoints
	read := router.Group("", auth.AuthorizeRequest("read-state", "room", auth.LookupResourceFromAddress))
	read.GET("/:address/power/status", handlers.GetPower)
	read.GET("/:address/input/current", handlers.GetInput)
	read.GET("/:address/input/list", handlers.GetInputList)
	read.GET("/:address/active/:port", handlers.GetActiveSignal)
	read.GET("/:address/volume/level", handlers.GetVolume)
	read.GET("/:address/volume/mute/status", handlers.GetMute)
	read.GET("/:address/display/status", handlers.GetBlank)
	read.GET("/:address/hardware", handlers.GetHardwareInfo)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
