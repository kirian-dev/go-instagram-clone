package utils

import "strconv"

func ParsePort(portStr string) uint16 {
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {

		return 0
	}
	return uint16(port)
}
