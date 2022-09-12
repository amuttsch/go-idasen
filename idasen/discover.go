package idasen

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"time"
	"tinygo.org/x/bluetooth"
)

func DiscoverDesk(timeout int64) (*desk, error) {
	result, err := getDesk("", timeout)
	if err != nil {
		return nil, err
	}

	return &desk{
		Name:    result.LocalName(),
		Address: result.Address.String(),
	}, nil
}

func getDesk(mac string, timeout int64) (*bluetooth.ScanResult, error) {
	log.Debug("Start discovery")
	err := adapter.Enable()
	if err != nil {
		return nil, fmt.Errorf("must enable adapter: %w", err)
	}

	resultCh := make(chan bluetooth.ScanResult, 1)

	deskNameRegex, err := regexp.Compile("Desk \\d+")
	if err != nil {
		log.Errorf("invalid regex: %s", err)
		return nil, errors.New("invalid regex")
	}

	// Start scanning.
	log.Debugln("Scanning...")
	err = adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if mac == "" {
			if deskNameRegex.MatchString(result.LocalName()) {
				log.Infoln("Found Desk:", result.Address.String(), result.RSSI, result.LocalName())
				adapter.StopScan()
				resultCh <- result
			}
		} else {
			if result.Address.String() == mac {
				adapter.StopScan()
				resultCh <- result
			}
		}

	})

	discoverTimeout := time.NewTimer(time.Duration(timeout) * time.Second)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)

	select {
	case result := <-resultCh:
		return &result, nil
	case <-discoverTimeout.C:
		e := fmt.Errorf("Discover timeout reached after %d seconds", timeout)
		log.Debugln(e)
		return nil, e
	case sig := <-ch:
		e := fmt.Errorf("Received signal [%v]; shutting down...", sig)
		log.Debugln(e)
		return nil, e
	}
}
