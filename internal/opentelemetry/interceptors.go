package opentelemetry

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func UnaryClientInterceptor(opts ...otelgrpc.Option) grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor(opts...)
}

func StreamClientInterceptor(opts ...otelgrpc.Option) grpc.StreamClientInterceptor {
	return otelgrpc.StreamClientInterceptor(opts...)
}

func UnaryServerInterceptor(opts ...otelgrpc.Option) grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor(opts...)
}

func StreamServerInterceptor(opts ...otelgrpc.Option) grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor(opts...)
}
