package v1

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/venqoi/racc-api/utils"
)

func GetRaccoonVideo(c *fiber.Ctx) error {
	var wantsJSON = utils.WantsJSON(c)
	randomIndex := utils.GetRandomIndexVideo()

	bytes, err := os.ReadFile("./raccs/videos/racc" + fmt.Sprint(randomIndex) + ".mp4")

	c.Set("X-Raccoon-Video-Index", fmt.Sprint(randomIndex))

	if err != nil {
		println("error while reading racc video", err.Error())
		if wantsJSON {
			return c.Status(500).JSON(utils.Response{
				Success: false,
				Message: "An error occurred whilst fetching video file",
			})
		}

		return c.SendStatus(500)
	}

	if wantsJSON {
		return c.JSON(utils.Response{
			Success: true,
			Data: utils.VideoStruct{
				URL:   utils.BaseURL(c) + "/v1/video/" + fmt.Sprint(randomIndex),
				Index: randomIndex,
				Alt:   utils.GetAlti(randomIndex),
			},
		})
	}

	c.Set("Content-Type", "video/mp4")
	return c.Send(bytes)
}
