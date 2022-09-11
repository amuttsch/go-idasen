package idasen

func (i *Idasen) HeightInMeters() (float64, error) {
	raw, err := i.readValue(_UUID_HEIGHT)
	if err != nil {
		log.Errorf("Cannot get height: %s", err)
		return 0, err
	}

	return heightBytesToMeter(raw), nil
}
