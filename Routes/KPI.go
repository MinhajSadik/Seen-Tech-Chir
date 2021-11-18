/* KPI Module
code: tinder-003
author: rrrokhtar
*/
package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func KPIRoute(route fiber.Router) {
	route.Post("/new", Controllers.KPINew)
	route.Get("/get_all", Controllers.KPIGet)
	route.Put("/modify/:id", Controllers.KPIModify)
}
