package idasen

import (
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Idasen struct {
	device *device.Device1
}

type Configuration struct {
	MacAddress string `mapstructure:"mac_address"`
	Positions map[string]float64 `mapstructure:"positions"`
}

type desk struct {
	Name    string
	Address string
}

func New(config Configuration) (*Idasen, error) {

	a, err := adapter.GetDefaultAdapter()
	if err != nil {
		api.Exit()
		return nil, err
	}

	d, err := a.GetDeviceByAddress(config.MacAddress)
	if err != nil || d == nil {
		log.Errorf("Could not find device %s: %s", config.MacAddress, err)
		api.Exit()
		return nil, err
	}

	err = d.Connect()
	if err != nil {
		log.Errorf("Cannot connect to device %s: %s", config.MacAddress, err)
		api.Exit()
		return nil, err
	}

	log.Debugf("Connected to desk %s", config.MacAddress)

	err = d.Pair()
	if err != nil {
		log.Errorf("Cannot connect to device %s: %s", config.MacAddress, err)
		api.Exit()
		return nil, err
	}

	log.Debugf("Paired with desk %s", config.MacAddress)

	return &Idasen{
		device: d,
	}, nil
}

func (i *Idasen) Close() {
	defer api.Exit()
	defer i.device.Disconnect()
}

func SetDebug() {
	log.SetLevel(logrus.DebugLevel)
}
