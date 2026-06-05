package middleware

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"go.uber.org/zap"
)

func NewLoggerInterceptor(logger *zap.Logger) connect.Interceptor {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			startedAt := time.Now()

			res, err := next(ctx, req)

			fields := []zap.Field{
				zap.String("procedure", req.Spec().Procedure),
				zap.Duration("duration", time.Since(startedAt)),
			}
			if err != nil {
				logger.Error("rpc request failed", append(
					fields,
					zap.String("code", connect.CodeOf(err).String()),
					zap.Error(err),
				)...)
				return res, err
			}

			logger.Info("rpc request completed", append(
				fields,
				zap.String("code", connect.CodeOf(err).String()),
			)...)
			return res, nil
		})
	})
}
