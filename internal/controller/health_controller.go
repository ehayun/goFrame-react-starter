package controller

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) Check(r *ghttp.Request) {
	r.Response.WriteJson(g.Map{
		"status":  "ok",
		"message": "Server is running",
	})
}
