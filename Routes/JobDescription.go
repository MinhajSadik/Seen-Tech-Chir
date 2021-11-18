package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func JobDescriptionRoute(route fiber.Router) {
	route.Put("/set_status/:id/:new_status", Controllers.JobDescriptionSetStatus)
	route.Put("/modify/:id", Controllers.JobDescriptionModify)
	route.Post("/get_all", Controllers.JobDescriptionGetAll)
	route.Post("/get_all_populated", Controllers.JobDescriptionGetAllPopulated)
	route.Post("/new", Controllers.JobDescriptionNew)

	route.Get("/job_title/get_all", Controllers.JobTitleGetAll)
	// just array of job titles
	route.Get("/job_title/get", Controllers.JobTitleGetAllTitles)
	route.Put("/job_title/new/:id", Controllers.JobTitleNew)
	route.Put("/job_title/modify/:id", Controllers.JobTitleModify)
	route.Get("/job_title/get_id", Controllers.JobTitleGetNewId)
}
