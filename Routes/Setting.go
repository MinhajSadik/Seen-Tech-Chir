/*
Author: Omar Tarek
code: Tinder-005
*/
package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func SettingRoute(route fiber.Router) {
	route.Put("/modify/:id", Controllers.SettingModify)
	route.Post("/get_all", Controllers.SettingGetAll)
	route.Post("/new", Controllers.SettingNew)
	route.Post("/add_attachment", Controllers.SettingAddAttachment)
	route.Post("/toggle_status", Controllers.SettingToggleStatus)
}
