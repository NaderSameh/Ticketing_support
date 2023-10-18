package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
	_ "github.com/naderSameh/ticketing_support/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) (*Server, error) {

	server := &Server{
		store: store,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://google.com"}
	// config.AllowOrigins = []string{"http://google.com", "http://facebook.com"}
	config.AllowOrigins = []string{"*"}
	config.AllowAllOrigins = true

	router.Use(cors.New(config))

	//tickets
	router.POST("/tickets", server.createTicket)              // Create new ticket
	router.DELETE("/tickets/:ticket_id", server.deleteTicket) // Delete ticket with its comments
	//comments
	router.DELETE("/tickets/:ticket_id/comments/:comment_id", server.deleteComment) // Delete a comment
	router.PUT("/tickets/:ticket_id/comments/:comment_id", server.updateComment)    //Edit comment text
	router.GET("/tickets/:ticket_id/comments", server.listComments)                 // Get comments for a ticket
	router.POST("/tickets/:ticket_id/comments", server.createComment)               //Add comment to a ticket
	//caterogries
	router.GET("/categories", server.listCategories) // Create Category
	//Swagger
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	authRoutes := router.Group("/").Use(authMiddleware())
	authRoutes.GET("/tickets", server.listTicket)              // Get list of tickets
	authRoutes.POST("/categories", server.createCategory)      // Create Category
	authRoutes.PUT("/tickets/:ticket_id", server.updateTicket) //Assign ticket or update its status

	server.router = router

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
