package idasen

import (
	"errors"
	"fmt"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"os"
	"os/signal"
	"regexp"
	"time"
)

const _DISCOVER_TIMEOUT = 5 * time.Second

func DiscoverDesk() (*desk, error) {
	//clean up connection on exit
	defer api.Exit()

	a, err := adapter.GetDefaultAdapter()
	if err != nil {
		return nil, err
	}

	log.Debug("Start discovery")
	discovery, cancel, err := api.Discover(a, nil)
	if err != nil {
		return nil, err
	}
	defer cancel()

	discoverTimeout := time.NewTimer(_DISCOVER_TIMEOUT)
	deskChan := make(chan *desk)

	deskNameRegex, err := regexp.Compile("Desk \\d+")
	if err != nil {
		log.Errorf("invalid regex: %s", err)
		return nil, errors.New("invalid regex")
	}

	go func() {

		for ev := range discovery {

			if ev.Type == adapter.DeviceRemoved {
				continue
			}

			dev, err := device.NewDevice1(ev.Path)
			if err != nil {
				log.Errorf("%s: %s", ev.Path, err)
				continue
			}

			if dev == nil {
				log.Errorf("%s: not found", ev.Path)
				continue
			}

			log.Debugf("name=%s addr=%s rssi=%d\n", dev.Properties.Name, dev.Properties.Address, dev.Properties.RSSI)

			if deskNameRegex.MatchString(dev.Properties.Name) {
				dev.Pair()
				log.Debugf("Paired device %s", dev.Properties.Name)

				dev.SetTrusted(true)
				log.Debugf("Trusting device %s", dev.Properties.Name)

				deskChan <- &desk{
					Name:    dev.Properties.Name,
					Address: dev.Properties.Address,
				}

				return
			}
		}

	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)

	select {
	case desk := <-deskChan:
		return desk, nil
	case <-discoverTimeout.C:
		e := fmt.Errorf("Discover timeout reached after %s", _DISCOVER_TIMEOUT.String())
		log.Debugln(e)
		return nil, e
	case sig := <-ch:
		e := fmt.Errorf("Received signal [%v]; shutting down...", sig)
		log.Debugln(e)
		return nil, e
	}
}
