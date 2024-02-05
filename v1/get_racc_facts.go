package v1

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/venqoi/racc-api/utils"
)

func GetRaccFacts(c *fiber.Ctx) error {
	var from = c.Query("from")
	var take = c.Query("take")

	var facts []string = make([]string, 0)

	if len(from) == 0 {
		from = "0"
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

	for i := 0 + parsedFrom; i < parsedTake+parsedFrom && i < len(utils.RaccoonFacts); i++ {
		facts = append(facts, utils.RaccoonFacts[i])
	}

	return c.JSON(utils.Response{
		Success: true,
		Data:    facts,
	})
}
