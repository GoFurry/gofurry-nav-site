package schedule

import (
	"fmt"
	"time"

	"github.com/GoFurry/gofurry-nav-backend/apps/schedule/task"
	"github.com/GoFurry/gofurry-nav-backend/common/log"
	cs "github.com/GoFurry/gofurry-nav-backend/common/service"
)

// 初始化
func InitScheduleOnStart() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("[InitScheduleOnStart] receive InitScheduleOnStart recover: %v", err))
		}
	}()
	log.Debug("[Schedule] init start 模块初始化开始...")

	//初始化后执行一次 Schedule
	go Schedule()
	go MetricsCache()
	go OneHourTask()
	// 定时任务执行 Schedule
	cs.AddCronJob(10*time.Minute, Schedule)
	cs.AddCronJob(1*time.Minute, MetricsCache)
	cs.AddCronJob(1*time.Hour, OneHourTask)

	log.Debug("[Schedule] init end 模块初始化结束...")
}

func OneHourTask() {
	task.UpdateChangeLog()
}

// 任务表
func Schedule() {
	// task 任务
	task.UpdateTopCountCache()
	task.UpdateLatestPingLog()
	task.UpdateSiteListCache()
	task.UpdateGroupListCache()
}

// 指标采集
func MetricsCache() {
	task.UpdateMetricsCache()
}
