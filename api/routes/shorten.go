package routes

import (
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"

	"github.com/ravikisha/url-shortener/helpers"
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
	// Converting request
	req := new(request)

	// Parse JSON into request struct
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// TODO: Implement rate limiting

	// Check if URL is empty
	if req.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL is required",
		})
	}

	// Validate the URL
	if !govalidator.IsURL(req.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	// Check the Domain Name Error
	if !helpers.RemoveDomainNameError(req.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Domain Name Error",
		})
	}

	// Set default expiry time (24 hours)
	if req.Expiry == 0 {
		req.Expiry = time.Duration(24) * time.Hour
	}

	// Generate short URL
	if req.CustomShort == "" {
		req.CustomShort = generateShortURL()
	}

	// Enforce Https, SSL
	req.URL = helpers.EnforceHTTP(req.URL)

	// Store the URL
	if err := storeURL(req.URL, req.CustomShort, req.Expiry); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot store URL",
		})
	}

	// Return response
	return c.JSON(response{
		URL:             req.URL,
		CustomShort:     req.CustomShort,
		Expiry:          req.Expiry,
		XRateRemaining:  os.Getenv("APP_QUOTA"),
		XRateLimitReset: time.Duration(24) * time.Hour,
	})
}

func generateShortURL() string {
	return "abc123"
}
