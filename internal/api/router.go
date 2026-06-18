package api

import (
	"net/http"

	"asset-core/internal/api/controller"
	"asset-core/internal/api/middleware"
	"asset-core/internal/config"
	"asset-core/internal/infrastructure/elasticsearch"
	"asset-core/internal/infrastructure/kafka"
	"asset-core/internal/module/asset"
	"asset-core/internal/module/auth"
	"asset-core/internal/module/data"
	"asset-core/internal/module/identity"
	"asset-core/internal/module/verification"
	"asset-core/internal/module/workflow"
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

	esClient := elasticsearch.New(deps.Config.ES)

	assetRepo := asset.NewRepository(deps.DB)
	if err := assetRepo.AutoMigrate(); err != nil {
		deps.Logger.Fatal("asset bootstrap failed", logger.Error(err))
	}
	assetService := asset.NewService(assetRepo, identityService, deps.EventProducer, esClient)

	verificationRepo := verification.NewRepository(deps.DB)
	verificationService := verification.NewService(verificationRepo, assetService, deps.EventProducer)

	dataRepo := data.NewRepository(deps.DB)
	dataService := data.NewService(dataRepo, assetService, deps.EventProducer)

	authRepo := auth.NewRepository(deps.DB)
	authService := auth.NewService(authRepo, deps.Config.JWT)
	if err := authService.Bootstrap(); err != nil {
		deps.Logger.Fatal("auth bootstrap failed", logger.Error(err))
	}
	workflowRepo := workflow.NewRepository(deps.DB)
	workflowService := workflow.NewService(workflowRepo, assetService)
	if err := workflowService.Bootstrap(); err != nil {
		deps.Logger.Fatal("workflow bootstrap failed", logger.Error(err))
	}

	assetCtl := controller.NewAssetController(assetService)
	identityCtl := controller.NewIdentityController(identityService)
	verificationCtl := controller.NewVerificationController(verificationService)
	dataCtl := controller.NewDataController(dataService)
	authCtl := controller.NewAuthController(authService)
	workflowCtl := controller.NewWorkflowController(workflowService)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": deps.Config.App.Name})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/login", authCtl.Login)

		authRoutes := v1.Group("")
		authRoutes.Use(middleware.AuthRequired(authService))
		authRoutes.GET("/auth/me", authCtl.Me)
		authRoutes.POST("/auth/register", middleware.RequirePermission(auth.PermissionUserCreate), authCtl.Register)
		authRoutes.GET("/auth/users", middleware.RequirePermission(auth.PermissionUserCreate), authCtl.ListUsers)
		authRoutes.PUT("/auth/users/:id", middleware.RequirePermission(auth.PermissionUserCreate), authCtl.UpdateUser)
		authRoutes.GET("/auth/permissions", middleware.RequirePermission(auth.PermissionUserCreate), authCtl.ListPermissions)

		assets := v1.Group("/assets")
		assets.Use(middleware.AuthRequired(authService))
		{
			assets.POST("", middleware.RequirePermission(auth.PermissionAssetCreate), assetCtl.Create)
			assets.GET("", middleware.RequirePermission(auth.PermissionAssetRead), assetCtl.List)
			assets.GET("/:id", middleware.RequirePermission(auth.PermissionAssetRead), assetCtl.Get)
			assets.PUT("/:id", middleware.RequirePermission(auth.PermissionAssetUpdate), assetCtl.Update)
			assets.POST("/:id/identity", middleware.RequirePermission(auth.PermissionAssetUpdate), assetCtl.GenerateIdentity)
			assets.DELETE("/:id", middleware.RequirePermission(auth.PermissionAssetDelete), assetCtl.Delete)
			assets.POST("/:id/status", middleware.RequirePermission(auth.PermissionAssetUpdate), assetCtl.ChangeStatus)
			assets.GET("/:id/changes", middleware.RequirePermission(auth.PermissionAssetRead), assetCtl.ChangeLogs)
			assets.POST("/:id/insurance", middleware.RequirePermission(auth.PermissionAssetUpdate), assetCtl.AddInsurance)
			assets.GET("/:id/insurance", middleware.RequirePermission(auth.PermissionAssetRead), assetCtl.ListInsurance)
			assets.POST("/:id/impairments", middleware.RequirePermission(auth.PermissionWorkflowApprove), assetCtl.RecordImpairment)
			assets.GET("/:id/impairments", middleware.RequirePermission(auth.PermissionAssetRead), assetCtl.ListImpairments)
			assets.POST("/:id/verify", middleware.RequirePermission(auth.PermissionAssetUpdate), verificationCtl.VerifyAsset)
			assets.GET("/:id/verification-result", middleware.RequirePermission(auth.PermissionAssetRead), verificationCtl.LatestByAsset)
		}

		identities := v1.Group("/identities")
		identities.Use(middleware.AuthRequired(authService))
		{
			identities.POST("/generate", identityCtl.Generate)
			identities.GET("", identityCtl.List)
			identities.GET("/:identity_id", identityCtl.Get)
			identities.POST("/:identity_id/bind", identityCtl.Bind)
			identities.POST("/:identity_id/unbind", identityCtl.Unbind)
			identities.GET("/:identity_id/features", identityCtl.Features)
		}

		verifications := v1.Group("/verifications")
		verifications.Use(middleware.AuthRequired(authService))
		{
			verifications.POST("", verificationCtl.Create)
			verifications.GET("", verificationCtl.List)
			verifications.GET("/:id", verificationCtl.Get)
		}

		dataRoutes := v1.Group("/data")
		dataRoutes.Use(middleware.AuthRequired(authService))
		{
			dataRoutes.POST("/import", middleware.RequirePermission(auth.PermissionAssetCreate), dataCtl.CreateImportTask)
			dataRoutes.POST("/import/assets", middleware.RequirePermission(auth.PermissionAssetCreate), dataCtl.ImportAssetsExcel)
			dataRoutes.GET("/import-tasks", dataCtl.ListImportTasks)
			dataRoutes.GET("/import-tasks/:id", dataCtl.GetImportTask)
			dataRoutes.GET("/import-tasks/:id/errors", dataCtl.ImportErrors)
			dataRoutes.GET("/export/assets", dataCtl.ExportAssets)
		}

		workflowRoutes := v1.Group("/workflows")
		workflowRoutes.Use(middleware.AuthRequired(authService))
		{
			workflowRoutes.GET("/definitions", middleware.RequirePermission(auth.PermissionWorkflowStart), workflowCtl.ListDefinitions)
			workflowRoutes.PUT("/definitions", middleware.RequirePermission(auth.PermissionWorkflowConfig), workflowCtl.SaveDefinition)
			workflowRoutes.DELETE("/definitions/:id", middleware.RequirePermission(auth.PermissionWorkflowConfig), workflowCtl.DeleteDefinition)
			workflowRoutes.POST("/instances", middleware.RequirePermission(auth.PermissionWorkflowStart), workflowCtl.Start)
			workflowRoutes.GET("/instances", middleware.RequirePermission(auth.PermissionWorkflowStart), workflowCtl.ListInstances)
			workflowRoutes.GET("/tasks", middleware.RequirePermission(auth.PermissionWorkflowApprove), workflowCtl.ListTasks)
			workflowRoutes.POST("/tasks/:id/approve", middleware.RequirePermission(auth.PermissionWorkflowApprove), workflowCtl.Approve)
		}
	}

	return r
}
