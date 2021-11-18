/*
Author: omartarek9984
Code: tinder-006
*/
package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func LocationRoute(route fiber.Router) {
	route.Put("/set_status/:id/:new_status", Controllers.LocationSetStatus)
	route.Put("/modify/:id", Controllers.LocationModify)
	route.Post("/get_all", Controllers.LocationGetAll)
	route.Post("/new", Controllers.LocationNew)
	route.Get("/allowed_employees/:id", Controllers.LocationGetEmployees)
}
