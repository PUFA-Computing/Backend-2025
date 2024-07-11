package api

import (
	"Backend/configs"
	"Backend/internal/handlers/aspirations"
	"Backend/internal/handlers/auth"
	"Backend/internal/handlers/event"
	"Backend/internal/handlers/news"
	"Backend/internal/handlers/permission"
	"Backend/internal/handlers/role"
	"Backend/internal/handlers/user"
	"Backend/internal/handlers/version"
	"Backend/internal/middleware"
	"Backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://computing.president.ac.id", "https://staging.computing.president.ac.id", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	maxTokens := 1000
	refillInterval := time.Minute
	r.Use(middleware.RateLimiterMiddleware(maxTokens, refillInterval, "general"))

	r.Static("/public", "./public")

	authService := services.NewAuthService()
	userService := services.NewUserService()
	eventService := services.NewEventService()
	newsService := services.NewNewsService()
	roleService := services.NewRoleService()
	permissionService := services.NewPermissionService()
	aspirationsService := services.NewAspirationService()
	AWSService, _ := services.NewAWSService()
	R2Service, _ := services.NewR2Service()
	MailgunService := services.NewMailgunService(
		configs.LoadConfig().MailGunDomain,
		configs.LoadConfig().MailGunApiKey,
		configs.LoadConfig().MailGunSenderEmail,
	)
	VersionService := services.NewVersionService(configs.LoadConfig().GithubAccessToken)

	eventStatusUpdater := services.NewEventStatusUpdater(eventService)
	go eventStatusUpdater.Run()

	versionUpdater := services.NewVersionUpdater(VersionService)
	go versionUpdater.Run()

	authHandlers := auth.NewAuthHandlers(authService, permissionService, MailgunService, userService)
	userHandlers := user.NewUserHandlers(userService, permissionService, AWSService, R2Service)
	eventHandlers := event.NewEventHandlers(eventService, permissionService, AWSService, R2Service)
	newsHandlers := news.NewNewsHandler(newsService, permissionService, AWSService, R2Service)
	roleHandlers := role.NewRoleHandler(roleService, userService, permissionService)
	permissionHandlers := permission.NewPermissionHandler(permissionService)
	aspirationHandlers := aspirations.NewAspirationHandlers(aspirationsService, permissionService)
	versionHandlers := version.NewVersionHandlers(VersionService)

	api := r.Group("/api/v1")

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandlers.RegisterUser)
		authRoutes.POST("/login", middleware.RateLimiterMiddleware(5, time.Minute, "login"), authHandlers.Login)
		authRoutes.POST("/logout", authHandlers.Logout)
		authRoutes.POST("/refresh-token", middleware.TokenMiddleware(), authHandlers.RefreshToken)
		authRoutes.GET("/verify-email", authHandlers.VerifyEmail)
	}

	userRoutes := api.Group("/user")
	{
		userRoutes.GET("/list", userHandlers.ListUsers)

		userRoutes.Use(middleware.TokenMiddleware())
		userRoutes.GET("/:userID", userHandlers.GetUserByID)
		userRoutes.PUT("/edit", userHandlers.EditUser)
		userRoutes.DELETE("/delete", userHandlers.DeleteUser)
		userRoutes.POST("/upload-profile-picture", userHandlers.UploadProfilePicture)
		userRoutes.PUT("/update-user", userHandlers.AdminUpdateRoleAndStudentIDVerified)
		userRoutes.POST("/2fa/enable", userHandlers.EnableTwoFA)
		userRoutes.POST("/2fa/verify", userHandlers.VerifyTwoFA)
		userRoutes.POST("/2fa/disable", userHandlers.DisableTwoFA)

		// ListEventsRegisteredByUser
		userRoutes.GET("/registered-events", eventHandlers.ListEventsRegisteredByUser)
	}

	eventRoutes := api.Group("/event")
	{
		eventRoutes.GET("/:eventID", eventHandlers.GetEventBySlug)
		eventRoutes.GET("/", eventHandlers.ListEvents)
		eventRoutes.GET("/:eventID/total-participant", eventHandlers.TotalRegisteredUsers)
		eventRoutes.Use(middleware.TokenMiddleware())
		eventRoutes.POST("/create", eventHandlers.CreateEvent)
		eventRoutes.PATCH("/:eventID/edit", eventHandlers.EditEvent)
		eventRoutes.DELETE("/:eventID/delete", eventHandlers.DeleteEvent)
		eventRoutes.POST("/:eventID/register", eventHandlers.RegisterForEvent)
		eventRoutes.GET("/:eventID/registered-users", eventHandlers.ListRegisteredUsers)
	}

	newsRoutes := api.Group("/news")
	{
		newsRoutes.GET("/", newsHandlers.ListNews)
		newsRoutes.GET("/:newsID", newsHandlers.GetNewsBySlug)
		newsRoutes.Use(middleware.TokenMiddleware())
		newsRoutes.POST("/create", newsHandlers.CreateNews)
		newsRoutes.PUT("/:newsID/edit", newsHandlers.EditNews)
		newsRoutes.DELETE("/:newsID/delete", newsHandlers.DeleteNews)
		newsRoutes.POST("/:newsID/like", newsHandlers.LikeNews)
	}

	roleRoutes := api.Group("/roles")
	{
		roleRoutes.Use(middleware.TokenMiddleware())
		roleRoutes.GET("/", roleHandlers.ListRoles)
		roleRoutes.POST("/create", roleHandlers.CreateRole)
		roleRoutes.GET("/:roleID", roleHandlers.GetRoleByID)
		roleRoutes.PUT("/:roleID/edit", roleHandlers.EditRole)
		roleRoutes.DELETE("/:roleID/delete", roleHandlers.DeleteRole)
		roleRoutes.POST("/:roleID/assign/:userID", roleHandlers.AssignRoleToUser)
	}
	permissionRoutes := api.Group("/permissions")
	{
		permissionRoutes.Use(middleware.TokenMiddleware())
		permissionRoutes.GET("/list", permissionHandlers.ListPermissions)
		permissionRoutes.POST("/assign/:roleID", permissionHandlers.AssignPermissionToRole)

	}

	aspirationRoutes := api.Group("/aspirations")
	{
		aspirationRoutes.GET("/", aspirationHandlers.GetAspirations)
		aspirationRoutes.GET("/:id", aspirationHandlers.GetAspirationByID)
		aspirationRoutes.Use(middleware.TokenMiddleware())
		aspirationRoutes.POST("/create", aspirationHandlers.CreateAspiration)
		aspirationRoutes.PATCH("/:id/close", aspirationHandlers.CloseAspiration)
		aspirationRoutes.DELETE("/:id/delete", aspirationHandlers.DeleteAspiration)
		aspirationRoutes.POST("/:id/upvote", aspirationHandlers.UpvoteAspiration)
		aspirationRoutes.GET("/:id/get_upvotes", aspirationHandlers.GetUpvotesByAspirationID)
		aspirationRoutes.POST("/:id/admin_reply", aspirationHandlers.AddAdminReply)
	}

	versionRoutes := api.Group("/version")
	{
		versionRoutes.GET("/", versionHandlers.GetVersion)
		versionRoutes.GET("/changelog", versionHandlers.GetChangelog)
	}

	return r
}
