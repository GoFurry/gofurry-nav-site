package schedule

import (
	"fmt"
	"time"

	"github.com/GoFurry/gofurry-game-backend/apps/schedule/task"
	"github.com/GoFurry/gofurry-game-backend/common/log"
	cs "github.com/GoFurry/gofurry-game-backend/common/service"
)

// 初始化
func InitScheduleOnStart() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("receive InitScheduleOnStart recover: %v", err))
		}
	}()
	log.Info("Schedule 模块初始化开始...")

	//初始化后执行一次 Schedule
	go ScheduleByTenMinutes()
	go ScheduleByOneHour()
	go ScheduleByHalfDay()
	// 定时任务执行 Schedule
	cs.AddCronJob(10*time.Minute, ScheduleByTenMinutes)
	cs.AddCronJob(1*time.Hour, ScheduleByOneHour)
	cs.AddCronJob(12*time.Hour, ScheduleByHalfDay)

	log.Info("Schedule 模块初始化结束...")
}

// 十分钟任务表
func ScheduleByTenMinutes() {
	// 缓存游戏模块主页分组内容
	task.UpdateMainInfoCache()
}

// 一小时任务表
func ScheduleByOneHour() {
	// 缓存游戏资讯面板数据
	task.UpdateGamePanelCache()
	// 缓存更新公告数据
	task.UpdateGameNewsCache()
	// 缓存创作者数据
	task.UpdateGameCreatorCache()
	// 缓存更多更新公告
	task.UpdateMoreGameNewsCache()

	// 抽奖系统感知兑奖
	task.UpdatePrizeWinner()
}

// 半天任务表
func ScheduleByHalfDay() {

}
