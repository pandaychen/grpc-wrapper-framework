package authentication

import (
	"context"

	"grpc-wrapper-framework/common/enums"
	"grpc-wrapper-framework/common/pyerrors"
	"grpc-wrapper-framework/common/vars"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// 提供服务端 RPC 的验证机制
type Authenticator struct {
	Ctx *context.Context
	//Appid    string		//通用的auth不应该包括私有信息
	//Apptoken string
	//AuthType enums.AuthType
}

func NewAuthenticator(ctx *context.Context) (*Authenticator, error) {
	return &Authenticator{
		Ctx: ctx,
	}, nil
}

// Authenticate authenticates the given ctx
func (a *Authenticator) Authenticate(ctx context.Context) error {
	var (
		appid    string
		apptoken string
		method   enums.AuthType
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		//not found context
		return status.Error(codes.Unauthenticated, pyerrors.AuthenticatorMissing)
	}

	//extractor auth method from ctx
	if val, ok := md[vars.APP_AUTH_METHOD]; ok {
		if len(val) != 0 {
			if len(val[0]) != 0 {
				method = enums.AuthType(val[0])
			} else {
				method = enums.AUTH_TYPE_APP
			}
		} else {
			method = enums.AUTH_TYPE_APP
		}
	} else {
		method = enums.AUTH_TYPE_APP
	}

	switch method {
	case enums.AUTH_TYPE_APP:
		if val, ok := md[vars.APPKEY_NAME]; ok {
			if len(val) == 0 {
				return status.Error(codes.Unauthenticated, pyerrors.AuthenticatorMissing)
			}
			appid = val[0]
		}
		if val, ok := md[vars.APPTOKEN_NAME]; ok {
			if len(val) == 0 {
				return status.Error(codes.Unauthenticated, pyerrors.AuthenticatorMissing)
			}
			apptoken = val[0]
		}
		if appid == "" || apptoken == "" {
			return status.Error(codes.Unauthenticated, pyerrors.AuthenticatorMissing)
		}

		return a.Validate(appid, apptoken, method)
	case enums.AUTH_TYPE_JWT:
		//TODO: fix jwt auth
		return status.Error(codes.Internal, pyerrors.InternalError)
	default:
		return status.Error(codes.Internal, pyerrors.InternalError)
	}
}

func (a *Authenticator) Validate(appid, token string, method enums.AuthType) error {
	switch method {
	case enums.AUTH_TYPE_APP:
		if appid != "test" && token != "test" {
			return status.Error(codes.Unauthenticated, pyerrors.TokenVerifyInvalid)
		}
	case enums.AUTH_TYPE_JWT:
		//TODO: fix jwt auth
		return status.Error(codes.Internal, pyerrors.InternalError)
	}

	return nil
}
