package api

import (
	"net/http"

	"asset-core/internal/api/controller"
	"asset-core/internal/api/middleware"
	"asset-core/internal/config"
	"asset-core/internal/infrastructure/kafka"
	"asset-core/internal/module/asset"
	"asset-core/internal/module/data"
	"asset-core/internal/module/identity"
	"asset-core/internal/module/verification"
	"asset-core/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Dependencies struct {
	Config        *config.Config
	Logger        *logger.Logger
	DB            *gorm.DB
	Redis         *goredis.Client
	EventProducer kafka.Producer
}

func NewRouter(deps Dependencies) *gin.Engine {
	if deps.Config.App.Env == "prod" || deps.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.TraceID())
	r.Use(middleware.AccessLog(deps.Logger))
	r.Use(middleware.Recovery(deps.Logger))

	identityRepo := identity.NewRepository(deps.DB)
	identityService := identity.NewService(identityRepo, deps.EventProducer)

	assetRepo := asset.NewRepository(deps.DB)
	assetService := asset.NewService(assetRepo, identityService, deps.EventProducer)

	verificationRepo := verification.NewRepository(deps.DB)
	verificationService := verification.NewService(verificationRepo, assetRepo, deps.EventProducer)

	dataRepo := data.NewRepository(deps.DB)
	dataService := data.NewService(dataRepo, deps.EventProducer)

	assetCtl := controller.NewAssetController(assetService)
	identityCtl := controller.NewIdentityController(identityService)
	verificationCtl := controller.NewVerificationController(verificationService)
	dataCtl := controller.NewDataController(dataService)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": deps.Config.App.Name})
	})

	v1 := r.Group("/api/v1")
	{
		assets := v1.Group("/assets")
		{
			assets.POST("", assetCtl.Create)
			assets.GET("", assetCtl.List)
			assets.GET("/:id", assetCtl.Get)
			assets.PUT("/:id", assetCtl.Update)
			assets.DELETE("/:id", assetCtl.Delete)
			assets.POST("/:id/status", assetCtl.ChangeStatus)
			assets.GET("/:id/changes", assetCtl.ChangeLogs)
			assets.POST("/:id/verify", verificationCtl.VerifyAsset)
			assets.GET("/:id/verification-result", verificationCtl.LatestByAsset)
		}

		identities := v1.Group("/identities")
		{
			identities.POST("/generate", identityCtl.Generate)
			identities.GET("/:identity_id", identityCtl.Get)
			identities.POST("/:identity_id/bind", identityCtl.Bind)
			identities.POST("/:identity_id/unbind", identityCtl.Unbind)
			identities.GET("/:identity_id/features", identityCtl.Features)
		}

		verifications := v1.Group("/verifications")
		{
			verifications.POST("", verificationCtl.Create)
			verifications.GET("/:id", verificationCtl.Get)
		}

		dataRoutes := v1.Group("/data")
		{
			dataRoutes.POST("/import", dataCtl.CreateImportTask)
			dataRoutes.GET("/import-tasks", dataCtl.ListImportTasks)
			dataRoutes.GET("/import-tasks/:id", dataCtl.GetImportTask)
			dataRoutes.GET("/import-tasks/:id/errors", dataCtl.ImportErrors)
			dataRoutes.GET("/export/assets", dataCtl.ExportAssets)
		}
	}

	return r
}
