package util

import (
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

/*
 * @Desc: gRPC连接池
 * @author: 福狼
 * @version: v1.0.0
 */

var (
	grpcConns           = make(map[string]*grpc.ClientConn)
	mu                  sync.RWMutex
	healthCheckInterval = 5 * time.Minute // 健康检查间隔
)

// 初始化时启动健康检查协程
func init() {
	go startHealthCheck()
}

// GetGrpcClientConn 获取指定服务的 gRPC 连接
func GetGrpcClientConn(serviceName string, creds *credentials.TransportCredentials) (newConn *grpc.ClientConn, err error) {
	// 先尝试读锁获取已有连接
	mu.RLock()
	conn, exists := grpcConns[serviceName]
	mu.RUnlock()

	// 连接存在且状态正常
	if exists && isConnHealthy(conn) {
		return conn, nil
	}

	// 连接不存在或不可用
	mu.Lock()
	defer mu.Unlock()

	// 防止并发创建
	if conn, exists = grpcConns[serviceName]; exists && isConnHealthy(conn) {
		return conn, nil
	}

	// 创建新连接
	target := "etcd:///" + serviceName
	if creds == nil {
		newConn, err = grpc.NewClient(
			target,
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	} else {
		newConn, err = grpc.NewClient(
			target,
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
			grpc.WithTransportCredentials(*creds),
		)
	}

	if err != nil {
		return nil, err
	}

	// 存入连接池
	grpcConns[serviceName] = newConn
	return newConn, nil
}

// CloseGrpcConns 程序退出时关闭所有 gRPC 连接
func CloseGrpcConns() {
	mu.Lock()
	defer mu.Unlock()
	for _, conn := range grpcConns {
		_ = conn.Close()
	}
	grpcConns = make(map[string]*grpc.ClientConn)
}

// 启动健康检查协程
func startHealthCheck() {
	ticker := time.NewTicker(healthCheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		checkAndRebuildConnections()
	}
}

// 检查所有连接状态
func checkAndRebuildConnections() {
	mu.RLock()
	// 复制当前连接列表
	serviceNames := make([]string, 0, len(grpcConns))
	for name := range grpcConns {
		serviceNames = append(serviceNames, name)
	}
	mu.RUnlock()

	// 逐个检查连接状态
	for _, serviceName := range serviceNames {
		mu.RLock()
		conn, exists := grpcConns[serviceName]
		mu.RUnlock()

		if !exists {
			continue
		}

		// 检查连接是否健康
		if !isConnHealthy(conn) {
			// 尝试重建
			if err := rebuildConnection(serviceName); err != nil {
				continue
			}
		}
	}
}

// 判断连接是否健康
func isConnHealthy(conn *grpc.ClientConn) bool {
	state := conn.GetState()
	// 健康状态：Idle（空闲）、Ready（就绪）、Connecting（连接中）
	// 不健康状态：TransientFailure（临时失败）、Shutdown（已关闭）
	return state != connectivity.TransientFailure && state != connectivity.Shutdown
}

// 重建指定服务的连接
func rebuildConnection(serviceName string) error {
	mu.Lock()
	defer mu.Unlock()

	// 再次检查连接状态
	oldConn, exists := grpcConns[serviceName]
	if exists && isConnHealthy(oldConn) {
		return nil
	}

	// 创建新连接
	target := "etcd:///" + serviceName
	newConn, err := grpc.NewClient(
		target,
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	// 关闭旧连接
	if exists {
		_ = oldConn.Close()
	}

	// 用新连接替换旧连接
	grpcConns[serviceName] = newConn
	return nil
}
