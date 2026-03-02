package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	BankServiceURl  string
	OrderServiceURl string
	UserServiceURl  string
	WebServiceURl   string
}

type ServiceProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

type APIGateway struct {
	config   *Config
	router   *mux.Router
	services map[string]*ServiceProxy
}

func (g *APIGateway) initService() {
	services := map[string]string{
		"bank-service":  g.config.BankServiceURl,
		"order-service": g.config.OrderServiceURl,
		"user-service":  g.config.UserServiceURl,
		"web-service":   g.config.WebServiceURl,
	}
	g.services = make(map[string]*ServiceProxy)

	for name, serviceURL := range services {
		target, err := url.Parse(serviceURL)
		if err != nil {
			log.Fatalf("failed to parse %s service url: %v", name, err)
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		g.services[name] = &ServiceProxy{
			target: target,
			proxy:  proxy,
		}
	}
}

func (g *APIGateway) proxyToService(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("proxing request to %s service: %s", serviceName, r.URL.Path)
		servis, exist := g.services[serviceName]
		if !exist {
			http.Error(w, "Service unavaible", http.StatusInternalServerError)
			return
		}

		servis.proxy.ServeHTTP(w, r)
	}
}

func (g *APIGateway) proxyToBankService(w http.ResponseWriter, r *http.Request) {
	g.proxyToService("bank-service")(w, r)
}

func (g *APIGateway) proxyToOrderService(w http.ResponseWriter, r *http.Request) {
	g.proxyToService("order-service")(w, r)
}

func (g *APIGateway) proxyToUserService(w http.ResponseWriter, r *http.Request) {
	g.proxyToService("user-service")(w, r)
}

func (g *APIGateway) proxyToWebService(w http.ResponseWriter, r *http.Request) {
	g.proxyToService("web-service")(w, r)
}

func (g *APIGateway) setRoutes() {
	g.router.HandleFunc("/", g.proxyToWebService).Methods("GET")
	g.router.HandleFunc("/about-us", g.proxyToWebService).Methods("GET")
	g.router.HandleFunc("/login", g.proxyToUserService).Methods("GET")
	g.router.HandleFunc("/register", g.proxyToUserService).Methods("GET")
	g.router.HandleFunc("/login", g.proxyToUserService).Methods("POST")
	g.router.HandleFunc("/register", g.proxyToUserService).Methods("POST")
	g.router.HandleFunc("/buy-merch", g.proxyToWebService).Methods("GET")
	g.router.HandleFunc("/donate", g.proxyToWebService).Methods("GET")
	g.router.HandleFunc("/donate/{category}/{username}/{moneySum}", g.proxyToWebService).Methods("POST")

	bank := g.router.PathPrefix("/").Subrouter()
	bank.HandleFunc("/my-wallet", g.proxyToBankService).Methods("GET")

	orders := g.router.PathPrefix("/").Subrouter()

	orders.HandleFunc("/add-to-cart/{product-id:[0-9]+}/{cart-id:[0-9]+}/{quantity:[0-9]+}", g.proxyToOrderService).Methods("POST")
	orders.HandleFunc("/get-all-from-cart", g.proxyToOrderService).Methods("GET")

}

func NewAPIGateway(cfg *Config) *APIGateway {
	router := mux.NewRouter()
	gateway := &APIGateway{
		config: cfg,
		router: router,
	}
	gateway.initService()
	gateway.setRoutes()
	return gateway
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment")
	}

	cfg := &Config{
		Port:            os.Getenv("PORT"),
		BankServiceURl:  os.Getenv("BANK_SERVICE_URL"),
		OrderServiceURl: os.Getenv("ORDER_SERVICE_URL"),
		UserServiceURl:  os.Getenv("USER_SERVICE_URL"),
		WebServiceURl:   os.Getenv("WEB_SERVICE_URL"),
	}

	gateway := NewAPIGateway(cfg)

	log.Printf("API Gateway starting on port %s <-", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, gateway.router); err != nil {
		log.Fatalf("Failed to start API Gateway: %s", err)
	}
}
