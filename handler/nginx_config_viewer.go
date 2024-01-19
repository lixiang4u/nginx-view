package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NginxConfigViewer(ctx *gin.Context) {
	nHttp, err := parse(NginxConfig)
	if err != nil {
		ctx.HTML(http.StatusOK, "error.html", gin.H{
			"config_config_file": NginxConfig,
			"msg":                err.Error(),
		})
	}

	ctx.HTML(http.StatusOK, "nginx_config_viewer_tpl.html", gin.H{
		"config_config_file": NginxConfig,
		"title":              "Main website",
		"nHttp":              ToJsonString(nHttp, true),
	})
}
