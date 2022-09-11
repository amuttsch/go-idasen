package idasen

import (
	"github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"
)

var log = logrus.New()

type Idasen struct {
	device *bluetooth.Device
}

type Configuration struct {
	MacAddress string             `mapstructure:"mac_address"`
	Positions  map[string]float64 `mapstructure:"positions"`
}

type desk struct {
	Name    string
	Address string
}

func New(config Configuration) (*Idasen, error) {
	desk, err := getDesk(config.MacAddress)
	if err != nil {
		log.Errorf("Could not find device %s: %s", config.MacAddress, err)
		return nil, err
	}

	d, err := adapter.Connect(desk.Address, bluetooth.ConnectionParams{})
	if err != nil {
		log.Errorf("Cannot connect to device %s: %s", config.MacAddress, err)
		return nil, err
	}

	log.Debugf("Connected to desk %s", config.MacAddress)

	return &Idasen{
		device: d,
	}, nil
}

func (i *Idasen) readValue(uuid string) ([]byte, error) {
	srvcs, _ := i.device.DiscoverServices(nil)

	for _, srvc := range srvcs {
		chars, err := srvc.DiscoverCharacteristics(nil)
		if err != nil {
			println(err)
		}
		for _, char := range chars {
			if char.UUID().String() == uuid {
				raw := make([]byte, 255)
				_, err = char.Read(raw)

				return raw, err
			}
		}
	}

	return nil, nil
}

func (i *Idasen) writeValue(uuid string, value []byte) (int, error) {
	srvcs, _ := i.device.DiscoverServices(nil)

	for _, srvc := range srvcs {
		chars, err := srvc.DiscoverCharacteristics(nil)
		if err != nil {
			println(err)
		}
		for _, char := range chars {
			if char.UUID().String() == uuid {
				return char.WriteWithoutResponse(value)
			}
		}
	}

	return -1, nil
}

func (i *Idasen) Close() {
	defer i.device.Disconnect()
}

func SetDebug() {
	log.SetLevel(logrus.DebugLevel)
}
