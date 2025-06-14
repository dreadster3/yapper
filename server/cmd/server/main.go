package main

import (
	"flag"
	"fmt"
	"maps"
	"slices"

	"github.com/dreadster3/yapper/server/internal/chat"
	"github.com/dreadster3/yapper/server/internal/platform/providers"
	"github.com/dreadster3/yapper/server/internal/platform/router"
	"github.com/dreadster3/yapper/server/internal/platform/router/middleware"
	"github.com/dreadster3/yapper/server/internal/user"
	"github.com/gin-gonic/gin/binding"
	en_locale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"go.uber.org/zap"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8000, "Port to listen on")
}

func _main() error {
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	userRepository := user.NewMongoRepository(logger.With(zap.String("repository", "user")))
	userHandler := user.NewUserHandler(userRepository)

	registeredProviders, err := providers.SetupProviders("http://localhost:11434")
	if err != nil {
		return err
	}

	chatHandler := chat.NewChatHandler(registeredProviders)

	jwtConfig := &middleware.JWTConfig{
		JWKSUrl: "http://localhost:9000/application/o/yapper/jwks/",
	}
	locale := en_locale.New()
	translator := ut.New(locale, locale)
	en_translator, ok := translator.GetTranslator("en")
	if !ok {
		return fmt.Errorf("translator for 'en' not found")
	}

	engine, err := router.SetupRouter(en_translator, jwtConfig, chatHandler, userHandler)
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
