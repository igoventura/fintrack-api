package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-core/internal/api/dto"
	"github.com/igoventura/fintrack-core/internal/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body dto.RegisterRequest true "Register User"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} handler.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorJSON(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := h.service.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		ErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Login godoc
// @Summary Login
// @Description Login with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body dto.LoginRequest true "Login User"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} handler.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorJSON(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		ErrorJSON(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	c.JSON(http.StatusOK, resp)
}
