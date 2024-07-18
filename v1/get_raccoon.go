package v1

import (
	"fmt"
	"image"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/venqoi/racc-api/utils"
)

func GetRaccoon(c *fiber.Ctx) error {
	var wantsJSON = utils.WantsJSON(c)
	var wantsTRANS = utils.WantsTRANS(c)
	randomIndex := utils.GetRandomIndex()
	randomIndexTrans := utils.GetRandomIndexTrans()
	var bytes []byte
	var err error

	if wantsTRANS { bytes, err = os.ReadFile("./raccs/transparent/racc" + fmt.Sprint(randomIndexTrans) + ".png")		
	} else { bytes, err = os.ReadFile("./raccs/racc" + fmt.Sprint(randomIndex) + ".jpg") }

	c.Set("X-Raccoon-Index", fmt.Sprint(randomIndex))

	if err != nil {
		println("error while reading racc photo", err.Error())
		if wantsJSON {
			return c.Status(500).JSON(utils.Response{
				Success: false,
				Message: "An error occurred whilst fetching file",
			})
		}

		return c.SendStatus(500)
	}

	if wantsJSON {
		var file *os.File
		var err error
		var url string

		if wantsTRANS { 
			file, err = os.Open("./raccs/racc/transparent" + fmt.Sprint(randomIndex) + ".png")		
			url = utils.BaseURL(c) + "/v1/raccoon/transparent/" + fmt.Sprint(randomIndexTrans)
		} else { 
			file, err = os.Open("./raccs/racc" + fmt.Sprint(randomIndex) + ".jpg")
			url = utils.BaseURL(c) + "/v1/raccoon/" + fmt.Sprint(randomIndex)
		}	

		if err != nil {
			println(err.Error())
		}

		defer file.Close()

		image, _, err := image.DecodeConfig(file)

		if err != nil {
			println(err.Error())
		}

		return c.JSON(utils.Response{
			Success: true,
			Data: utils.ImageStruct{
				URL:    url,
				Index:  randomIndex,
				Width:  image.Width,
				Height: image.Height,
				Alt:    utils.GetAlti(randomIndex),
			},
		})
	}

	c.Set("Content-Type", "image/jpeg")
	return c.Send(bytes)
}
