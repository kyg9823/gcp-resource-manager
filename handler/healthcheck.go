package handler

import (
	"github.com/gofiber/fiber/v2"
)

func Healthcheck(ctx *fiber.Ctx) error {

	return ctx.JSON(fiber.Map{
		"status":  200,
		"message": "OK",
	})

	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "application/json")

	// result := make(map[string]interface{})
	// result["StatusCode"] = http.StatusOK
	// result["Message"] = "OK"

	// jsonResult, err := json.Marshal(result)
	// if err != nil {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// }
	// w.Write(jsonResult)
}
