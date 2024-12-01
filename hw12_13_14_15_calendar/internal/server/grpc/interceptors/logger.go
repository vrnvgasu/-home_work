package validate

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Logger interface {
	File(msg string)
	Info(msg string)
	Error(msg string)
}

func UnaryServerRequestLoggerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		ip := ""
		if p, ok := peer.FromContext(ctx); ok {
			ip = p.Addr.String()
		}

		userAgent := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok && len(md["user-agent"]) > 0 {
			userAgent = md["user-agent"][0]
		}

		next, err := handler(ctx, req)
		if err != nil {
			logger.Error(err.Error())
		}

		msg := fmt.Sprintf(
			"%s [%s] %s %s %d %s '%s'",
			ip,
			time.Now().UTC(),
			info.FullMethod,
			"HTTP/2",
			http.StatusOK,
			time.Since(start),
			userAgent,
		)
		logger.Info(msg)
		logger.File(msg)

		return next, err
	}
}
