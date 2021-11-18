package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func TrainingRequestRoute(route fiber.Router) {

	route.Get("/get_all", Controllers.TrainingRequestGetAll)
	route.Get("/get_all_populated", Controllers.TrainingRequestGetAllPopulated)
	route.Post("/new", Controllers.TrainingRequestNew)
}
