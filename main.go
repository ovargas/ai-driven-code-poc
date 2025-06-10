package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ovargas/ai-driven-code-poc/product"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	// Parse flags
	repoFlag := flag.String("repository", "memory", "Repository type: memory or db")
	dsnFlag := flag.String("dsn", "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local", "MySQL DSN (used if repository=db)")
	flag.Parse()

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	http.HandleFunc("/health", healthHandler)

	// Select repository implementation
	var productRepo product.Repository
	switch *repoFlag {
	case "memory":
		productRepo = product.NewInMemoryRepository()
	case "db":
		db, err := gorm.Open(mysql.Open(*dsnFlag), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		productRepo = product.NewGormRepository(db)
	default:
		log.Fatalf("repository flag value '%s' is not supported (use 'memory' or 'db')", *repoFlag)
	}

	productHandler := product.NewHandler(productRepo)
	productHandler.RegisterRoutes()

	// Channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("HTTP server error:", err)
		}
	}()

	<-stop // Wait for interrupt

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	}

	fmt.Println("Server exiting")
}
