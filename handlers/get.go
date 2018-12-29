package handlers

import (
	"net/http"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/tcl-control-microservice/helpers"
	"github.com/labstack/echo"
)

// GetPower sends a command to get the power state of the TV
func GetPower(context echo.Context) error {
	address := context.Param("address")

	devInfo, err := helpers.GetDeviceInfo(address)
	if err != nil {
		log.L.Error(err.String())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, status.Power{Power: helpers.PowerStateMap[devInfo.PowerMode]})
}

// GetInput sends a command to get the current input of the TV
func GetInput(context echo.Context) error {
	address := context.Param("address")

	input, err := helpers.GetCurrentInput(address)
	if err != nil {
		log.L.Error(err.String())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, input)
}

// GetInputList sends a command to get the list of inputs for the TV
func GetInputList(context echo.Context) error {
	address := context.Param("address")

	list, err := helpers.GetInputList(address)
	if err != nil {
		log.L.Error(err.String())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, list)
}

// GetActiveSignal sends a command to determine if there is an active signal on the given port
func GetActiveSignal(context echo.Context) error {
	address := context.Param("address")
	port := context.Param("port")

	currentInput, err := helpers.GetCurrentInput(address)
	if err != nil {
		log.L.Error(err.String())
		return context.JSON(http.StatusInternalServerError, err)
	}

	var active structs.ActiveSignal

	if strings.EqualFold(port, currentInput.Input) {
		active.Active = true
	} else {
		active.Active = false
	}

	return context.JSON(http.StatusOK, active)
}

// GetVolume sends a command to get the volume level of the TV
// Roku TVs unfortunately cannot do this yet...
func GetVolume(context echo.Context) error {
	return context.JSON(http.StatusOK, "We unfortunately can't do this :(")
}

// GetMute sends a command to see if the TV is muted or not
// Roku TVs unfortunately cannot do this yet...
func GetMute(context echo.Context) error {
	return context.JSON(http.StatusOK, "We unfortunately can't do this :(")
}

// GetBlank sends a command to see if the TV is blanked or not
// There might be a dirty way to do this...
func GetBlank(context echo.Context) error {
	return context.JSON(http.StatusOK, "I'm hesitant on this one, but not hopeful :/")
}

// GetHardwareInfo sends a command to get the hardware information of the TV
func GetHardwareInfo(context echo.Context) error {
	address := context.Param("address")

	hardware, err := helpers.GetHardwareInfo(address)
	if err != nil {
		log.L.Error(err.String())
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, hardware)
}
