package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/internal/controller/middleware"
	"shop/internal/usecase"
	"time"
)

type Handler struct {
	service *usecase.UsecaseImplementation
}

func NewHandler(service *usecase.UsecaseImplementation) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle() http.Handler {
	router := gin.Default()

	router.POST("/api/auth", h.AuthHandler)
	router.GET("/api/info", middleware.AuthMiddleware(), h.InfoHandler)
	router.POST("/api/sendCoin", middleware.AuthMiddleware(), h.SendCoinHandler)
	router.POST("/api/buy/:item", middleware.AuthMiddleware(), h.BuyItemHandler)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented,
			gin.H{"code": http.StatusNotImplemented, "error": "not implemented"})
	})
	return router
}

func (h *Handler) AuthHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.service.Auth(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := middleware.JWT{}.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	expirationTime := time.Now().Add(5 * time.Hour)
	c.SetCookie("accessToken", token, int(expirationTime.Unix()), "/", "localhost", false, true)

	c.JSON(http.StatusOK,
		gin.H{"response": gin.H{"accessToken": token}})

}

func (h *Handler) SendCoinHandler(c *gin.Context) {
	var req struct {
		ReceiverUsername string  `json:"receiver_username"`
		Amount           float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	senderUsername := c.MustGet("username").(string)

	err, transaction := h.service.CreateTransaction(senderUsername, req.ReceiverUsername, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
func (h *Handler) BuyItemHandler(c *gin.Context) {
	itemName := c.Param("item")
	username := c.MustGet("username").(string)

	err, purchase := h.service.CreatePurchase(username, itemName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, purchase)
}
func (h *Handler) InfoHandler(c *gin.Context) {
	username := c.MustGet("username").(string)

	purchases, err := h.service.GetPurchasesForUserByUserGUID(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transactions, err := h.service.GetTransactionsForUserByUserGUID(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"purchases":    purchases,
		"transactions": transactions,
	})
}
