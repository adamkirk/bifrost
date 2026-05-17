package server

import (
	"context"

	"github.com/adamkirk/bifrost/api/internal/common"
	"github.com/danielgtaylor/huma/v2"
)

func ErrorHandler[Req any, Resp any](handler func(context.Context, *Req) (*Resp, error)) func(ctx context.Context, req *Req) (*Resp, error) {
	return func(ctx context.Context, req *Req) (*Resp, error) {
		resp, err := handler(ctx, req)

		if err == nil {
			return resp, nil
		}

		switch e := err.(type) {
		// case validation.Error:
		// 	return resp, buildValidationError(e)

		case common.ErrUnauthorised:
			return resp, huma.Error401Unauthorized(e.Message)

		default:
			// TODO: outside dev, return a generic errors message instead of the
			// actual error as it appears in the response
			return resp, err
		}
	}
}
