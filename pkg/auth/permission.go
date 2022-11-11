package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/shjting0510/sa_user/pkg/roles_permissions"
	"github.com/shjting0510/sa_user/pkg/user"
	"github.com/shjting0510/sa_user/utils"
	"net/http"
	"strings"
)

// CheckPermission 检查权限
func CheckPermission(c *gin.Context) {
	action := c.Query("action")
	authorization := c.GetHeader("Authorization")
	tokenArr := strings.SplitN(authorization, " ", 2)

	if len(tokenArr) != 2 {
		c.JSON(http.StatusUnauthorized, utils.Response{
			Msg: "无效token",
		})
		return
	}
	if tokenArr[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, utils.Response{
			Msg: "无效token",
		})
		return
	}

	res, err := user.GetUserInfoByToken(tokenArr[1])
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:   0,
			Msg:    "获取权限失败",
			Detail: err.Error(),
		})
		return
	}

	ok, err := roles_permissions.CheckPermission(res.Result.Roles, action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:   0,
			Msg:    "检查权限失败",
			Detail: err.Error(),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.Response{
			Code: 0,
			Msg:  "无操作权限",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code: 0,
		Msg:  "权限通过",
	})
}
