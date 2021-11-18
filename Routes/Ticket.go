package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func TicketRoute(route fiber.Router) {
	route.Put("/set_status/:id/:new_status", Controllers.TicketSetStatus)
	route.Put("/modify/:id", Controllers.TicketModify)
	route.Get("/get_all", Controllers.TicketGetAll)
	route.Get("/get_all_populated", Controllers.TicketGetAllPopulated)
	route.Get("/populated", Controllers.TicketPopulated)
	route.Post("/new", Controllers.TicketNew)
	route.Post("/ignore/:id", Controllers.TicketStatusIgnore)
	route.Post("/done/:id", Controllers.TicketStatusDone)
	route.Post("/selfassign/:id", Controllers.TicketAssign)
}
