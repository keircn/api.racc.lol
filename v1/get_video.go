package v1

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/venqoi/racc-api/utils"
)

func GetRaccoonVideo(c *fiber.Ctx) error {
	var wantsJSON = utils.WantsJSON(c)
	randomIndex := utils.GetRandomIndexVideo()

	supportedExtensions := []string{".mp4", ".mov"}

	var videoPath string
	var bytes []byte
	var err error

	for _, ext := range supportedExtensions {
		videoPath = fmt.Sprintf("./raccs/videos/racc%d%s", randomIndex, ext)
		bytes, err = os.ReadFile(videoPath)
		if err == nil {
			break
		}
	}

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
				URL:   utils.BaseURL(c) + "/v1/video/" + filepath.Base(videoPath),
				Index: randomIndex,
				Alt:   utils.GetAlti(randomIndex),
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