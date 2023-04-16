package middlewares

import (
	"NORSI-TRANS/appErrors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"time"
)

func LoggingMiddlewares() fiber.Handler {
	logConfig := logger.Config{
		Format: fmt.Sprintf("%s ${time} ${method} ${path} - ${latency} ${status} ${error}\n",
			time.Now().Format("2006/01/02")),
	}

	return logger.New(logConfig)
}

func ErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		if err == nil {
			return nil
		}

		if fiberErr, ok := err.(*fiber.Error); ok {
			return ctx.Status(fiberErr.Code).Format(fiberErr.Message)
		}

		appErr := appErrors.AppErrorFromError(err)

		return ctx.Status(appErr.Code).Format(appErr)
	}
}
