package helpers

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/byuoitav/common/status"
	"github.com/byuoitav/common/structs"

	"github.com/byuoitav/common/nerr"
)

// A list of key values
const (
	Home        = "Home"
	Rewind      = "Rev"
	Forward     = "Fwd"
	Play        = "Play"
	Pause       = "Play"
	Select      = "Select"
	Left        = "Left"
	Right       = "Right"
	Down        = "Down"
	Up          = "Up"
	Back        = "Back"
	Replay      = "InstantReplay"
	Star        = "Info"
	Backspace   = "Backspace"
	Search      = "Search"
	Enter       = "Enter"
	FindRemote  = "FindRemote"
	VolumeDown  = "VolumeDown"
	Mute        = "VolumeMute"
	VolumeUp    = "VolumeUp"
	PowerOn     = "PowerOn"
	Standby     = "PowerOff"
	ChannelUp   = "ChannelUp"
	ChannelDown = "ChannelDown"
	Tuner       = "InputTuner"
	HDMI1       = "InputHDMI1"
	HDMI2       = "InputHDMI2"
	HDMI3       = "InputHDMI3"
	HDMI4       = "InputHDMI4"
	Component   = "InputAV1"
)

// A list of query values
const (
	DeviceInfo   = "device-info"
	AllApps      = "apps"
	CurrentInput = "active-app"
)

// A list of Roku specific terms
const (
	TVInput      = "tvin"
	Application  = "appl"
	Menu         = "menu"
	NDKA         = "ndka"
	RSGA         = "rsga"
	InputWhenOff = "Davinci Channel"
)

// InputMap is a mapping of input names
var InputMap = map[string]string{
	"hdmi1":      HDMI1,
	"hdmi2":      HDMI2,
	"hdmi3":      HDMI3,
	"hdmi4":      HDMI4,
	"av":         Component,
	"tuner":      Tuner,
	InputWhenOff: "blanked",
}

// PowerStateMap maps their power state names to the ones used by this API
var PowerStateMap = map[string]string{
	"PowerOn":  "on",
	"Headless": "standby",
}

var client http.Client

// SendKeyPressRequest sends a command to press a key
func SendKeyPressRequest(address, key string) *nerr.E {
	url := fmt.Sprintf("http://%s:8060/keypress/%s", address, key)

	_, err := client.Post(url, "", bytes.NewReader([]byte{}))
	if err != nil {
		return nerr.Translate(err).Addf("failed to send %s keypress request to %s", key, address)
	}

	return nil
}

// SendKeyDownRequest sends a command to push a key down
func SendKeyDownRequest(address, key string) *nerr.E {
	url := fmt.Sprintf("http://%s:8060/keydown/%s", address, key)

	_, err := client.Post(url, "", bytes.NewReader([]byte{}))
	if err != nil {
		return nerr.Translate(err).Addf("failed to send %s keydown request to %s", key, address)
	}

	return nil
}

// SendKeyUpRequest sends a command to "release" a key down
func SendKeyUpRequest(address, key string) *nerr.E {
	url := fmt.Sprintf("http://%s:8060/keyup/%s", address, key)

	_, err := client.Post(url, "", bytes.NewReader([]byte{}))
	if err != nil {
		return nerr.Translate(err).Addf("failed to send %s keyup request to %s", key, address)
	}

	return nil
}

// sendQueryRequest sends a query request
func sendQueryRequest(address, query string) ([]byte, *nerr.E) {
	url := fmt.Sprintf("http://%s:8060/query/%s", address, query)

	resp, err := client.Get(url)
	if err != nil {
		return nil, nerr.Translate(err).Addf("failed to send %s query to %s", query, address)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, nerr.Translate(err).Addf("failed to read response body from %s", address)
	}

	return body, nil
}

// GetDeviceInfo returns the device information
func GetDeviceInfo(address string) (RokuDeviceResponse, *nerr.E) {
	var toReturn RokuDeviceResponse

	byteArray, ne := sendQueryRequest(address, DeviceInfo)
	if ne != nil {
		return toReturn, ne.Addf("failed to get device info from %s", address)
	}

	err := xml.Unmarshal(byteArray, &toReturn)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to unmarshal XML response from %s", address)
	}

	return toReturn, nil
}

// GetCurrentInput returns the active input information
func GetCurrentInput(address string) (status.Input, *nerr.E) {
	var currentApp RokuActiveApp

	body, ne := sendQueryRequest(address, CurrentInput)
	if ne != nil {
		return status.Input{}, ne.Addf("failed to get current input from %s", address)
	}

	err := xml.Unmarshal(body, &currentApp)
	if err != nil {
		return status.Input{}, nerr.Translate(err).Addf("failed to unmarshal XML response from %s", address)
	}

	currentInput := status.Input{
		Input: getInputReturnNameBasedOnType(currentApp.App),
	}

	return currentInput, nil
}

// GetInputList returns the list off all possible inputs and/or applications
func GetInputList(address string) ([]string, *nerr.E) {
	toReturn := []string{}
	var appList RokuAppList

	body, ne := sendQueryRequest(address, AllApps)
	if ne != nil {
		return toReturn, ne.Addf("failed to get the list of inputs from %s", address)
	}

	err := xml.Unmarshal(body, &appList)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to unmarshal XMl response from %s", address)
	}

	for _, app := range appList.Apps {
		toReturn = append(toReturn, getInputReturnNameBasedOnType(app))
	}

	return toReturn, nil
}

func getInputReturnNameBasedOnType(app RokuApp) string {
	var finalInput string

	switch app.Type {
	case TVInput:
		input := strings.Split(app.ID, ".")[1]
		finalInput = strings.ToUpper(input)
	case Application:
		if app.Text == InputWhenOff {
			finalInput = InputMap[InputWhenOff]
		} else {
			finalInput = app.Text
		}
	case Menu:
		finalInput = "Menu"
	default:
		finalInput = app.Text
	}

	return finalInput
}

// GetHardwareInfo takes all the device info and packages it into the HardwareInfo struct format
func GetHardwareInfo(address string) (structs.HardwareInfo, *nerr.E) {
	var hardware structs.HardwareInfo

	deviceInfo, ne := GetDeviceInfo(address)
	if ne != nil {
		return hardware, ne.Addf("failed to get the hardware info for %s", address)
	}

	// get the hostname
	addr, e := net.LookupAddr(address)
	if e != nil {
		hardware.Hostname = address
	} else {
		hardware.Hostname = strings.Trim(addr[0], ".")
	}

	hardware.ModelName = deviceInfo.ModelName
	hardware.SerialNumber = deviceInfo.SerialNumber
	hardware.FirmwareVersion = deviceInfo.SoftwareVersion

	hardware.NetworkInfo = structs.NetworkInfo{
		IPAddress:  address,
		MACAddress: deviceInfo.EthernetMAC,
	}

	hardware.PowerStatus = PowerStateMap[deviceInfo.PowerMode]
	return hardware, nil
}

// RokuDeviceResponse is the struct that the Roku device sends back when asking for device info
type RokuDeviceResponse struct {
	XMLName                     xml.Name `xml:"device-info"`
	UDN                         string   `xml:"udn"`
	SerialNumber                string   `xml:"serial-number"`
	DeviceID                    string   `xml:"device-id"`
	AdvertisingID               string   `xml:"advertising-id"`
	VendorName                  string   `xml:"vendor-name"`
	ModelName                   string   `xml:"model-name"`
	ModelNumber                 string   `xml:"model-number"`
	ModelRegion                 string   `xml:"model-region"`
	IsTV                        bool     `xml:"is-tv"`
	IsStick                     bool     `xml:"is-stick"`
	ScreenSize                  int      `xml:"screen-size"`
	PanelID                     int      `xml:"panel-id"`
	TunerType                   string   `xml:"tuner-type"`
	SupportsEthernet            bool     `xml:"supports-ethernet"`
	WiFiMAC                     string   `xml:"wifi-mac"`
	WiFiDriver                  string   `xml:"wifi-driver"`
	EthernetMAC                 string   `xml:"ethernet-mac"`
	NetworkType                 string   `xml:"network-type"`
	NetworkName                 string   `xml:"network-name"`
	FriendlyDeviceName          string   `xml:"friendly-device-name"`
	FriendlyModelName           string   `xml:"friendly-model-name"`
	DefaultDeviceName           string   `xml:"default-device-name"`
	UserDeviceName              string   `xml:"user-device-name"`
	SoftwareVersion             string   `xml:"software-version"`
	SoftwareBuild               string   `xml:"software-build"`
	SecureDevice                bool     `xml:"secure-device"`
	Language                    string   `xml:"language"`
	Country                     string   `xml:"country"`
	Locale                      string   `xml:"locale"`
	TimeZoneAuto                bool     `xml:"time-zone-auto"`
	TimeZone                    string   `xml:"time-zone"`
	TimeZoneName                string   `xml:"time-zone-name"`
	TimeZoneTZ                  string   `xml:"time-zone-tz"`
	TimeZoneOffset              int      `xml:"time-zone-offset"`
	ClockFormat                 string   `xml:"clock-format"`
	Uptime                      int      `xml:"uptime"`
	PowerMode                   string   `xml:"power-mode"`
	SupportsSuspend             bool     `xml:"supports-suspend"`
	SupportsFindRemote          bool     `xml:"supports-find-remote"`
	SupportsAudioGuide          bool     `xml:"supports-audio-guide"`
	SupportsRVA                 bool     `xml:"supports-rva"`
	DeveloperEnabled            bool     `xml:"developer-enabled"`
	KeyedDeveloperID            string   `xml:"keyed-developer-id"`
	SearchEnabled               bool     `xml:"search-enabled"`
	SearchChannelsEnabled       bool     `xml:"search-channels-enabled"`
	VoiceSearchEnabled          bool     `xml:"voice-search-enabled"`
	NotificationsEnabled        bool     `xml:"notifications-enabled"`
	NotificationsFirstUse       bool     `xml:"notifications-first-use"`
	SupportsPrivateListening    bool     `xml:"supports-private-listening"`
	SupportsPrivateListeningDTV bool     `xml:"supports-private-listening-dtv"`
	SupportsWarmStandby         bool     `xml:"supports-warm-standby"`
	HeadphonesConnected         bool     `xml:"headphones-connected"`
	ExpertPQEnabled             string   `xml:"expert-pq-enabled"`
	SupportsECSTextEdit         bool     `xml:"supports-ecs-textedit"`
	SupportsECSMicrophone       bool     `xml:"supports-ecs-microphone"`
	SupportsWakeOnWLAN          bool     `xml:"supports-wake-on-wlan"`
	HasPlayOnRoku               bool     `xml:"has-play-on-roku"`
	HasMobileScreensaver        bool     `xml:"has-mobile-screensaver"`
	SupportURL                  string   `xml:"support-url"`
}

// RokuAppList is a representation of how the Roku lists its possible apps
type RokuAppList struct {
	XMLName xml.Name  `xml:"apps"`
	Text    string    `xml:",chardata"`
	Apps    []RokuApp `xml:"app"`
}

// RokuApp is the representation for one single app/input port
type RokuApp struct {
	Text    string `xml:",chardata"`
	ID      string `xml:"id,attr"`
	Subtype string `xml:"subtype,attr"`
	Type    string `xml:"type,attr"`
	Version string `xml:"version,attr"`
}

// RokuActiveApp is a representation of the information presented to show the current active app/input
type RokuActiveApp struct {
	XMLName xml.Name `xml:"active-app"`
	Text    string   `xml:",chardata"`
	App     RokuApp  `xml:"app"`
}
