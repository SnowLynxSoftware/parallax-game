package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	mid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/snowlynxsoftware/parallax-game/config"
	"github.com/snowlynxsoftware/parallax-game/server/controllers"
	"github.com/snowlynxsoftware/parallax-game/server/database"
	"github.com/snowlynxsoftware/parallax-game/server/database/repositories"
	"github.com/snowlynxsoftware/parallax-game/server/middleware"
	"github.com/snowlynxsoftware/parallax-game/server/services"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

type AppServer struct {
	appConfig config.IAppConfig
	router    *chi.Mux
	dB        *database.AppDataSource
}

func NewAppServer(config config.IAppConfig) *AppServer {

	r := chi.NewRouter()
	r.Use(mid.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.GetCorsAllowedOrigin()},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-API-KEY"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	return &AppServer{
		appConfig: config,
		router:    r,
	}
}

func (s *AppServer) Start() {

	// Check if the app is running in production mode
	var isProductionMode = s.appConfig.GetCloudEnv() != "local"

	// Setup logger
	util.SetupZeroLogger(s.appConfig.IsDebugMode())

	// Connect to DB
	s.dB = database.NewAppDataSource()
	s.dB.Connect(s.appConfig.GetDBConnectionString())

	// Configure Repositories
	userRepository := repositories.NewUserRepository(s.dB)
	featureFlagRepository := repositories.NewFeatureFlagRepository(s.dB)

	// Configure Services
	featureFlagService := services.NewFeatureFlagService(featureFlagRepository)
	emailService := services.NewEmailService(s.appConfig.GetMJAPIKeyPublic(), s.appConfig.GetMJAPIKeyPrivate(), services.NewEmailTemplates())
	cryptoService := services.NewCryptoService(s.appConfig.GetAuthHashPepper())
	tokenService := services.NewTokenService(s.appConfig.GetJWTSecretKey())
	authService := services.NewAuthService(userRepository, tokenService, cryptoService, emailService, s.appConfig)
	userService := services.NewUserService(userRepository)
	templateService := services.NewTemplateService()
	staticService := services.NewStaticService()

	// Configure Middleware
	authMiddleware := middleware.NewAuthMiddleware(userRepository, tokenService, s.appConfig.GetSystemAPIKey())

	// Configure API Controllers (behind /api prefix)
	s.router.Mount("/api/health", controllers.NewHealthController().MapController())
	s.router.Mount("/api/auth", controllers.NewAuthController(authMiddleware, authService, isProductionMode, s.appConfig.GetCookieDomain()).MapController())
	s.router.Mount("/api/users", controllers.NewUserController(userService, authMiddleware).MapController())

	// Configure UI Controller (at root level)
	s.router.Mount("/", controllers.NewUIController(templateService, staticService, authMiddleware, featureFlagService).MapController())

	util.LogInfo("Starting server on localhost:3000")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", s.router))
}
