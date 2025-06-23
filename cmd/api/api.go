package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/olliekm/gorestapi/service/cart"
	"github.com/olliekm/gorestapi/service/order"
	"github.com/olliekm/gorestapi/service/product"
	"github.com/olliekm/gorestapi/service/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// Register user service routes
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// Register product service routes
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Adjust as needed
		Password: "",
		DB:       0, // Use default DB
	})
	if err := redisClient.Ping(redisClient.Context()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	productStore := product.NewStore(s.db, redisClient)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	// Cart service routes
	orderStore := order.NewStore(s.db)

	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
