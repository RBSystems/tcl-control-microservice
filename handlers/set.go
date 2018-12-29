package handlers

import (
	"net/http"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/byuoitav/tcl-control-microservice/helpers"
	"github.com/labstack/echo"
)

// PowerOn sends a command to turn on the TV
func PowerOn(context echo.Context) error {
	address := context.Param("address")

	err := helpers.SendKeyPressRequest(address, helpers.PowerOn)
	if err != nil {
		log.L.Error(err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, status.Power{Power: "on"})
}

// Standby sends a command to turn off the TV, or put it on standby
func Standby(context echo.Context) error {
	address := context.Param("address")

	err := helpers.SendKeyPressRequest(address, helpers.Standby)
	if err != nil {
		log.L.Error(err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, status.Power{Power: "standby"})
}

// SwitchInput sends a command to switch the input of the TV
func SwitchInput(context echo.Context) error {
	address := context.Param("address")
	port := context.Param("port")

	err := helpers.SendKeyPressRequest(address, helpers.InputMap[strings.ToLower(port)])
	if err != nil {
		log.L.Error(err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, status.Input{Input: strings.ToUpper(port)})
}

// SetVolume sends a command to set the volume of the TV
// Roku TVs unfortunately cannot do this yet...
func SetVolume(context echo.Context) error {
	return context.JSON(http.StatusOK, "We unfortunately can't do this :(")
}

// VolumeUp sends a command to move the volume up by one to the TV
func VolumeUp(context echo.Context) error {
	address := context.Param("address")

	err := helpers.SendKeyPressRequest(address, helpers.VolumeUp)
	if err != nil {
		log.L.Error(err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, "I don't know what the volume is...")
}

// VolumeDown sends a command to move the volume down by one to the TV
func VolumeDown(context echo.Context) error {
	address := context.Param("address")

	err := helpers.SendKeyPressRequest(address, helpers.VolumeDown)
	if err != nil {
		log.L.Error(err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, "I don't know what the volume is...")
}

// Mute sends a command to mute the TV
func Mute(context echo.Context) error {
	address := context.Param("address")

	// this is about to get real dirty
	err := helpers.SendKeyPressRequest(address, helpers.VolumeDown)
	if err != nil {
		log.L.Error(err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	err = helpers.SendKeyPressRequest(address, helpers.VolumeUp)
	if err != nil {
		log.L.Error(err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	err = helpers.SendKeyPressRequest(address, helpers.Mute)
	if err != nil {
		log.L.Error(err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, status.Mute{Muted: true})
}

// Unmute sends a command to unmute the TV
func Unmute(context echo.Context) error {
	address := context.Param("address")

	// this is about to get real dirty
	err := helpers.SendKeyPressRequest(address, helpers.VolumeDown)
	if err != nil {
		log.L.Error(err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	err = helpers.SendKeyPressRequest(address, helpers.VolumeUp)
	if err != nil {
		log.L.Error(err.Error())
		return context.JSON(http.StatusInternalServerError, err)
	}

	// err := helpers.SendKeyPressRequest(address, helpers.Mute)
	// if err != nil {
	// 	log.L.Error(err.Error())
	// 	return context.JSON(http.StatusInternalServerError, err)
	// }

	return context.JSON(http.StatusOK, status.Mute{Muted: false})
}

// BlankDisplay sends a command to blank the TV
// Roku TVs unfortunately cannot do this yet...
func BlankDisplay(context echo.Context) error {
	return context.JSON(http.StatusOK, "We unfortunately can't do this :(")
}

// UnblankDisplay sends a command to unblank the TV
// Roku TVs unfortunately cannot do this yet...
func UnblankDisplay(context echo.Context) error {
	return context.JSON(http.StatusOK, "We unfortunately can't do this :(")
}
