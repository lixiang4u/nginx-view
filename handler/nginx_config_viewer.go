package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"time"
)

func NginxConfigViewer(ctx *gin.Context) {
	ctx.File(filepath.Join(AppRoot(), "tpl/nginx_config_viewer_tpl.html"))
}
func NginxConfigJson(ctx *gin.Context) {
	nHttp, err := parse(NginxConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"time":        time.Now().Unix(),
		"config_file": NginxConfig,
		"config_data": nHttp,
	})
}
