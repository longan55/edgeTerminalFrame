package global

import "github.com/robfig/cron/v3"

var Cron *cron.Cron

func InitCron() {
	cron.New()
}
