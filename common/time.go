package common

import "time"

func GetCurrentTime() int64{
	t := time.Now()
	return t.Unix()
}

