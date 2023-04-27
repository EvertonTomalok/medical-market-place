package rest

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/EvertonTomalok/marketplace-health/docs"

	config "github.com/EvertonTomalok/marketplace-health/internal/app"

	"github.com/EvertonTomalok/marketplace-health/pkg/utils"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	router := gin.Default()

	injectRoutes(router)

	return router
}

// @title           Health Marketplace Api
// @version         0.0.1
// @description     Health Marketplace Api.

// @contact.name   Everton Tomalok
// @contact.email  evertontomalok123@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api
func RunServer(ctx context.Context, config config.Config) {
	done := utils.MakeDoneSignal()

	server := &http.Server{
		Addr:              net.JoinHostPort(config.App.Host, config.App.Port),
		Handler:           Router(),
		ReadTimeout:       15 * time.Minute,
		WriteTimeout:      15 * time.Minute,
		IdleTimeout:       30 * time.Minute,
		ReadHeaderTimeout: 30 * time.Minute,
	}

	go func() {
		log.Printf("Server started at %s:%s", config.App.Host, config.App.Port)

		if err := server.ListenAndServe(); err != nil {
			log.Panicf("Error trying to start server. %+v", err)
		}
	}()

	<-done
	log.Println("Stopping server...")
}

func injectRoutes(router *gin.Engine) {
	docs.SwaggerInfo.Title = "Health Marketplace Api"
	docs.SwaggerInfo.Description = "Health Marketplace Api."
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router.GET("/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	for _, route := range healthCheck {
		router.Handle(route.Method, route.Path, route.Handler)
	}

	apiGroup := router.Group("/api")
	for _, route := range routes {
		apiGroup.Handle(route.Method, route.Path, route.Handler)
	}
}
