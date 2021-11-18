package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func TrainingAllRequestsRoute(route fiber.Router) {
	route.Post("/get_all", Controllers.TrainingGetAll)
	route.Post("/get_all_populated", Controllers.TrainingGetAllPopulated)
	route.Post("/new", Controllers.TrainingRequestsNew)
	route.Post("/add_option/:id", Controllers.TrainingRequestsOptionNew)
	route.Put("/modify_option/:id/:optionId", Controllers.TraningRequestsOptionModify)
	route.Put("/modify_acceptno/:id", Controllers.TraningRequestsAcceptNoModify)
}
