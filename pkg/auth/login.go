package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/shjting0510/sa_user/pkg/user"
	"github.com/shjting0510/sa_user/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type LoginRes struct {
	Result struct {
		Token string `json:"token"`
	} `json:"result"`
	Success   bool   `json:"success"`
	ErrorInfo string `json:"errorInfo"`
}

func Login(c *gin.Context) {
	loginUri := user.BaseUrl + "/api/Authorize/Login"
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"userName":"admin", "password":"1q2w3E*"}`).
		Post(loginUri)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, utils.Response{Msg: "登录失败"})
		return
	}

	var res LoginRes
	if err := utils.Unmarshal[LoginRes](resp.Body(), &res); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, utils.Response{Msg: "登录失败"})
	}
	if resp.StatusCode() != http.StatusOK || !res.Success || len(res.Result.Token) == 0 {
		c.JSON(http.StatusInternalServerError, utils.Response{Msg: "登录失败"})
		return
	}

	user.GetUserInfoByToken(res.Result.Token)
	c.JSON(http.StatusOK, utils.Response{Msg: "登录成功"})
}
