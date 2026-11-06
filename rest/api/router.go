package api

import (
	"Quote/rest"
	"Quote/rest/api/controllers"
	"Quote/rest/api/controllers/apertureController"
	"Quote/rest/api/controllers/clientController"
	"Quote/rest/api/controllers/colorController"
	"Quote/rest/api/controllers/fileController"
	"Quote/rest/api/controllers/glassController"
	"Quote/rest/api/controllers/profileController"
	"Quote/rest/api/controllers/quoteController"
	"Quote/rest/api/controllers/supplierController"
	"Quote/rest/api/controllers/userController"
	"Quote/rest/middleware"

	"github.com/gin-gonic/gin"
)

func NewPublicAPIRouter() *gin.Engine {
	router := rest.NewRouter()

	publicGroup := router.Group("/public")
	{
		userGroup := publicGroup.Group("/users")
		{
			userGroup.POST("/create", userController.Create)
			userGroup.PATCH("/:user_id", userController.Update)
			userGroup.DELETE("/soft/:user_id", userController.SoftDelete)
			userGroup.GET("/", userController.Get)
			userGroup.GET("/by/id/:user_id", userController.GetById)
		}
	}

	developmentGroup := router.Group("/dev")
	{
		developmentGroup.GET("/test", controllers.Test)
	}

	basicAuthGroup := router.Group("/basic", middleware.AuthMiddleware)
	{
		clientGroup := basicAuthGroup.Group("/clients")
		{
			clientGroup.POST("/create", clientController.Create)
			clientGroup.PATCH("/:client_id", clientController.Update)
			clientGroup.DELETE("/soft/:client_id", clientController.SoftDelete)
			clientGroup.GET("/", clientController.Get)
			clientGroup.GET("/by/id/:client_id", clientController.GetById)
		}
		supplierGroup := basicAuthGroup.Group("/suppliers")
		{
			supplierGroup.POST("/create", supplierController.Create)
			supplierGroup.PATCH("/:supplier_id", supplierController.Update)
			supplierGroup.DELETE("/soft/:supplier_id", supplierController.SoftDelete)
			supplierGroup.GET("/", supplierController.Get)
			supplierGroup.GET("/by/id/:supplier_id", supplierController.GetById)
		}
		colorGroup := basicAuthGroup.Group("/colors")
		{
			colorGroup.POST("/create", colorController.Create)
			colorGroup.PATCH("/:color_id", colorController.Update)
			colorGroup.DELETE("/soft/:color_id", colorController.SoftDelete)
			colorGroup.GET("/", colorController.Get)
			colorGroup.GET("/by/id/:color_id", colorController.GetById)
		}
		glassGroup := basicAuthGroup.Group("/glasses")
		{
			glassGroup.POST("/create", glassController.Create)
			glassGroup.PATCH("/:glass_id", glassController.Update)
			glassGroup.DELETE("/soft/:glass_id", glassController.SoftDelete)
			glassGroup.GET("/", glassController.Get)
			glassGroup.GET("/by/id/:glass_id", glassController.GetById)
		}
		profileGroup := basicAuthGroup.Group("/profiles")
		{
			profileGroup.POST("/create", profileController.Create)
			profileGroup.PATCH("/:profile_id", profileController.Update)
			profileGroup.DELETE("/soft/:profile_id", profileController.SoftDelete)
			profileGroup.GET("/", profileController.Get)
			profileGroup.GET("/by/id/:profile_id", profileController.GetById)

			priceMatrixGroup := profileGroup.Group("/price-matrix")
			{
				priceMatrixGroup.GET("/max-dimensions/:profile_id/:width", profileController.GetMaxDimensionsInPriceMatrix)
			}
		}
		quoteGroup := basicAuthGroup.Group("/quotes")
		{
			quoteGroup.POST("/create", quoteController.Create)
			quoteGroup.PATCH("/:quote_id", quoteController.Update)
			quoteGroup.DELETE("/soft/:quote_id", quoteController.SoftDelete)
			quoteGroup.GET("/", quoteController.Get)
			quoteGroup.GET("/by/id/:quote_id", quoteController.GetById)

			calculationsGroup := quoteGroup.Group("/calculations")
			{
				calculationsGroup.GET("/:quote_id", quoteController.GetCalculationsById)
			}

			datasheetsGroup := quoteGroup.Group("/datasheets")
			{
				datasheetsGroup.GET("/:quote_id", quoteController.GetDatasheets)
			}

			invoicesGroup := quoteGroup.Group("/invoices")
			{
				invoicesGroup.GET("/download/client/:quote_id", quoteController.GenerateClientInvoice)
				invoicesGroup.GET("/download/supplier/:quote_id", quoteController.GenerateSupplierInvoice)
			}
		}
		apertureGroup := basicAuthGroup.Group("/apertures")
		{
			apertureGroup.POST("/create", apertureController.Create)
			apertureGroup.PATCH("/:aperture_id", apertureController.Update)
			apertureGroup.DELETE("/soft/:aperture_id", apertureController.SoftDelete)
			apertureGroup.GET("/", apertureController.Get)
			apertureGroup.GET("/by/id/:aperture_id", apertureController.GetById)
			apertureGroup.GET("/copy/:aperture_id", apertureController.Copy)
			apertureGroup.POST("/duplicate/:aperture_id", apertureController.Duplicate)
		}
		fileGroup := basicAuthGroup.Group("/files")
		{
			fileGroup.POST("/upload/single", fileController.AddFile)
			fileGroup.GET("/by/id/:file_id", fileController.GetById)
			fileGroup.GET("/download/single/:file_id", fileController.DownloadFile)
		}

		userGroup := basicAuthGroup.Group("/users")
		{
			userGroup.POST("/login", userController.Login)
		}
	}

	return router
}
