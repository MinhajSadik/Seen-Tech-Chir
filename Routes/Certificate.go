/*
Author: omartarek9984
Code: tinder-014
*/
package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func CertificateRoute(route fiber.Router) {
	route.Put("/set_status/:id/:new_status", Controllers.CertificateSetStatus)
	route.Put("/modify/:id", Controllers.CertificateModify)
	route.Post("/get_all", Controllers.CertificateGetAll)
	route.Post("/new", Controllers.CertificateNew)
}
