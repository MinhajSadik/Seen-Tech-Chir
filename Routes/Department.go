package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func DepartmentRoute(route fiber.Router) {
	route.Get("/", Controllers.DepartmentGetAll)
	route.Post("/new", Controllers.DepartmentNew)
	route.Post("/get_all", Controllers.DepartmentGetAll)
	route.Post("/get_all/populated", Controllers.DepartmentGetAllPopulated)
	route.Post("/get_parents_options", Controllers.DepartmentGetParentsOptions)
	route.Put("/modify", Controllers.DepartmentModify)
	route.Get("/employees/:id", Controllers.DepartmentGetEmployees)
	route.Get("/tree", Controllers.DepartmentTree)
	route.Get("/map", Controllers.DepartmentMap)
	route.Get("/access/:id", Controllers.DepartmentGetAccessibleChilds)
}
