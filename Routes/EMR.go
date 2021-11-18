package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func EMRRoute(route fiber.Router) {

	route.Post("/new", Controllers.EMRNew)
	route.Get("/get_all_EMR_in_Department/:id", Controllers.EMRGetPendingInDepartment)
	route.Post("/change_EMR_status/:EMRId/:employeeId/:status", Controllers.ChangeEMRStatus)
}
