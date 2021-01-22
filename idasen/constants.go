package idasen

import "encoding/binary"

var (
	// Minimum desk height in meters.
	_MIN_HEIGHT = 0.62

	// Maximum desk height in meters.
	_MAX_HEIGHT = 1.27

	_UUID_HEIGHT = "99fa0021-338a-1024-8a49-009c0215f78a"
	_UUID_COMMAND = "99fa0002-338a-1024-8a49-009c0215f78a"
	_UUID_REFERENCE_INPUT = "99fa0031-338a-1024-8a49-009c0215f78a"

    _COMMAND_REFERENCE_INPUT_STOP = []byte{0x01, 0x80}
	_COMMAND_UP = []byte{0x47, 0x00}
	_COMMAND_DOWN = []byte{0x46, 0x00}
	_COMMAND_STOP = []byte{0xFF, 0x00}
)

func getOptions() map[string]interface{} {
	options := make(map[string]interface{})
	return options
}

func heightBytesToMeter(b []byte) float64 {
	if len(b) < 2 {
		return -1
	}

	raw := binary.LittleEndian.Uint16(b[0:2])
	return float64(float64(raw) / 10000) + _MIN_HEIGHT
}
