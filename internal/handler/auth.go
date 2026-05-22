package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohammed-ayoub-js/endline-backend/internal/config"
	"github.com/mohammed-ayoub-js/endline-backend/internal/models"
	"github.com/mohammed-ayoub-js/endline-backend/internal/utils"
	"gorm.io/gorm"
)

type RegisterRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type AuthResponse struct {
    AccessToken string `json:"access_token"`
    User        models.User `json:"user"`
}

const sec = false;

func Register(db *gorm.DB, cfg *config.Config) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var req RegisterRequest
        if err := c.BodyParser(&req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
        }
        if req.Username == "" || req.Password == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username and password required"})
        }
        if len(req.Password) < 8 {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password must be at least 8 characters"})
        }

        var existing models.User
        if err := db.Where("username = ?", req.Username).First(&existing).Error; err == nil {
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already taken"})
        }

        hashed, err := utils.HashPassword(req.Password)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal error"})
        }

        user := models.User{
            Username: req.Username,
            Password: hashed,
            Points:   0,
            Level:    "beginner",
        }
        if err := db.Create(&user).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
        }

        accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, cfg)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
        }
        refreshToken, err := utils.GenerateRefreshToken(user.ID, cfg)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate refresh token"})
        }

        refreshModel := models.RefreshToken{
            UserID:    user.ID,
            Token:     refreshToken,
            ExpiresAt: time.Now().Add(cfg.RefreshTTL),
            Revoked:   false,
        }
        db.Create(&refreshModel)

        c.Cookie(&fiber.Cookie{
            Name:     "refresh_token",
            Value:    refreshToken,
            Expires:  time.Now().Add(cfg.RefreshTTL),
            HTTPOnly: true,
            Secure:   sec, 
            SameSite: "Strict",
            Path:     "/api/auth/refresh",
        })

        return c.Status(fiber.StatusCreated).JSON(AuthResponse{
            AccessToken: accessToken,
            User:        user,
        })
    }
}

func Login(db *gorm.DB, cfg *config.Config) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var req LoginRequest
        if err := c.BodyParser(&req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
        }

        var user models.User
        if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
        }

        ok, err := utils.VerifyPassword(user.Password, req.Password)
        if err != nil || !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
        }

        accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, cfg)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
        }
        refreshToken, err := utils.GenerateRefreshToken(user.ID, cfg)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate refresh token"})
        }

        refreshModel := models.RefreshToken{
            UserID:    user.ID,
            Token:     refreshToken,
            ExpiresAt: time.Now().Add(cfg.RefreshTTL),
            Revoked:   false,
        }
        db.Create(&refreshModel)

        c.Cookie(&fiber.Cookie{
            Name:     "refresh_token",
            Value:    refreshToken,
            Expires:  time.Now().Add(cfg.RefreshTTL),
            HTTPOnly: true,
            Secure:   sec,
            SameSite: "Strict",
            Path:     "/api/auth/refresh",
        })

        return c.JSON(AuthResponse{
            AccessToken: accessToken,
            User:        user,
        })
    }
}

func Refresh(db *gorm.DB, cfg *config.Config) fiber.Handler {
    return func(c *fiber.Ctx) error {
        refreshToken := c.Cookies("refresh_token")
        if refreshToken == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing refresh token"})
        }

        token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(cfg.RefreshSecret), nil
        })
        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
        }

        _, ok := token.Claims.(*jwt.RegisteredClaims)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
        }

        var storedToken models.RefreshToken
        if err := db.Where("token = ? AND revoked = ?", refreshToken, false).First(&storedToken).Error; err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token revoked or not found"})
        }

        if storedToken.ExpiresAt.Before(time.Now()) {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Refresh token expired"})
        }

        var user models.User
        if err := db.First(&user, storedToken.UserID).Error; err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
        }

        db.Model(&storedToken).Update("revoked", true)

        newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Username, cfg)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
        }
        newRefreshToken, err := utils.GenerateRefreshToken(user.ID, cfg)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate refresh token"})
        }

        newRefreshModel := models.RefreshToken{
            UserID:    user.ID,
            Token:     newRefreshToken,
            ExpiresAt: time.Now().Add(cfg.RefreshTTL),
            Revoked:   false,
        }
        db.Create(&newRefreshModel)

        c.Cookie(&fiber.Cookie{
            Name:     "refresh_token",
            Value:    newRefreshToken,
            Expires:  time.Now().Add(cfg.RefreshTTL),
            HTTPOnly: true,
            Secure:   sec,
            SameSite: "Strict",
            Path:     "/api/auth/refresh",
        })

        return c.JSON(fiber.Map{"access_token": newAccessToken})
    }
}

func Logout(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        refreshToken := c.Cookies("refresh_token")
        if refreshToken != "" {
            db.Model(&models.RefreshToken{}).Where("token = ?", refreshToken).Update("revoked", true)
        }
        c.Cookie(&fiber.Cookie{
            Name:     "refresh_token",
            Value:    "",
            Expires:  time.Now().Add(-time.Hour),
            HTTPOnly: true,
            Secure:   true,
            SameSite: "Strict",
            Path:     "/api/auth/refresh",
        })
        return c.JSON(fiber.Map{"message": "Logged out successfully"})
    }
}

func Profile(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)
    return c.JSON(fiber.Map{
        "user_id":  userID,
        "username": c.Locals("username"),
        "message":  "This is a protected route",
    })
}