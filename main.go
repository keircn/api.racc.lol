package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/venqoi/racc-api/utils"
	v1 "github.com/venqoi/racc-api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var requestCountCollection *mongo.Collection
var isDbConnected bool
var notConnectedMessages = []string{
	"a goose stole the database from the trash can",
	"a raccoon ate the database",
	"the database is in the trash",
	"the database scurried away",
	"missing raccoon- i mean database"

}

func main() {
	godotenv.Load()

	mongoURI := os.Getenv("MONGO_DB_URL")
	if mongoURI != "" {
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(context.TODO())

		requestCountCollection = client.Database("racc").Collection("requestCounts")
		isDbConnected = true
	} else {
		isDbConnected = false
	}

	raccImages, _ := os.ReadDir("raccs")
	raccVideos, _ := os.ReadDir("raccs/videos")
	raccTrans, _ := os.ReadDir("raccs/transparent")
	utils.NUMBER_OF_IMAGES = len(raccImages) + 1
	utils.NUMBER_OF_VIDEOS = len(raccVideos) + 1
	utils.NUMBER_OF_TRANS = len(raccTrans) + 1

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
		Max:        30,
		Expiration: 30 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(utils.Response{
				Success: false,
				Message: "You are being rate limited",
			})
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			xfwdFor := c.GetReqHeaders()["X-Forwarded-For"]
			if len(xfwdFor) > 0 {
				return strings.Join(xfwdFor, ",")
			}
			return c.IP()
		},
	}))

	app.Use(func(c *fiber.Ctx) error {
		if c.Path() != "/favicon.ico" && isDbConnected {
			_, err := requestCountCollection.UpdateOne(context.TODO(), bson.M{"_id": "requestCount"}, bson.M{"$inc": bson.M{"count": 1}}, options.Update().SetUpsert(true))
			if err != nil {
				return err
			}
		}
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(utils.Response{
			Success: true,
			Message: "trash panda discovered! you've found the api.",
		})
	})

	app.Get("/stats", func(c *fiber.Ctx) error {
		if isDbConnected {
			var result struct {
				Count int `bson:"count"`
			}
			err := requestCountCollection.FindOne(context.TODO(), bson.M{"_id": "requestCount"}).Decode(&result)
			if err != nil {
				return err
			}
			return c.JSON(utils.StatsResponse{
				Success:           true,
				Message:           "stats for the trash panda api :) ",
				Requests:          result.Count,
				Images:            utils.NUMBER_OF_IMAGES,
				Videos:            utils.NUMBER_OF_VIDEOS,
				TransparentImages: utils.NUMBER_OF_TRANS,
			})
		} else {
			rand.Seed(time.Now().UnixNano())
			randomMessage := notConnectedMessages[rand.Intn(len(notConnectedMessages))]
			return c.JSON(utils.StatsResponse{
				Success:           true,
				Message:           "stats for trash panda api :)",
				Requests:          randomMessage,
				Images:            utils.NUMBER_OF_IMAGES,
				Videos:            utils.NUMBER_OF_VIDEOS,
				TransparentImages: utils.NUMBER_OF_TRANS,
			})
		}
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
	v1Group.Get("/raccoon/transparent/:index", v1.GetRaccoonTransparentByIndex)
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
