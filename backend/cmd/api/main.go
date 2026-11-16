package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/db"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/handlers"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/repository"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/service"
)

func main() {
	ctx := context.Background()
	// Load backend/.env locally. In production (Lambda), the file does not exist
	// and variables come from the AWS configuration, so the error is ignored.
	_ = godotenv.Load()

	// Development mode: SIMULATED authentication and transfers. Local only.
	devAuth := os.Getenv("DEV_AUTH") == "true"
	if devAuth {
		if os.Getenv("APP_ENV") == "production" {
			log.Fatal("DEV_AUTH must never be enabled in production")
		}
		log.Println("WARNING: DEV_AUTH mode active (simulated authentication and transfers). Local only.")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "eu-west-1"
	}

	database, err := db.InitDynamoDB(ctx, region)
	if err != nil {
		log.Fatalf("DynamoDB connection failed: %v", err)
	}
	client := database.Client

	// Repositories (one per table)
	campaignRepo := repository.NewCampaignRepository(client, "campaigns")
	donationRepo := repository.NewDonationRepository(client, "donations")
	donorRepo := repository.NewDonorRepository(client, "donors")
	mediaRepo := repository.NewMediaRepository(client, "media")
	organiserRepo := repository.NewOrganizerRepository(client, "organizers")
	payoutRepo := repository.NewPayoutAccountRepository(client, "payout_accounts")

	// Transfer provider: simulated in dev, otherwise refuses everything until an aggregator is plugged in.
	var payoutProvider service.PayoutProvider = service.NoPayoutProvider{}
	if devAuth {
		payoutProvider = service.DevPayoutProvider{}
	}

	// Services
	campaignService := service.NewCampaignService(campaignRepo)
	mediaService := service.NewMediaService(mediaRepo)
	donorService := service.NewDonorService(donorRepo)
	organiserService := service.NewOrganiserService(organiserRepo)
	payoutService := service.NewPayoutAccountService(payoutRepo)
	donationService := service.NewDonationService(donationRepo, campaignRepo, payoutRepo, organiserRepo, payoutProvider)

	// Handlers
	campaignHandler := handlers.NewCampaignHandler(campaignService)
	viewSvc := service.NewCampaignViewService(campaignService, mediaRepo)
	campaignHandler.SetViewService(viewSvc)
	shareHandler := handlers.NewShareHandler(viewSvc)
	mediaHandler := handlers.NewMediaHandler(mediaService, campaignService)
	organizerHandler := handlers.NewOrganizerHandler(organiserService)
	payoutHandler := handlers.NewPayoutAccountHandler(payoutService)
	donationHandler := handlers.NewDonationHandler(donationService, donorService)
	historyHandler := handlers.NewDonationHistoryHandler(campaignService, donationService)
	donorHandler := handlers.NewDonorHandler(donorService, donationService)
	commentHandler := handlers.NewCommentHandler(campaignService, donationService)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization", "X-Dev-Organizer-ID", "X-Dev-Donor-ID"},
	}))

	// Maximum request body size: 1 MB.
	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	})

	// Public routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Demso API!"})
	})
	router.GET("/campaigns/:slug", campaignHandler.GetBySlug)
	router.GET("/campaigns", campaignHandler.ListPublic)
	router.GET("/share/:slug", shareHandler.Page)
	router.GET("/campaigns/:slug/media", mediaHandler.ListPublic)
	router.GET("/campaigns/:slug/donations", historyHandler.ListPublic)
	router.GET("/organizers/:id", organizerHandler.GetPublic)
	router.POST("/donations", donationHandler.Create)
	router.POST("/organizers", organizerHandler.Register) // closed outside dev mode

	// Connected organizer routes (401 until authentication is implemented)
	me := router.Group("/me")
	{
		me.POST("/campaigns", campaignHandler.Create)
		me.GET("/campaigns", campaignHandler.ListMine)
		me.POST("/campaigns/:id/media", mediaHandler.Attach)
		me.PUT("/campaigns/:id/media/:mediaID/cover", mediaHandler.SetCover)
		me.DELETE("/campaigns/:id/media/:mediaID", mediaHandler.Delete)
		me.GET("/payout-accounts", payoutHandler.List)
		me.POST("/payout-accounts", payoutHandler.Add)
		me.POST("/payout-accounts/:id/activate", payoutHandler.SetActive)
		me.GET("/donor", donorHandler.GetProfile)
		me.GET("/donations", donorHandler.ListDonations)
		me.PUT("/donor/photo", donorHandler.SetPhoto)
		me.DELETE("/donor/photo", donorHandler.RemovePhoto)
		me.GET("/campaigns/:id/donations", commentHandler.ListForOrganizer)
		me.PUT("/campaigns/:id/donations/:donationID/comment-visibility", commentHandler.SetVisibility)
		me.PUT("/donations/:id/comment", commentHandler.UpdateOwn)
		me.DELETE("/donations/:id/comment", commentHandler.DeleteOwn)
		me.PUT("/campaigns/:id/donations/:donationID/reply", commentHandler.Reply)
		me.DELETE("/campaigns/:id/donations/:donationID/reply", commentHandler.DeleteReply)
	}

	uploadSvc, err := service.NewUploadService(ctx)
	if err != nil {
		log.Println("WARNING: uploads disabled:", err)
	} else {
		mediaHandler.SetUploadService(uploadSvc)
		uploadHandler := handlers.NewUploadHandler(uploadSvc)
		me.POST("/uploads/presign", uploadHandler.Presign)
	}

	// Test route: simulates payment confirmation. Exists only in dev mode.
	if devAuth {
		devHandler := handlers.NewDevHandler(donationService)
		router.POST("/dev/donations/:id/complete", devHandler.CompleteDonation)
	}

	// The webhook route will be added once the aggregator is chosen:
	// it will need a signature verifier and an event translator for it.

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
