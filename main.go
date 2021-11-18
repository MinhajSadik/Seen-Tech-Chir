package main

import (
	auth "SEEN-TECH-CHIR/Auth"
	"SEEN-TECH-CHIR/DBManager"
	middleware "SEEN-TECH-CHIR/MIddleware"
	"SEEN-TECH-CHIR/Routes"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func SetupRoutes(app *fiber.App) {
	Routes.DepartmentRoute(app.Group("/department"))
	Routes.KPIRoute(app.Group("/kpi"))
	Routes.EMRRoute(app.Group("/emr"))
	Routes.JobDescriptionRoute(app.Group("/job_description"))
	Routes.EmployeeRoute(app.Group("/employee"))
	Routes.NewsFeedRoute(app.Group("/news_feed"))
	Routes.SettingRoute(app.Group("/setting"))
	Routes.LocationRoute(app.Group("/location"))
	Routes.TicketRoute(app.Group("/ticket"))
	Routes.AttendanceRoute(app.Group("/attendance"))
	Routes.CertificateRoute(app.Group("/certificate"))
	Routes.EMRSettingsRoute(app.Group("/emr_settings"))
	Routes.TrainingRequestRoute(app.Group("/training_request"))
	Routes.TrainingAllRequestsRoute(app.Group("/training_all_requests"))
	Routes.CashRequestRoute(app.Group("/cash_request"))
	app.Post("/login", auth.Login)
	app.Post("/logout", auth.Logout)
	app.Get("/init", auth.SeedAdmin)
}

func main() {
	fmt.Println(("Hello SEEN-TECH-CHIR"))
	fmt.Print("Initializing Database Connection ... ")
	initState := DBManager.InitCollections()
	if initState {
		fmt.Println("[OK]")
	} else {
		fmt.Println("[FAILED]")
		return
	}

	fmt.Print("Initializing the server ... ")
	app := fiber.New()
	app.Use(cors.New())
	app.Use(pprof.New())
	app.Use(middleware.AppGaurd)
	SetupRoutes(app)
	app.Static("/", "./public")
	fmt.Println("[OK]")
	app.Listen(":8080")
}
