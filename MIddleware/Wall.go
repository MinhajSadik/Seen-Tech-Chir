package middleware

import (
	Session "SEEN-TECH-CHIR/Auth"
	"SEEN-TECH-CHIR/Utils/Responses"
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

var protectedRoutes = map[string]bool{}
var managerRoutes = map[string]bool{}
var adminRoutes = map[string]bool{}
var publicAccessRegex = regexp.MustCompile(`[/images/*|/files/*]`)

var publicRoutes = map[string]bool{
	"/login": true,
}
var employeeRoutes = map[string]bool{
	"/employee/get_all": true,
}

func AppGaurd(c *fiber.Ctx) error {
	fmt.Println(c.OriginalURL())
	userAPIKey := string(c.Request().Header.Peek(Session.ApiKeyName))
	// Based on the user role, check if the route is protected
	if publicRoutes[c.OriginalURL()] || publicAccessRegex.MatchString(c.OriginalURL()) {
		return c.Next()
	}
	exist, session := Session.IsSessionExist(userAPIKey)
	if exist {
		if Session.IsSessionExpired(session) {
			Session.DeleteSession(session.ApiKey)
			return Responses.SessionExpired(c)
		} else {
			Session.UpdateSessionLastReqTime(session)
			if session.IsAdmin {
				// TODO Admin routes
				return c.Next()
			} else if session.IsManager {
				// TODO Manager routes
				return c.Next()
			} else {
				// TODO Employee routes
				return c.Next()
			}
		}
	}
	return Responses.Unauthenticated(c)
}
