package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func CashRequestRoute(route fiber.Router) {
	route.Get("/get_all", Controllers.CashRequestGetAll)
	route.Get("/get_all_populated", Controllers.CashRequestGetAllPopulated)
	route.Post("/new", Controllers.CashRequestNew)
}
