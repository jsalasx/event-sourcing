package infrastructure

import (
	"account-cmd/internal/application"

	"github.com/gofiber/fiber/v2"
)

type HTTPServer struct {
	app     *fiber.App
	service *application.AccountService
}

func NewHTTPServer(service *application.AccountService) *HTTPServer {
	app := fiber.New()
	server := &HTTPServer{app: app, service: service}

	app.Post("api/v1/accounts", server.createAccount)
	app.Post("api/v1/accounts/:id/deposit", server.deposit)
	app.Post("api/v1/accounts/:id/withdraw", server.withdraw)

	return server
}

func (h *HTTPServer) Listen(addr string) error {
	return h.app.Listen(addr)
}

func (h *HTTPServer) createAccount(c *fiber.Ctx) error {
	var req struct {
		Owner  string  `json:"owner"`
		Amount float64 `json:"amount"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	acc, err := h.service.CreateAccount(req.Owner, req.Amount)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(acc)
}

func (h *HTTPServer) deposit(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	acc, err := h.service.Deposit(id, req.Amount)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(acc)
}

func (h *HTTPServer) withdraw(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	acc, err := h.service.Withdraw(id, req.Amount)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(acc)
}
