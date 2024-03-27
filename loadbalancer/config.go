package loadbalancer

import "time"

type Config struct {
	Port    int
	Timeout time.Duration
}
