package response

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func New(ctx *fiber.Ctx) *Response {
	return &Response{
		ctx: ctx,
	}
}

func (r *Response) SetMessage(message string) *Response {
	r.Message = message
	return r
}

func (r *Response) SetData(data any) *Response {
	r.Data = data
	return r
}

func (r *Response) SetErr(err error) *Response {
	r.error = err
	return r
}

func (r *Response) SetHttpCode(httpCode int) *Response {
	r.HttpCode = httpCode
	return r
}

func (r *Response) Send(arg ...string) (resp error) {
	args := strings.Join(arg, "|")

	//valida http code
	if r.HttpCode <= 0 {
		r.HttpCode = fiber.StatusInternalServerError
	}

	if r.HttpCode < fiber.StatusContinue {
		r.HttpCode = fiber.StatusOK
	}

	//validate message for http code
	switch r.HttpCode / 100 {
	case fiber.StatusOK / 100:
		r.Message = r.Message + " successfully"
	case fiber.StatusBadRequest / 100:
		//replace message from args
		if strings.TrimSpace(args) != "" {
			r.Message = args
		}

		if strings.TrimSpace(args) == "" && r.error != nil {
			r.Message = errors.Cause(r.error).Error()
		}

	case fiber.StatusInternalServerError / 100:
		r.Message = "please try again"
	}

	resp = r.ctx.Status(r.HttpCode).JSON(&r)
	return resp
}
