package web

import (
	"github.com/fsync/web/config"
	"github.com/gin-gonic/gin"
)

type FsyncWebResponse interface {
}

func FsyncManageRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	baiduYunConfigManageWeb := config.NewBaiduYunConfigManageWeb()
	configManageApiV1 := r.Group("/api/manage/config")
	{
		configManageApiV1.POST("/baiduYun/name/:name", baiduYunConfigManageWeb.AddNewAccessToken)
		configManageApiV1.GET("/baiduYun/url", baiduYunConfigManageWeb.GetAccessTokenUrl)
		configManageApiV1.DELETE("/baiduYun/name/:name", baiduYunConfigManageWeb.DeleteBaiduYunOAuthDataByUniqueName)
		configManageApiV1.GET("/baiduYun/list", baiduYunConfigManageWeb.GetAllBaiduYunOAuthData)
	}
	return r
}
