package auth

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type Authorizer struct {
	Username, Password string
}

func NewAuthorizer(username string, password string) *Authorizer {
	return &Authorizer{Username: username, Password: password}
}

func (a *Authorizer) HandleUnary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx, err := a.Context(ctx)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func (a *Authorizer) HandleStream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	wrap := grpc_middleware.WrapServerStream(stream)
	ctx := wrap.Context()
	ctx, err := a.Context(ctx)
	if err != nil {
		return err
	}
	wrap.WrappedContext = ctx
	return handler(srv, stream)
}

func (a *Authorizer) authorize(ctx context.Context) (context.Context, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("md = %#v\n", md)
		if len(md["username"]) > 0 && md["username"][0] == "admin" && len(md["password"]) > 0 && md["password"][0] == "admin123" {
			return ctx, nil
		}
		return nil, fmt.Errorf("Access denied")
	}
	return ctx, fmt.Errorf("Metadata is empty")
}

func (a *Authorizer) Verify(username string, password string) error {
	if username == "admin" && password == "admin123" {
		return nil
	}
	return fmt.Errorf("AccessDeniedError")
}

func (a *Authorizer) Context(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "Authorization required: no context metadata")
	}

	if err := a.Verify(md["username"][0], md["password"][0]); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, err.Error())
	}
	return ctx, nil
}
