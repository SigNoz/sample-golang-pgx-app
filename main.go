package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
)

func initTracer() func(context.Context) error {

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

func main() {

	shutdown := initTracer()
	defer shutdown(context.Background())

	connectDB() // Initialize DB connection
	defer conn.Close()
	r := gin.Default()

	// Use otelgin middleware to instrument Gin
	r.Use(otelgin.Middleware("myService"))

	r.GET("/todos", func(c *gin.Context) {
		todos, err := getAllTodos(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving todos"})
			return
		}
		c.JSON(http.StatusOK, todos)
	})

	// Inside func main()

	// Create a todo
	r.POST("/todos", func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		createdTodo, err := createTodo(c.Request.Context(), todo.Title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating todo"})
			log.Println(err)
			return
		}
		c.JSON(http.StatusCreated, createdTodo)
	})

	// Update a todo
	r.PUT("/todos/:id", func(c *gin.Context) {
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		updatedTodo, err := updateTodo(c.Request.Context(), id, todo.Title, todo.Completed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating todo"})
			return
		}
		c.JSON(http.StatusOK, updatedTodo)
	})

	// Delete a todo
	r.DELETE("/todos/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := deleteTodo(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting todo"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
	})

	r.Run() // Listen and serve on 0.0.0.0:8080
}
