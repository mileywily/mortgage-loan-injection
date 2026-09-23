package main

import (
	"log/slog"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"mortgage-loan-injection/internal/app/api/handler"
	"mortgage-loan-injection/internal/core/usecases"
	"mortgage-loan-injection/internal/infra/client"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// Setup JSON Logger (Replicando la estructura de Logback JSON Layout del legado)
	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Key = "timestamp"
			a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02T15:04:05.000Z"))
		}
		if a.Key == slog.MessageKey {
			a.Key = "message"
		}
		if a.Key == slog.LevelKey {
			// keep "level"
		}
		return a
	}
	
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: replaceAttr,
		Level:       slog.LevelDebug, // Legado usa Debug en logging.level (ver application.yml)
	}))
	slog.SetDefault(logger)

	// Start Datadog tracer
	tracer.Start(
		tracer.WithService("mortgage-loan-injection"),
		tracer.WithEnv(os.Getenv("ENV")),
	)
	defer tracer.Stop()

	// Fetch required environment variables
	port := os.Getenv("PORT")
	if port == "" {
		slog.Error("PORT environment variable is required")
		os.Exit(1)
	}

	// Initialize Gin
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("customRutPattern", func(fl validator.FieldLevel) bool {
			matched, _ := regexp.MatchString(`^[0-9]{1,2}\.[0-9]{3}\.[0-9]{3}-[0-9kK]$`, fl.Field().String())
			return matched
		})
		_ = v.RegisterValidation("customFechaPattern", func(fl validator.FieldLevel) bool {
			matched, _ := regexp.MatchString(`^[0-9]{2}-[0-9]{2}-[0-9]{4}$`, fl.Field().String())
			return matched
		})
	}

	// Setup middlewares (Datadog gintrace)
	router.Use(gintrace.Middleware("mortgage-loan-injection"))

	// Fetch FinnFlow configuration
	finnflowURL := os.Getenv("FINNFLOW_URL")
	if finnflowURL == "" {
		slog.Error("FINNFLOW_URL environment variable is required")
		os.Exit(1)
	}
	finnflowKey := os.Getenv("FINNFLOW_KEY")
	finnflowSecret := os.Getenv("FINNFLOW_SECRET")

	// Initialize dependencies
	finnflowClient := client.NewFinnFlowClient(finnflowURL, finnflowKey, finnflowSecret)
	uc := usecases.NewInjectionUseCaseImpl(finnflowClient)
	injectionHandler := handler.NewInjectionHandler(uc)

	v1 := router.Group("/v1/bfcl/mortgage-loan")
	{
		v1.POST("/injections", injectionHandler.InjectLoanApplication)
	}

	// Health check

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Start server
	slog.Info("Starting server", "port", port)
	if err := router.Run(":" + port); err != nil {
		slog.Error("Failed to start server", "error", err.Error())
		os.Exit(1)
	}
}
