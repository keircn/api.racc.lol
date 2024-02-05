package v1

import (
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/venqoi/racc-api/utils"
)

func GetRaccFact(c *fiber.Ctx) error {
	factIndex := rand.Intn(len(utils.RaccoonFacts))

	return c.JSON(utils.Response{
		Success: true,
		Data: utils.FactStruct{
			Fact: utils.RaccoonFacts[factIndex],
		},
	})
}
