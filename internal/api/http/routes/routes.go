package routes

import (
	"spy-cat-agency/internal/api/http/handlers"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func SetupRoutes(e *echo.Echo, catHandler *handlers.CatHandler, missionHandler *handlers.MissionHandler) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api/v1")

	api.GET("/cats", catHandler.ListCats)
	api.GET("/cats/:id", catHandler.GetCat)
	api.GET("/cats/breeds", catHandler.GetBreeds)

	agency := api.Group("/agency")

	agencyCats := agency.Group("/cats")
	agencyCats.POST("", catHandler.CreateCat)
	agencyCats.PUT("/:id/salary", catHandler.UpdateCatSalary)
	agencyCats.DELETE("/:id", catHandler.DeleteCat)

	agencyMissions := agency.Group("/missions")
	agencyMissions.POST("", missionHandler.CreateMission)
	agencyMissions.GET("", missionHandler.ListMissions)
	agencyMissions.GET("/:id", missionHandler.GetMission)
	agencyMissions.DELETE("/:id", missionHandler.DeleteMission)
	agencyMissions.POST("/:id/assign", missionHandler.AssignCatToMission)
	agencyMissions.GET("/free-cats", missionHandler.GetFreeCats)

	agencyMissions.POST("/:id/targets", missionHandler.AddTargetToMission)
	agencyMissions.DELETE("/:id/targets/:targetId", missionHandler.DeleteTargetFromMission)

	spyCats := api.Group("/spy-cats/:catId")
	spyCats.GET("/mission", missionHandler.GetCatMission)
	spyCats.PUT("/mission/targets/:targetId/status", missionHandler.UpdateTargetStatus)
	spyCats.PUT("/mission/targets/:targetId/notes", missionHandler.UpdateTargetNotes)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
}
