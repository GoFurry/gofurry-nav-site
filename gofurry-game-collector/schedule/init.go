package schedule

import (
	game "github.com/GoFurry/gofurry-game-collector/collector/game/controller"
	"github.com/GoFurry/gofurry-game-collector/common/log"
)

func InitSchedule() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	// 初始化 Game 采集模块
	game.GameApi.InitGameCollection()

}
