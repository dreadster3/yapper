package main

import (
	"context"
	"flag"
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"

	"github.com/dreadster3/yapper/server/internal/chats"
	"github.com/dreadster3/yapper/server/internal/messages"
	"github.com/dreadster3/yapper/server/internal/platform/database"
	"github.com/dreadster3/yapper/server/internal/platform/providers"
	"github.com/dreadster3/yapper/server/internal/platform/router"
	"github.com/dreadster3/yapper/server/internal/platform/router/middleware"
	"github.com/dreadster3/yapper/server/internal/profiles"
	"github.com/dreadster3/yapper/server/internal/steps"
	"github.com/gin-gonic/gin/binding"
	en_locale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	port   int
	jwkUrl string
	dbHost string
	dbPort int
	dbUser string
	dbPass string
)

func GetEnvDefault(key string, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}

func GetEnvIntDefault(key string, defaultValue int) int {
	if val, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(val); err != nil {
			return intVal
		}
	}

	return defaultValue
}

func initFlags() {
	flag.IntVar(&port, "port", GetEnvIntDefault("PORT", 8000), "Port to listen on")
	flag.StringVar(&jwkUrl, "jwk-url", os.Getenv("JWK_URL"), "The URL to the JWKS endpoint")
	flag.StringVar(&dbHost, "db-host", GetEnvDefault("DB_HOST", "mongo"), "The hostname of the database")
	flag.IntVar(&dbPort, "db-port", GetEnvIntDefault("DB_PORT", 27017), "The port of the database")
	flag.StringVar(&dbUser, "db-user", os.Getenv("DB_USER"), "The username of the database")
	flag.StringVar(&dbPass, "db-pass", os.Getenv("DB_PASS"), "The password of the database")
	flag.Parse()
}

func _main() error {
	godotenv.Load()
	initFlags()

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	if os.Getenv("GIN_MODE") != "release" {
		logger = zap.Must(zap.NewDevelopment())
	}
	defer logger.Sync()

	ctx := context.Background()

	db, closeDatabase, err := database.ConnectDatabase(ctx, dbHost, dbPort, dbUser, dbPass)
	if err != nil {
		return err
	}
	defer closeDatabase(ctx)

	profileRepository := profiles.NewProfileRepository(db, logger.With(zap.String("repository", "profile")))
	profileHandler := profiles.NewProfileHandler(profileRepository)

	registeredProviders, err := providers.SetupProviders("http://localhost:11434", logger)
	if err != nil {
		return err
	}

	chatRepository := chats.NewChatRepository(db, logger.With(zap.String("repository", "chat")))
	chatHandler := chats.NewChatHandler(registeredProviders, chatRepository)

	stepsRepository := steps.NewStepRepository(db, logger.With(zap.String("repository", "step")))
	messageRepository := messages.NewMessageRepository(db, logger.With(zap.String("repository", "message")))
	messageHandler := messages.NewMessageHandler(messageRepository, chatRepository, stepsRepository, registeredProviders)

	jwtConfig := &middleware.JWTConfig{
		JWKSUrl: jwkUrl,
	}
	locale := en_locale.New()
	translator := ut.New(locale, locale)
	en_translator, ok := translator.GetTranslator("en")
	if !ok {
		return fmt.Errorf("translator for 'en' not found")
	}

	engine, err := router.SetupRouter(en_translator, jwtConfig, profileRepository, chatHandler, profileHandler, messageHandler)
	if err != nil {
		return err
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("registered_provider", providers.ValidateRegisteredProvider(registeredProviders))

		if err := en_translations.RegisterDefaultTranslations(v, en_translator); err != nil {
			return err
		}

		validProviders := slices.Collect(maps.Keys(registeredProviders))
		v.RegisterTranslation("registered_provider", en_translator, func(ut ut.Translator) error {
			return ut.Add("registered_provider", fmt.Sprintf("{0} is not a registered provider. Valid options are: %s", validProviders), true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("registered_provider", fe.Field())
			return t
		})
	}

	engine.Run(fmt.Sprintf(":%d", port))

	return nil
}

func main() {
	if err := _main(); err != nil {
		panic(err)
	}
}
