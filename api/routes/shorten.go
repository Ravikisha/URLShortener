package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"

	"github.com/ravikisha/url-shortener/database"
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

	// Implement rate limiting
	rdb := database.NewClient(2)
	defer rdb.Close()

	// Get the number of requests made
	val, err := rdb.Get(database.Ctx, c.IP()).Result()

	if err == redis.Nil {
		_ = rdb.Set(database.Ctx, c.IP(), os.Getenv("APP_QUOTA"), 30*60*time.Second).Err()
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	} else {
		// Check if the number of requests is greater than the quota
		limit, _ := strconv.Atoi(os.Getenv("APP_QUOTA"))
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":              "Rate limit exceeded",
				"x-rate-limit-reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

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

	// Decrement the number of requests made
	rdb.Decr(database.Ctx, c.IP())

	// Return response
	return c.JSON(response{
		URL:             req.URL,
		CustomShort:     req.CustomShort,
		Expiry:          req.Expiry,
		XRateRemaining:  1,
		XRateLimitReset: time.Duration(24) * time.Hour,
	})
}

func generateShortURL() string {
	return "abc123"
}
