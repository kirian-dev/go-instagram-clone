package utils

import "strconv"

func ParsePort(portStr string) int {
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {

		return 0
	}
	return int(port)
}
