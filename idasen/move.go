package idasen

import (
	"fmt"
)

func (i *Idasen) MoveUp() error {
	char, err := i.device.GetCharByUUID(_UUID_COMMAND)
	if err != nil {
		log.Errorf("Cannot get command char %s", err)
		return err
	}
	return char.WriteValue(_COMMAND_UP, getOptions())
}

func (i *Idasen) MoveDown() error {
	char, err := i.device.GetCharByUUID(_UUID_COMMAND)
	if err != nil {
		log.Errorf("Cannot get command char %s", err)
		return err
	}
	return char.WriteValue(_COMMAND_DOWN, getOptions())
}

func (i *Idasen) MoveStop() error {
	char_cmd, err := i.device.GetCharByUUID(_UUID_COMMAND)
	if err != nil {
		log.Errorf("Cannot get command char %s", err)
		return err
	}
	
	char_ref, err := i.device.GetCharByUUID(_UUID_REFERENCE_INPUT)
	if err != nil {
		log.Errorf("Cannot get reference input char %s", err)
		return err
	}

	_ = char_cmd.WriteValue(_COMMAND_STOP, getOptions())
	return char_ref.WriteValue(_COMMAND_REFERENCE_INPUT_STOP, getOptions())
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
