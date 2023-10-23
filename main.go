package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/naderSameh/ticketing_support/api"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
	"github.com/naderSameh/ticketing_support/mail"
	"github.com/naderSameh/ticketing_support/util"
	worker "github.com/naderSameh/ticketing_support/woker"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

//	@title			Gin Swagger Example API
//	@version		1.0
//	@description	Ticketing support microservice
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Cypodsolutions
//	@contact.url	http://www.cypod.com/
//	@contact.email	naders@cypodsolutions.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used

// @host		localhost:8080
// @BasePath	/
// @schemes	http
func main() {

	err := util.Loadconfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	conn, err := sql.Open(viper.GetString("DB_DRIVER"), viper.GetString("DB_SOURCE"))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	store := db.NewStore(conn)
	taskDistributor := worker.NewRedisDistributor(viper.GetString("REDDIS_ADDR"))

	server, err := api.NewServer(store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	go runTaskProcessor(viper.GetString("REDDIS_ADDR"))

	server.Start(viper.GetString("SERVER_ADDRESS"))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server:")
	}

}

func runTaskProcessor(RedisAddress string) {
	mailer := mail.NewGmailSender(viper.GetString("GMAIL_NAME"), viper.GetString("GMAIL_EMAIL"), viper.GetString("GMAIL_PASS"))
	taskProcessor := worker.NewRedisTaskProcessor(RedisAddress, mailer)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
