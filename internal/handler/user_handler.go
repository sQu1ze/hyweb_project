package handler

import (
	"hyweb-api/internal/model"
	"hyweb-api/internal/response"
	"hyweb-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Register 處理使用者註冊
// @Summary      使用者註冊
// @Description  建立新帳號並將密碼加密儲存
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body     model.RegisterRequest  true  "註冊資料"
// @Success      200     {object}  response.Response      "成功"
// @Failure      400     {object}  response.Response      "請求參數錯誤"
// @Router       /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters", err.Error())
		return
	}

	if err := h.svc.Register(req); err != nil {
		if err.Error() == "email already exists" {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, "Registration successful", nil)
}

// Login 處理使用者登入
// @Summary      使用者登入
// @Description  驗證身分並回傳 JWT
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body     model.LoginRequest  true  "登入資料"
// @Success      200     {object}  response.Response    "成功"
// @Router       /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters", err.Error())
		return
	}

	token, err := h.svc.Login(req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			response.Error(c, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, "Login successful", gin.H{"token": token})
}

// ChangePassword 修改密碼
// @Summary      修改密碼
// @Description  驗證舊密碼後更換新密碼
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body      model.ChangePasswordRequest  true  "修改密碼資料"
// @Success      200     {object}  response.Response            "成功"
// @Router       /users/changePassword [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	email, _ := c.Get("email")
	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters", err.Error())
		return
	}

	if err := h.svc.ChangePassword(email.(string), req); err != nil {
		if err.Error() == "old password incorrect" || err.Error() == "new password cannot be the same as old password" {
			response.Error(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	response.Success(c, "Password changed successfully", nil)
}

// HealthCheck godoc
// @Summary 健康檢查
// @Description 檢查服務是否正常運作
// @Tags System
// @Produce json
// @Success 200 {object} response.Response
// @Router /health [get]
func (h *UserHandler) HealthCheck(c *gin.Context) {
	response.Success(c, "Service is healthy", nil)
}
