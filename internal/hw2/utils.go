package hw2

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func getUserIDFromCtx(c *fiber.Ctx) (int64, error) {
	rawUserID := c.Params("userID")
	if rawUserID == "" {
		return 0, errEmptyUserID
	}

	userID, err := strconv.ParseInt(rawUserID, 10, 64)
	if err != nil {
		return 0, errIncorrectUserIDFormat
	}

	return userID, nil
}
