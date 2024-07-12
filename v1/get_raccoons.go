package v1

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
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
	var wantsTRANS = utils.WantsTRANS(c)

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

		var filePath string
		var index = i
		var url string
		if random == "true" {
			index = rand.Intn(utils.NUMBER_OF_IMAGES-parsedFrom) + parsedFrom
		}

		if wantsTRANS {
			filePath = fmt.Sprintf("./raccs/transparent/racc%d.png", index)
			url = utils.BaseURL(c) + "/v1/raccoon/transparent/" + fmt.Sprint(index)
		} else {
			filePath = fmt.Sprintf("./raccs/racc%d.jpg", index)
			url = utils.BaseURL(c) + "/v1/raccoon/" + fmt.Sprint(index)
		}

		file, err := os.Open(filePath)
		if err != nil {
			println(err.Error())
			continue
		}

		imageConfig, _, err := image.DecodeConfig(file)
		file.Close()

		if err != nil {
			println(err.Error())
			continue
		}

		photos = append(photos, utils.ImageStruct{
			URL:    url,
			Index:  index,
			Width:  imageConfig.Width,
			Height: imageConfig.Height,
			Alt:    utils.GetAlti(index),
		})
	}

	return c.JSON(utils.Response{
		Success: true,
		Data:    photos,
	})
}