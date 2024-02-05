package v1

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"math/rand"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/venqoi/racc-api/utils"
)

func GetRaccoons(c *fiber.Ctx) error {
	var from = c.Query("from")
	var take = c.Query("take")
	var random = c.Query("random")

	if len(from) == 0 {
		from = "1"
	}

	if len(take) == 0 {
		take = "25"
	}

	parsedTake, err := strconv.Atoi(take)
	if err != nil {
		return c.Status(500).JSON(utils.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	parsedFrom, err := strconv.Atoi(from)
	if err != nil {
		return c.Status(500).JSON(utils.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	var photos []utils.ImageStruct
	for i := 0 + parsedFrom; i < parsedTake+parsedFrom && i < utils.NUMBER_OF_IMAGES; i++ {

		/* if user wants random index */
		var index = i
		if random == "true" {
			index = rand.Intn(utils.NUMBER_OF_IMAGES-parsedFrom) + parsedFrom
		}

		file, err := os.Open("./raccs/racc" + fmt.Sprint(index) + ".jpg")

		if err != nil {
			println(err.Error())
		}

		image, _, err := image.DecodeConfig(file)

		if err != nil {
			println(err.Error())
		}

		photos = append(photos, utils.ImageStruct{
			URL:    utils.BaseURL(c) + "/v1/raccoon/" + fmt.Sprint(index),
			Index:  index,
			Width:  image.Width,
			Height: image.Height,
			Alt:    utils.GetAlti(index),
		})

		file.Close()
	}

	return c.JSON(utils.Response{
		Success: true,
		Data:    photos,
	})
}
