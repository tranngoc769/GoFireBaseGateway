package api

import (
	"net/http"
	"time"

	mdw "go-firebase-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	serviceName = "go-firebase-gateway-api"
	version     = "v1.0"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer() *Server {
	engine := gin.New()
	// engine.Use(cors.Default())
	engine.Use(gin.Recovery())
	// engine.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{"*"},
	// 	AllowMethods: []string{"*"},
	// }))
	engine.Use(CORSMiddleware())
	mdw.SetupGoGuardian()
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": serviceName,
			"version": version,
			"time":    time.Now().Unix(),
		})
	})
	server := &Server{Engine: engine}
	return server
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (server *Server) Start(port string) {
	v := make(chan struct{})
	go func() {
		if err := server.Engine.Run(":" + port); err != nil {
			log.WithError(err).Error("failed to start service")
			close(v)
		}
	}()
	log.Infof("service %v listening on port %v", serviceName, port)
	<-v
}
