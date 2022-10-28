package handlers

import (
	"time"

	"github.com/byvko-dev/am-stats-api/logic"
	api "github.com/byvko-dev/am-types/api/generic/v1"
	stats "github.com/byvko-dev/am-types/api/stats/v1"
	"github.com/gofiber/fiber/v2"
)

func GetPlayerSession(c *fiber.Ctx) error {
	var response api.ResponseWithError
	var request stats.RequestPayload
	err := c.BodyParser(&request)
	if err != nil {
		response.Error = api.ResponseError{
			Message: "Invalid request body",
			Context: err.Error(),
			Code:    "INVALID_REQUEST",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	if request.AccountID == 0 {
		response.Error = api.ResponseError{
			Message: "AccountID is required",
			Context: "accountId is missing",
			Code:    "INVALID_REQUEST",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	profile, session, snapshot, bigErr := logic.GetPlayerSession(request.AccountID, request.Days, false)
	if bigErr != nil {
		response.Error = api.ResponseError{
			Code:    "SESSION_REQUEST_FAILED",
			Message: bigErr.Message,
			Context: bigErr.Raw.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Add player data from cache or cache it
	profile = logic.CheckPlayerCache(profile)

	response.Data = stats.ResponsePayload{
		AccountID: profile.AccountID,
		Timestamp: time.Now(),
		Snapshot:  snapshot,
		Session:   session,
		Account:   profile,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
