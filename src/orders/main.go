package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ShippingAddress struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zipCode"`
}

type OrderItem struct {
	ID        uint   `json:"-"         gorm:"primaryKey;autoIncrement"`
	OrderID   string `json:"-"         gorm:"column:order_id;index"`
	ProductID string `json:"productId" gorm:"column:product_id"`
	Name      string `json:"name"      gorm:"column:name"`
	Quantity  int    `json:"quantity"`
	UnitCost  int    `json:"price"     gorm:"column:unit_cost"`
	TotalCost int    `json:"totalCost" gorm:"column:total_cost"`
}

type Order struct {
	ID        string      `json:"id"          gorm:"primaryKey"`
	CreatedAt time.Time   `json:"createdDate"`
	FirstName string      `json:"firstName"   gorm:"column:first_name"`
	LastName  string      `json:"lastName"    gorm:"column:last_name"`
	Email     string      `json:"email"`
	Address1  string      `json:"address1"    gorm:"column:address1"`
	Address2  string      `json:"address2"    gorm:"column:address2"`
	City      string      `json:"city"`
	State     string      `json:"state"`
	ZipCode   string      `json:"zipCode"     gorm:"column:zip_code"`
	Items     []OrderItem `json:"items"       gorm:"foreignKey:OrderID"`
}

type CreateOrderRequest struct {
	ShippingAddress ShippingAddress `json:"shippingAddress"`
	Items           []OrderItem     `json:"items"`
}

var db *gorm.DB

func main() {

	// Logger en JSON: cada línea que se escriba va a stdout como un objeto JSON
	// y siempre lleva el campo "service".
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "orders")
	slog.SetDefault(logger)

	dsn := buildDSN()

	var err error
	for i := 0; i < 6; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newGormLogger()})
		if err == nil {
			break
		}
		slog.Warn("esperando a PostgreSQL", "intento", i+1, "error", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		fatal("no se pudo conectar a la base de datos", err)
	}

	if err := db.AutoMigrate(&Order{}, &OrderItem{}); err != nil {
		fatal("no se pudo migrar el esquema", err)
	}

	// gin.New() arranca sin middlewares (gin.Default() traía el logger de texto).
	// Primero nuestro logger JSON y después Recovery, para que si un handler
	// entra en pánico, Recovery responda 500 y el logger registre ese 500.
	r := gin.New()
	r.Use(requestLogger())
	r.Use(gin.Recovery())

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.GET("/orders", listOrders)
	r.POST("/orders", createOrder)

	port := getEnv("PORT", "8080")
	slog.Info("servidor escuchando", "port", port)
	if err := r.Run(":" + port); err != nil {
		fatal("el servidor no pudo escuchar en el puerto", err)
	}
}

func createOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := Order{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		FirstName: req.ShippingAddress.FirstName,
		LastName:  req.ShippingAddress.LastName,
		Email:     req.ShippingAddress.Email,
		Address1:  req.ShippingAddress.Address1,
		Address2:  req.ShippingAddress.Address2,
		City:      req.ShippingAddress.City,
		State:     req.ShippingAddress.State,
		ZipCode:   req.ShippingAddress.ZipCode,
	}

	for i := range req.Items {
		req.Items[i].OrderID = order.ID
		req.Items[i].TotalCost = req.Items[i].UnitCost * req.Items[i].Quantity
	}
	order.Items = req.Items

	if err := db.Create(&order).Error; err != nil {
		slog.Error("no se pudo guardar el pedido", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func listOrders(c *gin.Context) {
	var orders []Order
	db.Preload("Items").Order("created_at desc").Find(&orders)
	c.JSON(http.StatusOK, orders)
}

// fatal registra el error como JSON con nivel ERROR y termina el programa con código 1
// (lo mismo que hacía log.Fatal, pero en el formato nuevo).
func fatal(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}

func buildDSN() string {
	endpoint := getEnv("RETAIL_ORDERS_PERSISTENCE_ENDPOINT", "localhost:5432")
	host, port := endpoint, "5432"
	if idx := strings.LastIndex(endpoint, ":"); idx >= 0 {
		host = endpoint[:idx]
		port = endpoint[idx+1:]
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port,
		getEnv("RETAIL_ORDERS_PERSISTENCE_USERNAME", "retail_user"),
		os.Getenv("RETAIL_ORDERS_PERSISTENCE_PASSWORD"),
		getEnv("RETAIL_ORDERS_PERSISTENCE_NAME", "orders"),
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
