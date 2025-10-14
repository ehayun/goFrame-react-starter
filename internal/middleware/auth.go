package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"tzlev/internal/session"
)

func Auth() func(r *ghttp.Request) {
	sessionManager := session.NewSessionManager()

	return func(r *ghttp.Request) {
		ctx := r.Context()

		// Get session ID from cookie
		sessionID, err := r.Session.Id()
		if err != nil || sessionID == "" {
			r.Response.Status = 401
			r.Response.WriteJson(g.Map{
				"error": "Not authenticated",
			})
			return
		}

		// Verify session in Redis
		sess, err := sessionManager.Get(ctx, sessionID)
		if err != nil {
			g.Log().Warning(ctx, "Invalid session:", err)
			r.Response.Status = 401
			r.Response.WriteJson(g.Map{
				"error": "Invalid or expired session",
			})
			return
		}

		// Refresh session TTL
		if err := sessionManager.Refresh(ctx, sessionID); err != nil {
			g.Log().Warning(ctx, "Failed to refresh session:", err)
		}

		// Store user info in request context
		r.SetCtxVar("user_id", sess.UserID)
		r.SetCtxVar("user_email", sess.Email)
		r.SetCtxVar("user_name", sess.Name)

		r.Middleware.Next()
	}
}
