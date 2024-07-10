package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/venqoi/racc-api/utils"
	v1 "github.com/venqoi/racc-api/v1"
)

func main() {
	godotenv.Load()

	raccImages, _ := os.ReadDir("raccs")
	raccVideos, _ := os.ReadDir("raccs/videos")
	utils.NUMBER_OF_IMAGES = len(raccImages)
	utils.NUMBER_OF_VIDEOS = len(raccVideos)

	if err := utils.LoadRaccAlts("utils/alt.json"); err != nil {
		log.Printf("could not load alt text, using default response: %s", err)
	}

	app := fiber.New(fiber.Config{
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"10.50.0.0/24"},
	})

	app.Use(recover.New(recover.Config{
		Next:             nil,
		EnableStackTrace: true,
	}))

	app.Use(logger.New(logger.Config{
		Format: "${time} |   ${cyan}${status} ${reset}|   ${latency} | ${ip} on ${cyan}${ua} ${reset}| ${cyan}${method} ${reset}${path} \n",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        500,
		Expiration: 30 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(utils.Response{
				Success: false,
				Message: "You are being rate limited",
			})
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.GetReqHeaders()["X-Forwarded-For"]
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(utils.Response{
			Success: true,
			Message: "trash panda discovered! you've found the api.",
		})
	})

	v1Group := app.Group("/v1")
	v1Group.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(utils.Response{
			Success: true,
			Message: "trash pandas favourite v1, welcome to the trash can.",
		})
	})

	v1Group.Get("/raccoons", v1.GetRaccoons)
	v1Group.Get("/raccoon", v1.GetRaccoon)
	v1Group.Get("/raccoon/:index", v1.GetRaccoonByIndex)
	v1Group.Get("/raccoftheday", v1.GetRaccoonOfTheDay)
	v1Group.Get("/racchour", v1.GetRaccHour)
	v1Group.Get("/raccofthehour", v1.GetRaccHour)
	v1Group.Get("/video", v1.GetRaccoonVideo)
	v1Group.Get("/video/:index", v1.GetRaccoonVideoByIndex)

	v1Group.Get("/fact", v1.GetRaccFact)
	v1Group.Get("/facts", v1.GetRaccFacts)

	var port = os.Getenv("PORT")

	if len(port) == 0 {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
