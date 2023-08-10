package middleware

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func LoggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	err := c.Next()

	// Logging
	log.Printf("%s - %s - %s", c.IP(), c.Method(), c.Path(), time.Since(start))

	return err
}
