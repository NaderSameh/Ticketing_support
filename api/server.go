package api

import (
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
	_ "github.com/naderSameh/ticketing_support/docs"
	worker "github.com/naderSameh/ticketing_support/worker"
	"github.com/penglongli/gin-metrics/ginmetrics"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type Server struct {
	store           db.Store
	router          *gin.Engine
	taskDistributor worker.TaskDistributor
}

func NewServer(store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {

	server := &Server{
		store:           store,
		taskDistributor: taskDistributor,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	// For debug mode
	// router := gin.Default()

	//For release mode
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	SetupCORS(router)
	SetupLogger(router)
	SetupMetrics(router)
	//tickets
	router.POST("/tickets", server.createTicket)              // Create new ticket
	router.DELETE("/tickets/:ticket_id", server.deleteTicket) // Delete ticket with its comments
	router.GET("/tickets", server.listTicket)                 // Get list of tickets
	router.GET("/tickets/:ticket_id", server.getTicket)       // Get single ticket
	router.PUT("/tickets/:ticket_id", server.updateTicket)    //Assign ticket or update its status
	//comments
	router.DELETE("/tickets/:ticket_id/comments/:comment_id", server.deleteComment) // Delete a comment
	router.PUT("/tickets/:ticket_id/comments/:comment_id", server.updateComment)    //Edit comment text
	router.GET("/tickets/:ticket_id/comments", server.listComments)                 // Get comments for a ticket
	router.POST("/tickets/:ticket_id/comments", server.createComment)               //Add comment to a ticket
	//caterogries
	router.GET("tickets/categories", server.listCategories)
	router.POST("tickets/categories", server.createCategory) // Create Category
	//Swagger
	url := ginSwagger.URL("http://localhost:8080/tickets/swagger/doc.json") // The url pointing to API definition
	router.GET("/tickets/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	server.router = router

}

func SetupMetrics(router *gin.Engine) {

	// get global Monitor object
	m := ginmetrics.GetMonitor()

	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// set middleware for gin
	m.Use(router)

}

func SetupCORS(router *gin.Engine) {

	config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://google.com"}
	// config.AllowOrigins = []string{"http://google.com", "http://facebook.com"}
	config.AllowOrigins = []string{"*"}
	// config.AllowAllOrigins = true

	// router.Use(cors.New(config))
	router.Use(CORSMiddleware())

}
func SetupLogger(router *gin.Engine) {

	logger, _ := zap.NewProduction()
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	router.Use(ginzap.RecoveryWithZap(logger, true))

}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
