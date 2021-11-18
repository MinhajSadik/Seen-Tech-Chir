package Responses

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Created(c *fiber.Ctx, resource string, data interface{}) {
	msg := resource + " has been created successfully!"
	c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": msg, "result": data})
}

func ModifiedSuccess(c *fiber.Ctx, resource string) {
	msg := resource + " has been modified successfully!"
	c.Status(fiber.StatusAccepted).JSON(fiber.Map{"success": true, "message": msg})
}

func Get(c *fiber.Ctx, resource string, data interface{}) {
	msg := resource + " has been retrieved successfully!"
	if fmt.Sprint(data) == "[]" {
		data = []interface{}{}
	}
	c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": msg, "result": data})
}

func ModifiedFail(c *fiber.Ctx, resource string, trace string) error {
	msg := resource + " has been not modified unfortunately!"
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": msg, "trace": trace})
}

func ResourceAlreadyExist(c *fiber.Ctx, resource string, data interface{}) error {
	msg := resource + " has not been saved because this name already exist!"
	return c.Status(fiber.StatusConflict).JSON(fiber.Map{"success": false, "message": msg, "result": data})
}

func NotFound(c *fiber.Ctx, resource string) error {
	msg := "Requested " + resource + " is not found!"
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": msg})
}

func NotValid(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": msg})
}

func ValidationError(c *fiber.Ctx, errs interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Validation error", "errors": errs})
}

func BadRequest(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": msg})
}

func SomethingGoneWrong(c *fiber.Ctx) error {
	msg := "Something gone wrong please try again later"
	return c.Status(fiber.StatusGone).JSON(fiber.Map{"success": false, "message": msg})
}

func Unauthorized(c *fiber.Ctx) error {
	// TODO after unifying the response, we can use the same response for all the cases
	return c.Status(fiber.StatusUnauthorized).Context().Err()
	msg := "Unauthorized request!"
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": msg})
}

func Unauthenticated(c *fiber.Ctx) error {
	// TODO after unifying the response, we can use the same response for all the cases
	return c.Status(fiber.StatusUnauthorized).Context().Err()
	msg := "You are unauthenticated, need to login first!"
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": msg})
}

func NotAllowed(c *fiber.Ctx) error {
	msg := "Not allowed request!"
	return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"success": false, "message": msg})
}

func SessionExpired(c *fiber.Ctx) error {
	msg := "Session expired!"
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": msg})
}

// Custom response
func Response(c *fiber.Ctx, statusCode int, success bool, msg string, data interface{}) {
	c.Status(statusCode).JSON(fiber.Map{"success": success, "message": msg, "result": data})
}
