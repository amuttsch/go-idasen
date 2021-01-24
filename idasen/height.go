package idasen

func (i *Idasen) HeightInMeters() (float64, error) {
	height, err := i.device.GetCharByUUID(_UUID_HEIGHT)
	if err != nil {
		log.Errorf("Cannot get height: %s", err)
		return 0, err
	}

	raw, err := height.ReadValue(getOptions())
	if err != nil {
		log.Errorf("Cannot ReadValue: %s", err)
		return 0, err
	}

	meters := heightBytesToMeter(raw)

	select {
	case i.HeightCh <- meters:
	default:
	}

	return meters, nil
}
