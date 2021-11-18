package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func AttendanceRoute(route fiber.Router) {
	//route.Put("/set_status/:id/:new_status", Controllers.AttendanceSetStatus)
	//route.Put("/modify/:id", Controllers.AttendanceModify)
	//route.Post("/get_all_populated", Controllers.AttendanceGetAllPopulated)
	route.Get("/get_all", Controllers.AttendanceGetAll)
	route.Get("/get_all_populated", Controllers.AttendanceGetPopulated)
	route.Post("/new/:id/:type", Controllers.AttendanceNew)
	route.Get("/report/:id/:year/:month", Controllers.AttendaceReport)
}
