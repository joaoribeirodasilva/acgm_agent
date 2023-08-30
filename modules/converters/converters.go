package converters

import (
	"strconv"
	"time"
)

func StringToInterval(interval string) (*time.Duration, error) {

	milliseconds, err := strconv.ParseInt(interval, 10, 64)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(milliseconds) * time.Millisecond
	return &duration, nil
}
