package permission

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/shjting0510/sa_user/inits"
	"github.com/shjting0510/sa_user/models"
	"github.com/shjting0510/sa_user/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const table = "permissions"

func GetPermissionByCode(code string) (models.Permission, error) {
	var permission models.Permission
	queryUser := `select id, name, code, description, action, parent_id from users where code = $1 limit 1`
	err := pgxscan.Get(context.Background(), inits.GetDB(), &permission, queryUser, code)
	if err != nil && !pgxscan.NotFound(err) {
		return models.Permission{}, err
	}

	return permission, nil
}

func AddPermission(c *gin.Context) {
	var form models.Permission
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg:    "添加权限失败",
			Detail: err.Error(),
		})
		inits.Log.WithFields(log.Fields{"action": "【permission.permissions】"}).Debug(err)
		return
	}
	ct, err := inits.GetDB().Exec(context.Background(),
		fmt.Sprintf(`insert into %s(name, code, description, action, parent_id) values($1, $2, $3)`, table), form.Name, form.Code, form.Description, form.Action, form.ParentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg:    "添加权限失败",
			Detail: "添加权限失败",
		})
		inits.Log.WithFields(log.Fields{"action": "【permission.permissions】"}).Error(err)
		return
	}
	if ct.RowsAffected() != 1 {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg:    "添加权限失败",
			Detail: "添加权限失败",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{Msg: "添加权限成功"})
}

func DelPermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg: "id 不能为空",
		})
		inits.Log.WithFields(log.Fields{"action": "【permission.DelPermission】"}).Debug(id)
		return
	}
	ct, err := inits.GetDB().Exec(context.Background(),
		fmt.Sprintf(`delect from %s where id = $1`, table), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg: "删除权限失败",
		})
		inits.Log.WithFields(log.Fields{"action": "【permission.DelPermission】"}).Error(err)
		return
	}
	if ct.RowsAffected() != 1 {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg: "删除权限失败",
		})
		return
	}
	// TODO 删除角色权限表
	c.JSON(http.StatusOK, utils.Response{Msg: "删除权限成功"})
}
