package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func EmployeeRoute(route fiber.Router) {
	route.Put("/set_status/:id/:new_status", Controllers.EmployeeSetStatus)
	route.Put("/modify/:id", Controllers.EmployeeModify)
	route.Put("/modify_info/:id", Controllers.EmployeeModifyInfo)
	route.Get("/get_all", Controllers.EmployeeGetAll)
	route.Get("/get_all_location_status/:id", Controllers.EmployeeGetAllLocationStatus)
	route.Post("/get_all_populated", Controllers.EmployeeGetAllPopulated)
	route.Post("/new", Controllers.EmployeeNew)
	route.Post("/add_working_info/:id", Controllers.EmployeeAddWorkingInfo)
	route.Post("/add_image_attachment/:id", Controllers.EmployeeAddImageAttachments)
	route.Post("/image", Controllers.UploadImage)
	route.Get("/get_pending_requests", Controllers.GetAllPendingRequests)
	route.Post("/change_request_status/:id/:type/:rid/:status", Controllers.EmployeeChangeRequestStatus)
	route.Post("/receivings/:id", Controllers.EmployeeSetRecivings)
	// Additional contacts endpoints
	route.Post("/additionalcontact/:id", Controllers.EmployeeAddContact)
	route.Put("/additionalcontact/:id/:cid", Controllers.EmployeeEditContact)
	route.Delete("/additionalcontact/:id/:cid", Controllers.EmployeeRemoveContact)
	route.Get("/get_by_department/:departmentid", Controllers.EmployeeGetAllByDepartment)
	route.Post("/:id/:type", Controllers.EmployeeAddDuration)
	route.Get("/:id/:type", Controllers.EmployeeGetDurations)
}
