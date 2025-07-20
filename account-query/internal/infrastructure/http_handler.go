package infrastructure

import (
	"account-query/internal/infrastructure/repository"

	"github.com/gofiber/fiber/v2"
)

type HTTPServer struct {
	app  *fiber.App
	repo *repository.MongoAccountRepository
}

func NewHTTPServer(repo *repository.MongoAccountRepository) *HTTPServer {
	app := fiber.New()
	server := &HTTPServer{app: app, repo: repo}

	app.Get("/api/v1/accounts-query/:id", server.getAccount)
	return server
}

func (h *HTTPServer) Listen(addr string) error {
	return h.app.Listen(addr)
}

func (h *HTTPServer) getAccount(c *fiber.Ctx) error {
	id := c.Params("id")
	acc, err := h.repo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "account not found"})
	}
	return c.JSON(acc)
}
