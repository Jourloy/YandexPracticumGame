package app

import "time"

type service struct{}

var startTime = time.Now()

func InitService() *service {
	return &service{}
}

type healthCheckResp struct {
	Uptime int
}

func (s *service) HealthCheck() healthCheckResp {
	return healthCheckResp{
		Uptime: int(time.Since(startTime)),
	}
}
