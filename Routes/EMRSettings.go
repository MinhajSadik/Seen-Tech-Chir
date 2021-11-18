package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func EMRSettingsRoute(route fiber.Router) {

	route.Post("/new", Controllers.EMRSettingsNew)
	route.Post("/getall", Controllers.EMRSettingsGetAll)
	route.Get("/get/:id", Controllers.EMRSettingsGetById)
	route.Post("/delete/:id", Controllers.EMRSettingsRemove)
	route.Put("/modify/:id", Controllers.EMRSettingsModify)
}
