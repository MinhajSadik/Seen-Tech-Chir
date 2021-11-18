/* NewsFeed Module
code: tinder-002
author: omartarek9984
*/
package Routes

import (
	"SEEN-TECH-CHIR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func NewsFeedRoute(route fiber.Router) {
	route.Post("/get_all", Controllers.NewsFeedGetAll)
	route.Post("/new", Controllers.NewsFeedNew)
}
