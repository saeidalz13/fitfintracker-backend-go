package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	cn "github.com/saeidalz13/LifeStyle2/lifeStyleBack/config"
	database "github.com/saeidalz13/LifeStyle2/lifeStyleBack/db"

	h "github.com/saeidalz13/LifeStyle2/lifeStyleBack/handlers"
	"github.com/saeidalz13/LifeStyle2/lifeStyleBack/routes"
	"github.com/saeidalz13/LifeStyle2/lifeStyleBack/token"
)

func main() {
	mustPrepareReqVars()
	psqlDb := database.MustConnectToDb()
	tokenManager, err := token.NewPasetoMaker(cn.EnvVars.PasetoKey)
	if err != nil {
		log.Fatalln("Failed to extract Paseto Key!")
	}
	handlersConfig := h.NewHandlersConfig(
		&h.AuthHandlersManager{Db: psqlDb, TokenManager: tokenManager},
		&h.FinanceHandlersManager{Db: psqlDb},
		&h.FitnessHandlersManager{Db: psqlDb},
	)

	app := fiber.New()
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,X-CSRF-Token,Set-Cookie,Authorization",
		AllowOrigins:     cn.EnvVars.FrontEndUrl,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	routes.Setup(app, handlersConfig)

	log.Printf("Listening to port %v...", cn.EnvVars.Port)
	app.Listen(cn.EnvVars.Port)
}

func mustGetEnvVars() *cn.DotEnvVars {
	if os.Getenv("DEV_STAGE") != cn.DefaultDevStages.Production {
		if err := godotenv.Load(".env"); err != nil {
			panic("no .env file found")
		}
	}

	return &cn.DotEnvVars{
		FrontEndUrl: os.Getenv("FRONTENDURL"),
		Port:        os.Getenv("PORT"),
		PasetoKey:   os.Getenv("PASETO_KEY"),
		DbUrl:       os.Getenv("DATABASE_URL"),
		DevStage:    os.Getenv("DEV_STAGE"),
		GClientId:   os.Getenv("GOOGLE_CLIENT_ID"),
		GClientSec:  os.Getenv("GOOGLE_CLIENT_SEC"),
		GRedirUrl:   os.Getenv("GOOGLE_REDIRECT_URL"),
		GptApiKey:   os.Getenv("GPT_API_KEY"),
	}
}

func mustPrepareReqVars() {
	envVars := mustGetEnvVars()
	googleOAuthConfig := &oauth2.Config{
		RedirectURL:  envVars.GRedirUrl,
		ClientID:     envVars.GClientId,
		ClientSecret: envVars.GClientSec,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	cn.EnvVars = envVars
	cn.OAuthConfigFitFin = googleOAuthConfig
}
