package routers

/*
 * @Desc: 路由层
 * @author: 福狼
 * @version: v1.0.0
 */

import (
	"os"
	"sync"
	"time"

	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/GoFurry/gofurry-game-backend/middleware"
	"github.com/GoFurry/gofurry-game-backend/roof/env"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var Router *router

type router struct{}

func NewRouter() *router {
	return &router{}
}

func init() {
	Router = NewRouter()
}

var once = sync.Once{}

func (router *router) Init() *fiber.App {
	once.Do(func() {
	})

	cfg := env.GetServerConfig()

	app := fiber.New(fiber.Config{
		Network:                 cfg.Server.Network, // tcp tcp4 tcp6 三种模式
		AppName:                 common.COMMON_PROJECT_NAME,
		ServerHeader:            "GoFurry-Nav",
		Prefork:                 cfg.Server.EnablePrefork,   // 多核cpu处理计算密集型任务 业务量小、IO密集型需关闭
		EnablePrintRoutes:       cfg.Server.Mode == "debug", // 在生产环境禁用错误堆栈跟踪
		ErrorHandler:            customErrorHandler,         // 统一错误处理
		EnableTrustedProxyCheck: true,                       // 信任 Nginx 反向代理
		ReadTimeout:             5 * time.Second,
		WriteTimeout:            10 * time.Second,
	})

	// 注册全局中间件
	registerMiddlewares(app)

	// 路由分组
	gameApi(app.Group("/api/game"))
	recommendApi(app.Group("/api/recommend"))
	searchApi(app.Group("/api/search"))
	reviewApi(app.Group("/api/review"))
	prizeApi(app.Group("/api/prize"))

	app.Get("/api/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	return app
}

// registerMiddlewares 注册中间件
func registerMiddlewares(app *fiber.App) {
	cfg := env.GetServerConfig()
	// 恢复 panic
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.Server.Mode == "debug", // 仅调试模式打印堆栈
	}))

	// 跨域中间件
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Middleware.Cors.AllowOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400, // 预检请求缓存 24 小时
	}))

	// 请求限流
	if cfg.Middleware.Limiter.IsOn {
		app.Use(limiter.New(limiter.Config{
			Max:        cfg.Middleware.Limiter.MaxRequests,              // 单位时间最大请求数
			Expiration: cfg.Middleware.Limiter.Expiration * time.Second, // 时间窗口
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP() // 按 IP 限流
			},
			LimitReached: func(c *fiber.Ctx) error {
				return common.NewResponse(c).ErrorWithCode("请求过于频繁, 请稍后再试", fiber.StatusTooManyRequests)
			},
		}))
	}

	// WAF 中间件
	if cfg.Waf.WafSwitch {
		app.Use(middleware.CorazaMiddleware())
	}

	// 调试模式专属
	if cfg.Server.Mode == "debug" {
		// pprof 性能分析
		app.Use(pprof.New())

		// Swagger 文档
		if cfg.Middleware.Swagger.IsOn {
			// 校验 Swagger 文件是否存在
			if _, err := os.Stat(cfg.Middleware.Swagger.FilePath); os.IsNotExist(err) {
				panic("Swagger 文件不存在: " + cfg.Middleware.Swagger.FilePath)
			}
			swaggerCfg := swagger.Config{
				BasePath: cfg.Middleware.Swagger.BasePath,
				FilePath: cfg.Middleware.Swagger.FilePath,
				Path:     cfg.Middleware.Swagger.Path,
				Title:    cfg.Middleware.Swagger.Title,
			}
			app.Use(swagger.New(swaggerCfg))
		}
	}

	// Prometheus
	app.Use(middleware.PrometheusMiddleware)
	app.Get("/metrics", middleware.MetricsHandler)

	// IP地理位置统计 本地GeoIP + API接入 跳过/metrics
	app.Use(func(c *fiber.Ctx) error {
		if c.Path() == "/metrics" {
			return c.Next()
		}
		return middleware.GeoIPStat(c)
	})
}

// customErrorHandler 自定义错误处理
func customErrorHandler(c *fiber.Ctx, err error) error {
	// 获取错误状态码
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// 标准化错误响应
	response := common.NewResponse(c)
	switch code {
	case fiber.StatusNotFound:
		return response.ErrorWithCode("链接不存在", code)
	case fiber.StatusMethodNotAllowed:
		return response.ErrorWithCode("方法不存在", code)
	case fiber.StatusRequestTimeout:
		return response.ErrorWithCode("请求超时", code)
	default:
		// 生产环境隐藏具体错误信息
		if env.GetServerConfig().Server.Mode != "debug" {
			return response.ErrorWithCode("服务器内部错误", code)
		}
		return response.ErrorWithCode(err.Error(), code)
	}
}
