package idasen

import (
	"fmt"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Idasen struct {
	device *device.Device1
	HeightCh chan float64
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

	if d.Properties.Connected {
		return nil, fmt.Errorf("Device already connected, cannot use go-idasen concurrently!")
	}

	err = d.Connect()
	if err != nil {
		log.Errorf("Cannot connect to device %s: %s", config.MacAddress, err)
		api.Exit()
		return nil, err
	}

	log.Debugf("Connected to desk %s", config.MacAddress)

	if !d.Properties.Paired {
		err = d.Pair()
		if err != nil {
			log.Errorf("Cannot connect to device %s: %s", config.MacAddress, err)
			api.Exit()
			return nil, err
		}

		log.Debugf("Paired with desk %s", config.MacAddress)
	}

	uuids, err := d.GetUUIDs()
	if err != nil {
		log.Errorf("Cannot detect uuids: %s", err)
		api.Exit()
		return nil, err
	}
	log.Trace("uuids: %v, err: %s\n", uuids, err)

	return &Idasen{
		device: d,
		HeightCh: make(chan float64),
	}, nil
}

func (i *Idasen) Close() {
	defer api.Exit()
	defer i.Disconnect()

	log.Debugln("Closing connection and exit api.")
}

func  (i *Idasen) Disconnect() {
	close(i.HeightCh)
	err := i.device.Disconnect()

	log.Debugf("Disconnected. Error: %s", err)
}

func SetDebug() {
	log.SetLevel(logrus.DebugLevel)
}
