package handlers

import (
	"errors"
	"time"

	"github.com/byvko-dev/am-stats-api/core/database"
	"github.com/byvko-dev/am-stats-api/logic"
	api "github.com/byvko-dev/am-types/api/generic/v1"
	stats "github.com/byvko-dev/am-types/api/stats/v1"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetPlayerSession(c *fiber.Ctx) error {
	var response api.ResponseWithError
	var request stats.RequestPayload
	err := c.BodyParser(&request)
	if err != nil {
		response.Error = api.ResponseError{
			Message: "Invalid request body",
			Context: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	if request.AccountID == 0 {
		response.Error = api.ResponseError{
			Message: "AccountID is required",
			Context: "accountId is missing",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	profile, session, snapshot, bigErr := logic.GetPlayerSession(request.AccountID, request.Days, false)
	if bigErr != nil {
		response.Error = api.ResponseError{
			Message: bigErr.Message,
			Context: bigErr.Raw.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	account, err := database.FindAccountByID(request.AccountID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		response.Error = api.ResponseError{
			Message: "ACCOUNT_CACHE_ERROR",
			Context: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	if account.AccountID == 0 {
		account.AccountID = request.AccountID
		account.Nickname = profile.Nickname
	}

	response.Data = stats.ResponsePayload{
		AccountID: account.AccountID,
		Timestamp: time.Now(),
		Snapshot:  snapshot,
		Session:   session,
		Account:   account,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
