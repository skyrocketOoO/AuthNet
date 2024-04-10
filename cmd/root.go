/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"time"

	"github.com/skyrocketOoO/AuthNet/api"
	"github.com/skyrocketOoO/AuthNet/config"
	"github.com/skyrocketOoO/AuthNet/docs"
	"github.com/skyrocketOoO/AuthNet/domain"
	"github.com/skyrocketOoO/AuthNet/internal/delivery/rest"
	"github.com/skyrocketOoO/AuthNet/internal/delivery/rest/middleware"
	"github.com/skyrocketOoO/AuthNet/internal/infra/graph"
	"github.com/skyrocketOoO/AuthNet/internal/infra/repository/mongo"
	"github.com/skyrocketOoO/AuthNet/internal/infra/repository/redis"
	"github.com/skyrocketOoO/AuthNet/internal/infra/repository/sql"
	"github.com/skyrocketOoO/AuthNet/internal/usecase"

	errors "github.com/rotisserie/eris"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func workFunc(cmd *cobra.Command, args []string) {
	zerolog.TimeFieldFormat = time.RFC3339
	// human-friendly logging without efficiency
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Logger initialized")

	if err := config.ReadConfig(); err != nil {
		log.Fatal().Msg(errors.ToString(err, true))
	}

	docs.SwaggerInfo.Title = "Swagger API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http"}

	var dbRepo domain.DbRepository
	mode, err := cmd.Flags().GetInt("mode")
	if err != nil {
		log.Fatal().Msg(errors.ToString(err, true))
	}
	switch mode {
	case 1:
		sqlDb, disconnectDb, err := sql.InitDB("pg")
		if err != nil {
			log.Fatal().Msg(errors.ToString(err, true))
		}
		defer disconnectDb()
		dbRepo, err = sql.NewSqlRepository(sqlDb)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	case 2:
		mongoClient, disconnectDb, err := mongo.InitDb()
		if err != nil {
			log.Fatal().Msg(errors.ToString(err, true))
		}
		defer disconnectDb()
		dbRepo, err = mongo.NewMongoRepository(mongoClient)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	case 3:
		rdsCli, disconnectDb, err := redis.InitDb()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		defer disconnectDb()
		dbRepo, err = redis.NewRedisRepository(rdsCli)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	case 4:
		rdsCli, disconnectDb, err := redis.InitDb()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		defer disconnectDb()
		dbRepo, err = redis.NewRedis2Repository(rdsCli)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	case 5:
		sqlDb, disconnectDb, err := sql.InitDB("cockroachdb")
		if err != nil {
			log.Fatal().Msg(errors.ToString(err, true))
		}
		defer disconnectDb()
		dbRepo, err = sql.NewSqlRepository(sqlDb)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
	case 6:

	default:
		log.Fatal().Msg("mode not supported")
	}

	var graphInfra domain.GraphInfra
	switch mode {
	case 6:
		rdsCli, disconnectDb, err := redis.InitDb()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		defer disconnectDb()
		dbRepo, err = redis.NewRedis2Repository(rdsCli)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		graphInfra = graph.NewRedisLuaGraphInfra(dbRepo, rdsCli)
	default:
		graphInfra = graph.NewGraphInfra(dbRepo)
	}

	usecase := usecase.NewUsecase(dbRepo, graphInfra)
	delivery := rest.NewDelivery(usecase)

	router := gin.Default()
	router.Use(middleware.CORS())
	api.Binding(router, delivery)

	port, _ := cmd.Flags().GetString("port")
	router.Run(":" + port)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "A brief description of your application",
	Long:  `The longer description`,
	Run:   workFunc,
}

// Adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.Flags().StringP("port", "p", "8081", "port")
	rootCmd.Flags().IntP("mode", "m", 1, "mode")
}
