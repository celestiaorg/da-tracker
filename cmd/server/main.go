package main

import (
	"context"
	"log"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/celestiaorg/validator-da-tracker/config"
	"github.com/celestiaorg/validator-da-tracker/pkg/database"
	"github.com/celestiaorg/validator-da-tracker/pkg/handlers"
	"github.com/celestiaorg/validator-da-tracker/pkg/metrics"
	"github.com/celestiaorg/validator-da-tracker/pkg/repository"
)

const (
	scrapeInterval  = time.Second * 30
	processInterval = time.Second * 30
)

func main() {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// we don't define the file here as it is already defined in the OS env in the CP side
	// check serverless.yml. You are free to define it here if you want to run it locally
	cfg := config.LoadEnv()

	log.Println("Initalizing DB...")
	database.InitDB(cfg)   // Initialize the database connection
	db := database.GetDB() // Get the database instance

	log.Println("Initalizing Validator Repository & Handler...")
	validatorRepoForHandlers := repository.NewValidatorRepository(db)
	validatorHandler := handlers.NewValidatorHandler(validatorRepoForHandlers)

	prometheusClient := metrics.NewPrometheusClient()
	// Start the metrics scraper with the context
	log.Println("PROMETHEUS_URL: ", cfg.PromURL)
	log.Println("Starting metrics scraper...")
	g.Go(func() error {
		return metrics.StartMetricsScraper(ctx, db, prometheusClient, cfg.PromURL, cfg.PromToken, scrapeInterval)
	})

	log.Println("Starting metrics processor...")
	g.Go(func() error {
		return metrics.StartMetricsProcessor(ctx, db, prometheusClient, cfg.PromURL, cfg.PromToken, processInterval)
	})

	g.Go(func() error {
		router := gin.Default()

		// Setup your routes
		router.GET("/validators", validatorHandler.GetAllValidators())
		router.GET("/validators/:id", validatorHandler.GetValidatorByID())
		router.GET("/validators/email", validatorHandler.GetValidatorByEmail())
		router.GET("/validators/name", validatorHandler.GetValidatorByName())
		// e.g. /validators/peerids?peerid=id1&peerid=id2
		router.GET("/validators/peerids", validatorHandler.GetValidatorsByPeerIDs())
		// listen and serve on 0.0.0.0:8080
		router.GET("/metrics", gin.WrapH(promhttp.Handler()))
		err := router.Run(":8080")
		if err != nil {
			log.Fatalf("Error running router: %v", err)
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("Encountered error: %v", err)
	}
}
