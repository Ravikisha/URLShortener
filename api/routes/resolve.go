package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"

	"github.com/ravikisha/url-shortener/database"
)

func ResolveURL(c *fiber.Ctx) error {

	url := c.Params("url")

	// Create a Database Client
	rdb := database.NewClient(0)

	// Close the connection
	defer rdb.Close()

	// Get the URL from the database
	val, err := rdb.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	// Create a Database Client
	rInr := database.NewClient(1)
	defer rInr.Close()

	// Increment the number of times the URL has been resolved
	_ = rInr.Incr(database.Ctx, "counter")
	return c.Redirect(val, fiber.StatusMovedPermanently)
}
