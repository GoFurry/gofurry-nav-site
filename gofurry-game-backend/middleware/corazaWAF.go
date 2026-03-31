package middleware

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/GoFurry/gofurry-game-backend/common"
	"github.com/corazawaf/coraza/v3"
	"github.com/corazawaf/coraza/v3/debuglog"
	"github.com/corazawaf/coraza/v3/experimental"
	"github.com/corazawaf/coraza/v3/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

/*
 * @Desc: Fiber Coraza-WAF 中间件
 * @author: 福狼
 * @version: v1.0.2
 */

// 全局 WAF 单例与初始化锁
var (
	globalWAF  coraza.WAF
	wafOnce    sync.Once
	wafInitErr error
)

// CorazaCfg 中间件配置文件
type CorazaCfg struct {
	// 核心配置
	DirectivesFile string // WAF 规则文件路径
	RuleEngine     string // 规则引擎模式 (On/Off/DetectionOnly)
	RootFS         fs.FS  // 文件系统根目录(用于加载规则中的外部文件)

	// 请求体配置
	RequestBodyAccess        bool // 是否启用请求体访问
	RequestBodyLimit         int  // 请求体最大字节数 (默认 10M)
	RequestBodyInMemoryLimit int  // 内存中缓存请求体最大字节数 (默认 128K)

	// 响应体配置
	ResponseBodyAccess    bool     // 是否启用响应体访问
	ResponseBodyLimit     int      // 响应体最大字节数 (默认 512K)
	ResponseBodyMimeTypes []string // 需要处理的响应体 MIME 类型

	// 日志配置
	DebugLogger    debuglog.Logger // 调试日志器 (可选)
	EnableErrorLog bool            // 是否启用错误日志回调
}

// DefaultCorazaCfg 默认配置
func DefaultCorazaCfg() CorazaCfg {
	return CorazaCfg{
		// 核心默认配置
		DirectivesFile: "./conf/coraza.conf", // 默认规则文件路径
		RuleEngine:     "On",                 // 默认启用规则引擎 (拦截模式)
		RootFS:         nil,                  // 默认不设置

		// 请求体默认配置
		RequestBodyAccess:        true,
		RequestBodyLimit:         10 * 1024 * 1024, // 10MB
		RequestBodyInMemoryLimit: 128 * 1024,       // 128KB

		// 响应体默认配置
		ResponseBodyAccess:    false,      // 默认不启用响应体处理 (性能更好)
		ResponseBodyLimit:     512 * 1024, // 512KB
		ResponseBodyMimeTypes: []string{"text/html", "text/plain", "application/json", "application/xml"},

		// 日志默认配置
		EnableErrorLog: true, // 默认启用错误日志
	}
}

// InitGlobalWAFWithCfg 基于配置初始化全局 WAF 单例
func InitGlobalWAFWithCfg(cfg CorazaCfg) {
	wafOnce.Do(func() {
		globalWAF, wafInitErr = createWAFWithCfg(cfg)
		if wafInitErr != nil {
			slog.Error("[CorazaWAF] InitGlobalWAFWithCfg Error", wafInitErr.Error())
		}
	})
}

// InitGlobalWAF 传入 Coraza 配置文件路径完成初始化
func InitGlobalWAF(path ...string) {
	if len(path) > 0 {
		InitGlobalWAFWithCfg(CorazaCfg{
			DirectivesFile: path[0],
		})
	} else {
		InitGlobalWAFWithCfg(DefaultCorazaCfg())
	}
}

// CorazaMiddleware 中间件
func CorazaMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		// 直接使用全局 WAF 实例
		if wafInitErr != nil {
			slog.Error("WAF全局实例初始化失败: ", wafInitErr)
			return common.NewResponse(c).ErrorWithCode("WAF 初始化失败", http.StatusInternalServerError)
		}
		if globalWAF == nil {
			return common.NewResponse(c).ErrorWithCode("WAF 实例未初始化", http.StatusInternalServerError)
		}

		// 事件句柄匿名函数
		newTX := func(*http.Request) types.Transaction {
			return globalWAF.NewTransaction()
		}
		// 事件句柄匿名函数
		if ctxwaf, ok := globalWAF.(experimental.WAFWithOptions); ok {
			newTX = func(r *http.Request) types.Transaction {
				return ctxwaf.NewTransactionWithOptions(experimental.Options{
					Context: r.Context(),
				})
			}
		}

		// FastHTTP 转换为 HTTP
		stdReq, err := convertFasthttpToStdRequest(c)
		if err != nil {
			return common.NewResponse(c).ErrorWithCode("请求转换失败", http.StatusInternalServerError)
		}

		// 开启事件
		tx := newTX(stdReq)
		defer func() {
			// 捕获 panic
			if r := recover(); r != nil {
				slog.Error(fmt.Sprintf("WAF transaction panicked: %v", r))
			}
			// 打印日志
			tx.ProcessLogging()
			// 关闭事件
			if err = tx.Close(); err != nil {
				tx.DebugLogger().Error().Err(err).Msg("Failed to close the transaction")
			}
		}()

		// 没开规则就放行
		if tx.IsRuleEngineOff() {
			return c.Next()
		}

		// 处理请求
		if it, err := processRequest(tx, stdReq); err != nil {
			// 处理失败
			tx.DebugLogger().Error().Err(err).Msg("Failed to process request")
			return common.NewResponse(c).ErrorWithCode("WAF 处理请求失败", http.StatusInternalServerError)
		} else if it != nil {
			// 拦截成功
			status := obtainStatusCodeFromInterruptionOrDefault(it, http.StatusOK)
			c.Status(status)
			c.Set("X-WAF-Blocked", "true")
			return common.NewResponse(c).ErrorWithCode("您的请求存在安全风险, 已被系统拦截.", status)
		}

		// 放行
		return c.Next()
	}
}

// logError WAF 错误日志
func logError(error types.MatchedRule) {
	slog.Warn("WAF rule matched",
		slog.String("severity", string(error.Rule().Severity())),
		slog.String("error_log", error.ErrorLog()),
		slog.Int("rule_id", error.Rule().ID()),
	)
}

// processRequest 拦截部分具体实现操作
func processRequest(tx types.Transaction, req *http.Request) (*types.Interruption, error) {
	var (
		client string
		cport  int
	)
	// IMPORTANT: Some http.Request.RemoteAddr implementations will not contain port or contain IPV6: [2001:db8::1]:8080
	idx := strings.LastIndexByte(req.RemoteAddr, ':')
	if idx != -1 {
		client = req.RemoteAddr[:idx]
		cport, _ = strconv.Atoi(req.RemoteAddr[idx+1:])
	}

	var in *types.Interruption
	// There is no socket access in the request object, so we neither know the server client nor port.
	tx.ProcessConnection(client, cport, "", 0)
	tx.ProcessURI(req.URL.String(), req.Method, req.Proto)
	for k, vr := range req.Header {
		for _, v := range vr {
			tx.AddRequestHeader(k, v)
		}
	}

	// Host will always be removed from req.Headers() and promoted to the
	// Request.Host field, so we manually add it
	if req.Host != "" {
		tx.AddRequestHeader("Host", req.Host)
		// This connector relies on the host header (now host field) to populate ServerName
		tx.SetServerName(req.Host)
	}

	// Transfer-Encoding header is removed by go/http
	// We manually add it to make rules relying on it work (E.g. CRS rule 920171)
	if req.TransferEncoding != nil {
		tx.AddRequestHeader("Transfer-Encoding", req.TransferEncoding[0])
	}

	in = tx.ProcessRequestHeaders()
	if in != nil {
		return in, nil
	}

	if tx.IsRequestBodyAccessible() {
		// We only do body buffering if the transaction requires request
		// body inspection, otherwise we just let the request follow its
		// regular flow.
		if req.Body != nil && req.Body != http.NoBody {
			it, _, err := tx.ReadRequestBodyFrom(req.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to append request body: %w", err)
			}

			if it != nil {
				return it, nil
			}

			rbr, err := tx.RequestBodyReader()
			if err != nil {
				return nil, fmt.Errorf("failed to get the request body: %w", err)
			}

			// Adds all remaining bytes beyond the coraza limit to its buffer
			// It happens when the partial body has been processed and it did not trigger an interruption
			bodyReader := io.MultiReader(rbr, req.Body)
			// req.Body is transparently reinizialied with a new io.ReadCloser.
			// The http handler will be able to read it.
			req.Body = io.NopCloser(bodyReader)
		}
	}

	return tx.ProcessRequestBody()
}

// obtainStatusCodeFromInterruptionOrDefault "deny" Action 拒绝时设置状态码 403
func obtainStatusCodeFromInterruptionOrDefault(it *types.Interruption, defaultStatusCode int) int {
	if it.Action == "deny" {
		statusCode := it.Status
		if statusCode == 0 {
			statusCode = 403
		}

		return statusCode
	}
	return defaultStatusCode
}

// convertFasthttpToStdRequest 转换请求类型
func convertFasthttpToStdRequest(c *fiber.Ctx) (*http.Request, error) {
	stdReq, err := adaptor.ConvertRequest(c, false) // false 表示不自动关闭请求体, 后续需由 WAF 处理
	if err != nil {
		return nil, err
	}

	stdReq.RemoteAddr = net.JoinHostPort(c.IP(), c.Port())

	// 手动设置完整 Host
	if stdReq.Host == "" {
		stdReq.Host = c.Hostname()
	}

	return stdReq, nil
}

// createWAFWithCfg 基于配置创建 WAF 实例
func createWAFWithCfg(cfg CorazaCfg) (coraza.WAF, error) {
	// 环境变量
	if envDirectivesFile := os.Getenv("CORAZA_DIRECTIVES_FILE"); envDirectivesFile != "" {
		cfg.DirectivesFile = envDirectivesFile
	}

	// 验证核心配置有效性
	if cfg.DirectivesFile == "" {
		slog.Warn("WAF 规则文件路径未配置 (directives_file 为空)")
		panic("WAF 规则文件路径未配置 (directives_file 为空)")
	} else {
		if _, err := os.Stat(cfg.DirectivesFile); os.IsNotExist(err) {
			slog.Warn("WAF 规则文件不存在: ", cfg.DirectivesFile)
			panic("WAF 规则文件路径未配置 (directives_file 为空)")
		}
	}

	wafConfig := coraza.NewWAFConfig()

	// 错误回调配置
	if cfg.EnableErrorLog {
		wafConfig = wafConfig.WithErrorCallback(logError)
	}

	// 请求体相关配置
	if cfg.RequestBodyAccess {
		wafConfig = wafConfig.WithRequestBodyAccess()
	}
	if cfg.RequestBodyLimit > 0 {
		wafConfig = wafConfig.WithRequestBodyLimit(cfg.RequestBodyLimit)
	}
	if cfg.RequestBodyInMemoryLimit > 0 {
		wafConfig = wafConfig.WithRequestBodyInMemoryLimit(cfg.RequestBodyInMemoryLimit)
	}

	// 响应体相关配置
	if cfg.ResponseBodyAccess {
		wafConfig = wafConfig.WithResponseBodyAccess()
	}
	if cfg.ResponseBodyLimit > 0 {
		wafConfig = wafConfig.WithResponseBodyLimit(cfg.ResponseBodyLimit)
	}
	if len(cfg.ResponseBodyMimeTypes) > 0 {
		wafConfig = wafConfig.WithResponseBodyMimeTypes(cfg.ResponseBodyMimeTypes)
	}

	// 其他可选配置
	if cfg.RootFS != nil {
		wafConfig = wafConfig.WithRootFS(cfg.RootFS)
	}
	if cfg.DebugLogger != nil {
		wafConfig = wafConfig.WithDebugLogger(cfg.DebugLogger)
	}

	// 加载规则文件
	if cfg.DirectivesFile != "" {
		wafConfig = wafConfig.WithDirectivesFromFile(cfg.DirectivesFile)
	}

	// 创建并返回 WAF 实例
	return coraza.NewWAF(wafConfig)
}
