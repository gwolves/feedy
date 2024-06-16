package http

import (
	"fmt"
	"log/slog"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/gwolves/feedy/internal/service"
)

func NewServer(port int, uc *service.UseCase, logger *slog.Logger) *Server {
	handler := &functionHandler{u: uc, logger: logger}
	return &Server{
		port:    port,
		handler: handler,
		logger:  logger,
	}
}

type Server struct {
	port    int
	handler *functionHandler
	logger  *slog.Logger
}

// Note: error response 시 channeltalk에서 에러 응답에 대한 피드백이 불가능하여
// 요청을 항상 성공하게 하고, 필요한 경우 알림을 통해 피드백
func (s *Server) Serve() error {
	app := fiber.New(fiber.Config{
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	app.Use(
		recover.New(),
		slogfiber.New(s.logger),
	)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	app.Put("/channeltalk/function", func(c *fiber.Ctx) error {
		var req functionRequest
		if err := c.BodyParser(&req); err != nil {
			s.logger.Error("invalid parameter", "error", err)
			return c.JSON(succeedResponse)
		}

		ctx := c.Context()
		res, err := s.handler.Handle(ctx, &req)
		if err != nil {
			s.logger.Error("unexpected error", "error", err)
			return c.JSON(succeedResponse)
		}

		return c.JSON(res)
	})

	return app.Listen(fmt.Sprintf(":%d", s.port))
}
