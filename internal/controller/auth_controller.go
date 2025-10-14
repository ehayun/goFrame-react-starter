package controller

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"tzlev/internal/auth"
	"tzlev/internal/oauth"
	"tzlev/internal/repository"
	"tzlev/internal/session"
)

type AuthController struct {
	userRepo       *repository.UserRepository
	sessionManager *session.SessionManager
}

func NewAuthController() *AuthController {
	return &AuthController{
		userRepo:       repository.NewUserRepository(),
		sessionManager: session.NewSessionManager(),
	}
}

// Login handles Zehut/Password authentication
func (c *AuthController) Login(r *ghttp.Request) {
	ctx := r.Context()

	type LoginRequest struct {
		Zehut    string `json:"zehut" v:"required"`
		Password string `json:"password" v:"required"`
	}

	var req LoginRequest
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(g.Map{
			"error": "Invalid request",
		})
		return
	}

	// Find user by zehut
	user, err := c.userRepo.FindByZehut(ctx, req.Zehut)
	if err != nil {
		g.Log().Warning(ctx, "User not found:", req.Zehut)
		r.Response.WriteJson(g.Map{
			"error": "Invalid credentials",
		})
		return
	}

	// Verify password
	if !auth.CheckPassword(req.Password, user.HashedPassword) {
		g.Log().Warning(ctx, "Invalid password for user:", req.Zehut)
		r.Response.WriteJson(g.Map{
			"error": "Invalid credentials",
		})
		return
	}

	// Create session
	sessionID, _ := r.Session.Id()
	fullName := user.FirstName + " " + user.LastName
	sess := &session.Session{
		UserID:    0, // Legacy field
		Zehut:     user.Zehut,
		Email:     user.Email,
		Name:      fullName,
		CreatedAt: time.Now(),
	}

	if err := c.sessionManager.Create(ctx, sessionID, sess); err != nil {
		g.Log().Error(ctx, "Failed to create session:", err)
		r.Response.WriteJson(g.Map{
			"error": "Failed to create session",
		})
		return
	}

	// Set session cookie
	r.Session.Set("user_zehut", user.Zehut)
	r.Session.Set("user_email", user.Email)
	r.Session.Set("user_name", fullName)

	r.Response.WriteJson(g.Map{
		"status": "ok",
		"message": "Login successful",
	})
}

func (c *AuthController) GoogleLogin(r *ghttp.Request) {
	// Generate random state for CSRF protection
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Store state in session
	r.Session.Set("oauth_state", state)

	// Redirect to Google OAuth
	url := oauth.GoogleOAuthConfig.AuthCodeURL(state)
	r.Response.RedirectTo(url)
}

func (c *AuthController) GoogleCallback(r *ghttp.Request) {
	ctx := r.Context()

	// Verify state
	state := r.Get("state").String()
	sessionState, _ := r.Session.Get("oauth_state")
	sessionStateStr := sessionState.String()

	if state == "" || state != sessionStateStr {
		r.Response.WriteJson(g.Map{
			"error": "Invalid state parameter",
		})
		return
	}

	// Exchange code for token
	code := r.Get("code").String()
	token, err := oauth.GoogleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		g.Log().Error(ctx, "Failed to exchange token:", err)
		r.Response.WriteJson(g.Map{
			"error": "Failed to exchange token",
		})
		return
	}

	// Get user info from Google
	googleUser, err := oauth.GetGoogleUserInfo(ctx, token)
	if err != nil {
		g.Log().Error(ctx, "Failed to get user info:", err)
		r.Response.WriteJson(g.Map{
			"error": "Failed to get user info",
		})
		return
	}

	// Find user by email (match against existing users)
	user, err := c.userRepo.FindByEmail(ctx, googleUser.Email)
	if err != nil {
		// User doesn't exist - this means they're not authorized
		// In a real system, you might want to auto-create users or show a registration page
		g.Log().Warning(ctx, "User not found in system:", googleUser.Email)
		r.Response.WriteJson(g.Map{
			"error": "User not authorized. Please contact administrator.",
		})
		return
	}

	// User exists, update their confirmed_at if not set
	if user.ConfirmedAt == nil {
		now := time.Now()
		user.ConfirmedAt = &now
		user.Avatar = googleUser.Picture // Update avatar from Google
		if err := c.userRepo.Update(ctx, user); err != nil {
			g.Log().Error(ctx, "Failed to update user:", err)
		}
	}

	// Create session
	sessionID, _ := r.Session.Id()
	fullName := user.FirstName + " " + user.LastName
	sess := &session.Session{
		UserID: 0, // We'll use zehut instead
		Email:  user.Email,
		Name:   fullName,
	}

	if err := c.sessionManager.Create(ctx, sessionID, sess); err != nil {
		g.Log().Error(ctx, "Failed to create session:", err)
		r.Response.WriteJson(g.Map{
			"error": "Failed to create session",
		})
		return
	}

	// Set session cookie
	r.Session.Set("user_zehut", user.Zehut)
	r.Session.Set("user_email", user.Email)
	r.Session.Set("user_name", fullName)

	// Redirect to home page
	r.Response.RedirectTo("/")
}

func (c *AuthController) Logout(r *ghttp.Request) {
	ctx := r.Context()

	// Delete session from Redis
	sessionID, _ := r.Session.Id()
	if err := c.sessionManager.Delete(ctx, sessionID); err != nil {
		g.Log().Error(ctx, "Failed to delete session:", err)
	}

	// Clear session cookie
	r.Session.RemoveAll()

	r.Response.WriteJson(g.Map{
		"status":  "ok",
		"message": "Logged out successfully",
	})
}

func (c *AuthController) GetCurrentUser(r *ghttp.Request) {
	ctx := r.Context()

	// Get session from Redis
	sessionID, _ := r.Session.Id()
	sess, err := c.sessionManager.Get(ctx, sessionID)
	if err != nil {
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"error": "Not authenticated",
		})
		return
	}

	// Get full user data from database
	user, err := c.userRepo.FindByZehut(ctx, sess.Zehut)
	if err != nil {
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"error": "User not found",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"user": g.Map{
			"zehut":      user.Zehut,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"avatar":     user.Avatar,
			"role":       user.Role,
			"is_admin":   user.IsAdmin,
		},
	})
}
