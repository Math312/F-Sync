package config

import (
	"fmt"
	"github.com/fsync/common"
	"github.com/fsync/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type BaiduYunConfigManageWeb struct {
}

func NewBaiduYunConfigManageWeb() BaiduYunConfigManageWeb {
	return BaiduYunConfigManageWeb{}
}

type AddNewAccessTokenRequest struct {
	Url string `json:"url"`
}

func (BaiduYunConfigManageWeb) AddNewAccessToken(c *gin.Context) {
	requestBody := AddNewAccessTokenRequest{}
	err := common.GetJsonBody(c, &requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Parameters Parse Error."})
		return
	}
	uniqueName := c.Param("name")
	parseResult, err := config.ParseBaiduOAuthUrl(requestBody.Url)
	if err != nil {
		c.JSON(200, gin.H{"message": "Input BaiduYun Access Url Incorrect"})
		return
	}

	config.WriteBaiduYunOAuthDataByUniqueName(uniqueName, parseResult.AccessToken, time.Now().Unix()+int64(parseResult.ExpiresIn))
	c.JSON(200, gin.H{"message": "Insert Successfully"})
}

func (BaiduYunConfigManageWeb) GetAccessTokenUrl(c *gin.Context) {
	url := fmt.Sprintf("http://openapi.baidu.com/oauth/2.0/authorize?response_type=token&client_id=%s&redirect_uri=oob&scope=basic,netdisk", config.APP_KEY)
	c.PureJSON(200, gin.H{"url": url})
}

func (BaiduYunConfigManageWeb) GetAllBaiduYunOAuthData(c *gin.Context) {
	result, _ := config.ListAllBaiduYunOAuthData()
	c.PureJSON(http.StatusOK, gin.H{"baiduYunOAuthData": result})
}

func (BaiduYunConfigManageWeb) DeleteBaiduYunOAuthDataByUniqueName(c *gin.Context) {
	uniqueName := c.Param("name")
	config.DeleteBaiduYunOAuthDataByUniqueName(uniqueName)
	c.PureJSON(http.StatusOK, gin.H{"message": "Delete Successfully"})
}
