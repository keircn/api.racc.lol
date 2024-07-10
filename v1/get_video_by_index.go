package v1

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/venqoi/racc-api/utils"
)

func GetRaccoonVideoByIndex(c *fiber.Ctx) error {
	var index = c.Params("index")
	var wantsJSON = utils.WantsJSON(c)

	parsedIndex, err := strconv.Atoi(index)
	if err != nil {
		return c.Status(400).JSON(utils.Response{
			Success: false,
			Message: "Invalid index parameter",
		})
	}

	supportedExtensions := []string{".mp4", ".mov"}

	var videoPath string
	var bytes []byte
	found := false
	for _, ext := range supportedExtensions {
		videoPath = fmt.Sprintf("./raccs/videos/racc%d%s", parsedIndex, ext)
		bytes, err = os.ReadFile(videoPath)
		if err == nil {
			found = true
			break
		}
	}

	c.Set("X-Raccoon-Video-Index", fmt.Sprint(parsedIndex))

	if !found {
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
				URL:   utils.BaseURL(c) + "/v1/video/" + index,
				Index: parsedIndex,
				Alt:   utils.GetAltv(parsedIndex),
			},
		})
	}

	contentType := "video/mp4"
	if filepath.Ext(videoPath) == ".mov" {
		contentType = "video/quicktime"
	}

	c.Set("Content-Type", contentType)
	return c.Send(bytes)
}
