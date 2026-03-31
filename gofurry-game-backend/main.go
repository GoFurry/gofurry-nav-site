package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/GoFurry/gofurry-game-backend/apps/schedule"
	"github.com/GoFurry/gofurry-game-backend/common"
	gfLog "github.com/GoFurry/gofurry-game-backend/common/log"
	cs "github.com/GoFurry/gofurry-game-backend/common/service"
	"github.com/GoFurry/gofurry-game-backend/middleware"
	"github.com/GoFurry/gofurry-game-backend/roof/env"
	"github.com/GoFurry/gofurry-game-backend/routers"
	"github.com/kardianos/service"
)

//@title GoFurry-Game-Backend
//@version v1.0.0
//@description GoFurry-Game-Backend

var (
	errChan = make(chan error)
)

func main() {
	dir, _ := os.Getwd()

	svcConfig := &service.Config{
		Name:        common.COMMON_PROJECT_NAME,
		DisplayName: "gf-game",
		Description: "gf-game",
		Option: service.KeyValue{
			"SystemdScript": `[Unit]
Description=gf-game
After=network.target
Requires=network.target

[Service]
Type=simple
WorkingDirectory=` + dir + `/
ExecStart=` + dir + `/gf-game
Restart=always
RestartSec=30
LogOutput=true
LogDirectory=/var/log/gf-game
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target`,
		},
	}
	prg := &goFurry{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		slog.Error(err.Error())
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			err = s.Install()
			if err != nil {
				slog.Error("服务安装失败: ", err)
			} else {
				slog.Info("服务安装成功.")
			}
			return
		case "uninstall":
			err = s.Uninstall()
			if err != nil {
				slog.Error("服务卸载失败: ", err)
			} else {
				slog.Info("服务卸载成功.")
			}
			return
		case "version":
			slog.Info("gf-game V1.0.0")
			return
		case "help":
			slog.Info(common.COMMON_PROJECT_HELP)
			return
		}
		return
	}

	// 内存限制和 GC 策略
	debug.SetGCPercent(env.GetServerConfig().Server.GCPercent)
	debug.SetMemoryLimit(int64(env.GetServerConfig().Server.MemoryLimit << 30))

	// 初始化系统服务
	InitOnStart()
	// 启动系统
	err = s.Run()
	if err != nil {
		slog.Error(err.Error())
	}
}

type goFurry struct{}

func InitOnStart() {
	cfg := env.GetServerConfig()
	// 初始化自定义日志
	logCfg := &gfLog.Config{
		ShowLine:   true,
		TimeFormat: common.TIME_FORMAT_DATE,
	}
	if cfg.Server.Mode == "debug" {
		logCfg.Level = "debug"
		logCfg.Mode = "dev"
		logCfg.EncodeJson = false
	} else {
		logCfg.Level = cfg.Log.LogLevel
		logCfg.Mode = cfg.Log.LogMode
		logCfg.FilePath = cfg.Log.LogPath
		logCfg.MaxSize = cfg.Log.LogMaxSize
		logCfg.MaxBackups = cfg.Log.LogMaxBackups
		logCfg.MaxAge = cfg.Log.LogMaxAge
		logCfg.Compress = true
		logCfg.EncodeJson = true
		logCfg.TimeFormat = common.TIME_FORMAT_LOG
	}

	// 初始化自定义日志
	err := gfLog.InitLogger(logCfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// 初始化 Prometheus 中间件
	middleware.InitPrometheus(middleware.FiberPromConf{
		SkipPaths:         []string{},
		IgnoreStatusCodes: []int{},
	})
	// 初始化 GeoIP 中间件
	middleware.InitGeoIP()
	// 初始化 Coraza 中间件
	if cfg.Waf.WafSwitch {
		middleware.InitGlobalWAF(cfg.Waf.ConfPath)
	}
	// 初始化 redis
	cs.InitRedisOnStart()
	// 初始化时间调度
	cs.InitTimeWheelOnStart()

	// 初始化定时任务
	schedule.InitScheduleOnStart()
}

func (gf *goFurry) Start(s service.Service) error {
	go gf.run()
	return nil
}

func (gf *goFurry) run() {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	// 启动 web
	go func() {
		app := routers.Router.Init()

		addr := env.GetServerConfig().Server.IPAddress + ":" + env.GetServerConfig().Server.Port
		// nginx 完成 https 就不使用 TLS
		//pem := env.GetServerConfig().Key.TlsPem
		//key := env.GetServerConfig().Key.TlsKey
		//if err := app.ListenTLS(addr, pem, key); err != nil {
		//	fmt.Println(err)
		//	errChan <- err
		//}
		if err := app.Listen(addr); err != nil {
			fmt.Println(err)
			errChan <- err
		}
	}()
	if err := <-errChan; err != nil {
		slog.Error(err.Error())
	}
}

func (gf *goFurry) Stop(s service.Service) error {
	return nil
}
