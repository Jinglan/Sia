package modules

import (
	"time"

	"github.com/NebulousLabs/Sia/build"
)

const (
	NotifyBuffer = 3
)

var (
	SafeMutexDelay time.Duration
)

func init() {
	if build.Release == "dev" {
		SafeMutexDelay = 5 * time.Second
	} else if build.Release == "standard" {
		SafeMutexDelay = 8 * time.Second
	} else if build.Release == "testing" {
		SafeMutexDelay = 3 * time.Second
	}
}
