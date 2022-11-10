package user

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/shjting0510/sa_user/inits"
	"github.com/shjting0510/sa_user/pkg/auth/jwt"
	"github.com/shjting0510/sa_user/utils"

	log "github.com/sirupsen/logrus"
	"net/http"
)

const BaseUrl = "http://192.168.1.253:44313"

type UsersModel struct {
	UserID   int64  `form:"user_id" json:"user_id"`
	UserName string `form:"user_name" json:"user_name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterForm struct {
	UsersModel
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required"`
}

func Register(c *gin.Context) {
	form := RegisterForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg:    "注册失败",
			Detail: err.Error(),
			Data:   nil,
		})
		inits.Log.WithFields(log.Fields{"action": "user.Register"}).Debug(err)
		return
	}

	if form.Password != form.ConfirmPassword {
		c.JSON(http.StatusBadRequest, utils.Response{Msg: "两次密码不一致"})
		return
	}
	user, err := GetUserByName(form.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Msg: "注册失败"})
		inits.Log.WithFields(log.Fields{"action": "user.Register"}).Error(err)
		return
	}
	if user.UserID > 0 {
		c.JSON(http.StatusBadRequest, utils.Response{Msg: "该用户已存在"})
		return
	}

	passwd, err := utils.GenPassword(form.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg: "注册失败",
		})
	}

	storyUser := `insert into users(user_name, password) values($1, $2)`
	ct, err := inits.GetDB().Exec(context.Background(), storyUser, form.UserName, passwd)
	if err != nil || ct.RowsAffected() < 1 {
		inits.Log.WithFields(log.Fields{"action": "user.Register"}).Error(err)
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg: "注册失败 ",
		})
	}
	c.JSON(http.StatusOK, utils.Response{
		Msg: "注册成功",
	})
}

type LoginForm struct {
	UsersModel
}

func Login(c *gin.Context) {
	form := LoginForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{Msg: "登录失败", Detail: err.Error()})
		inits.Log.WithFields(log.Fields{"action": "user.Login"}).Debug(err)
		return
	}
	user, err := GetUserByName(form.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Msg: "登录失败"})
		inits.Log.WithFields(log.Fields{"action": "user.Login"}).Error(err)
		return
	}
	if user.UserID < 1 {
		c.JSON(http.StatusBadRequest, utils.Response{Msg: "用户不存在"})
		return
	}
	if !utils.CheckPassword(form.Password, user.Password) {
		c.JSON(http.StatusBadRequest, utils.Response{Msg: "密码不正确"})
		return
	}

	token, err := jwt.GenToken(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Msg: "登录失败"})
		inits.Log.WithFields(log.Fields{"action": "user.Login"}).Error(err)
		return
	}

	c.JSON(http.StatusOK, utils.Response{Msg: "登录成功", Data: token})
}

// GetUserByName gets user's detail by the given username.
func GetUserByName(userName string) (UsersModel, error) {
	var user UsersModel
	queryUser := `select user_id, user_name, password from users where user_name = $1 limit 1`
	err := pgxscan.Get(context.Background(), inits.GetDB(), &user, queryUser, userName)
	if err != nil && !pgxscan.NotFound(err) {
		return UsersModel{}, err
	}

	return user, nil
}

type Info struct {
	Result struct {
		Id          string   `json:"id"`
		UserName    string   `json:"userName"`
		Email       string   `json:"email"`
		PhoneNumber string   `json:"phoneNumber"`
		IsActive    bool     `json:"isActive"`
		Roles       []string `json:"roles"`
	} `json:"result"`
	Success   bool   `json:"success"`
	ErrorInfo string `json:"errorInfo"`
}

func GetUserInfoByToken(token string) {
	loginUri := BaseUrl + "/api/Authorize/userinfo"
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).Get(loginUri)
	if err != nil {
		log.Error(err)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		log.WithFields(log.Fields{"http status code": "resp.StatusCode()"}).Error("获取用户信息失败")
		return
	}

	var res Info
	if err := utils.Unmarshal[Info](resp.Body(), &res); err != nil {
		log.Error(err)
		return
	}
	fmt.Printf("%+v\n", res)
}
