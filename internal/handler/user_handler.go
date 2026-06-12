package handler

import (
	"net/http"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
	jwtutil "github.com/miank1/ecommerce_backend/pkg/jwt"
)

type UserHandler struct {
	Svc *service.UserService
}

type apiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{Svc: s}
}

func writeError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, apiResponse{
		Status:  "error",
		Message: message,
	})
}

func writeSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, apiResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

/* Register */
type registerReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	u, err := h.Svc.Register(req.Name, req.Email, req.Password)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccess(c, http.StatusCreated, "user registered successfully", gin.H{"user": u})
}

/* Login */
type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Step 1: Validate user credentials
	user, err := h.Svc.Authenticate(req.Email, req.Password)
	if err != nil {
		writeError(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// Step 2: Generate JWT Token
	token, err := jwtutil.GenerateToken(user.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// Step 3: Respond with token and user info
	writeSuccess(c, http.StatusOK, "login successful", gin.H{
		"token": token,
		"user":  user,
	})
}

/* Me */
func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		writeError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := h.Svc.GetByID(userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	if user == nil {
		writeError(c, http.StatusNotFound, "user not found")
		return
	}
	writeSuccess(c, http.StatusOK, "user fetched successfully", gin.H{"user": user})
}
