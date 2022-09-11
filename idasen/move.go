package idasen

import "fmt"

func (i *Idasen) MoveUp() error {
	_, err := i.writeValue(_UUID_COMMAND, _COMMAND_UP)
	if err != nil {
		log.Errorf("Cannot move desk: %s", err)
		return err
	}

	return nil
}

func (i *Idasen) MoveDown() error {
	_, err := i.writeValue(_UUID_COMMAND, _COMMAND_DOWN)
	if err != nil {
		log.Errorf("Cannot move desk: %s", err)
		return err
	}

	return nil
}

func (i *Idasen) MoveStop() error {
	_, err := i.writeValue(_UUID_COMMAND, _COMMAND_STOP)
	if err != nil {
		log.Errorf("Cannot stop desk: %s", err)
		return err
	}

	_, err = i.writeValue(_UUID_REFERENCE_INPUT, _COMMAND_STOP)
	if err != nil {
		log.Errorf("Cannot stop desk: %s", err)
		return err
	}

	return nil
}

func (i *Idasen) MoveToTarget(targetInMeters float64) error {
	if targetInMeters > _MAX_HEIGHT {
		return fmt.Errorf("target heigth %.4fm exceeds maximum height of %.2fm", targetInMeters, _MAX_HEIGHT)
	}

	if targetInMeters < _MIN_HEIGHT {
		return fmt.Errorf("target heigth %.4fm exceeds minimum height of %.2fm", targetInMeters, _MIN_HEIGHT)
	}

	previousHeight, err := i.HeightInMeters()
	if err != nil {
		return err
	}

	willMoveUp := previousHeight < targetInMeters

	log.Debugf("Will move to target %.4f\n", targetInMeters)

	for true {
		height, err := i.HeightInMeters()
		if err != nil {
			return err
		}

		// Check if the safety stop was triggered
		if (height > previousHeight && !willMoveUp) || (height < previousHeight && willMoveUp) {
			_ = i.MoveStop()
			return fmt.Errorf("safety stop was trigged")
		}

		difference := 0.0
		if willMoveUp {
			difference = targetInMeters - height
		} else {
			difference = height - targetInMeters
		}

		log.Debugf("Moving desk. Current height %.4fm, difference %.4f\n", height, difference)

		if difference < 0.005 {
			log.Debugf("Reached target height")
			return i.MoveStop()
		}

		if willMoveUp {
			err = i.MoveUp()
		} else {
			err = i.MoveDown()
		}
		if err != nil {
			return err
		}

		previousHeight = height
	}

	return nil
}
