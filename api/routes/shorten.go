package routes

import (
	"time"

	"github.com/gofiber/fiber"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"x-rate-remaining"`
	XRateLimitReset time.Duration `json:"x-rate-limit-reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	req := new(request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if req.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL is required",
		})
	}

	if req.Expiry == 0 {
		req.Expiry = time.Duration(24) * time.Hour
	}

	if req.CustomShort == "" {
		req.CustomShort = generateShortURL()
	}

	if err := storeURL(req.URL, req.CustomShort, req.Expiry); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot store URL",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"short":  req.CustomShort,
		"url":    req.URL,
		"expiry": req.Expiry,
	})
}

func generateShortURL() string {
	return "abc123"
}
