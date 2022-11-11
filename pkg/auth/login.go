package auth

import (
	"fmt"
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

type LoginForm struct {
	UserName string `json:"userName" form:"userName" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{Msg: "登录失败", Detail: err.Error()})
		return
	}
	loginUri := user.BaseUrl + "/api/Authorize/Login"
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"userName":"%s", "password":"%s"}`, form.UserName, form.Password)).
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

	c.JSON(http.StatusOK, utils.Response{Msg: "登录成功", Data: map[string]string{"token": res.Result.Token}})
}
