package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/time/rate"
)

var (
	metricsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of requests processed, labeled by path and method.",
	}, []string{"path", "method"})
	limiter = rate.NewLimiter(1, 5) // Allow 1 request per second with a burst of 5

	// Google OAuth2 configuration
	oauthConfig = &oauth2.Config{
		ClientID:     "1022579920436-rtt97h38dnm09k92kvl46r9qqgchjpg9.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:7171/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(metricsTotal)
}

func main() {
	router := gin.Default()

	// Middleware for Prometheus metrics
	router.Use(prometheusMiddleware)
	// Middleware for rate limiting
	router.Use(rateLimiterMiddleware)

	// Register routes
	router.GET("/users", getUsers)
	router.POST("/create", createUser)
	router.PUT("/update", updateUser)
	router.DELETE("/delete", deleteUser)
	router.GET("/health", healthCheck)

	// OAuth2 login and callback routes
	router.GET("/login", handleLogin)
	router.GET("/callback", handleCallback)

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Root route
	router.GET("/", handleRoot)

	// Graceful shutdown setup
	server := &http.Server{
		Addr:    ":7171",
		Handler: router,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Println("Server is running on http://localhost:7171")

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	wg.Wait()
	log.Println("Server stopped gracefully.")
}

// Handlers
func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get Users"})
}

func createUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User Created"})
}

func updateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User Updated"})
}

func deleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User Deleted"})
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func handleLogin(c *gin.Context) {
	// Redirect to OAuth2 login
	url := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func handleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
		return
	}
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func handleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Gin Server"})
}

// Middleware for Prometheus metrics
func prometheusMiddleware(c *gin.Context) {
	metricsTotal.WithLabelValues(c.FullPath(), c.Request.Method).Inc()
	c.Next()
}

// Middleware for rate limiting
func rateLimiterMiddleware(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
		c.Abort()
		return
	}
	c.Next()
}
