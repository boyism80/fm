package fault

import (
	"context"
	"path"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (in *Injector) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		delay, mode := in.find(path.Base(method))
		if err := in.wait(ctx, delay); err != nil {
			return status.FromContextError(err).Err()
		}

		switch mode {
		case Unreachable:
			return status.Error(codes.Unavailable, "fault: unreachable")
		case Lost:
			if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
				return err
			}
			return status.Error(codes.Unavailable, "fault: response lost")
		case Timeout:
			<-ctx.Done()
			return status.FromContextError(ctx.Err()).Err()
		default:
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}
}
