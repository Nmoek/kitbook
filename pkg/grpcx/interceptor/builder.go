package interceptor

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"net"
	"strings"
)

type Builder struct {
}

// @func: PeerName
// @date: 2024-01-23 01:15:33
// @brief: 获取对端应用名称
// @author: Kewin Li
// @receiver b
// @param ctx
// @return string
func (b *Builder) PeerName(ctx context.Context) string {
	return b.grpcHeaderValue(ctx, "app")
}

// @func: PeerIP
// @date: 2024-01-23 01:15:27
// @brief: 获取对端ip
// @author: Kewin Li
// @receiver b
// @param ctx
// @return string
func (b *Builder) PeerIP(ctx context.Context) string {
	// 如果在 ctx 里面传入。或者说客户端里面设置了，就直接用它设置的
	// 有些时候你经过网关之类的东西，就需要客户端主动设置，防止后面拿到网关的 IP
	clientIP := b.grpcHeaderValue(ctx, "client-ip")
	if clientIP != "" {
		return clientIP
	}

	// 从grpc里取对端ip
	pr, ok2 := peer.FromContext(ctx)
	if !ok2 {
		return ""
	}
	if pr.Addr == net.Addr(nil) {
		return ""
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	if len(addSlice) > 1 {
		return addSlice[0]
	}
	return ""
}

// @func: grpcHeaderValue
// @date: 2024-01-23 01:15:42
// @brief: 获取指定头部字段
// @author: Kewin Li
// @receiver b
// @param ctx
// @param key
// @return string
func (b *Builder) grpcHeaderValue(ctx context.Context, key string) string {
	if key == "" {
		return ""
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	return strings.Join(md.Get(key), ";")
}
