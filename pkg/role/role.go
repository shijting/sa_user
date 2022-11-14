package role

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

const table = "roles"

func GetRoleByCode(code string) (models.Role, error) {
	var role models.Role
	queryStatement := fmt.Sprintf(`select id, name, code,description from %s where code = $1 limit 1`, table)
	err := pgxscan.Get(context.Background(), inits.GetDB(), &role, queryStatement, code)
	if err != nil && !pgxscan.NotFound(err) {
		return models.Role{}, err
	}

	return role, nil
}

func AddRoute(c *gin.Context) {
	var form models.Role
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg:    "添加角色失败",
			Detail: err.Error(),
		})
		inits.Log.WithFields(log.Fields{"action": "【role.AddRole】"}).Debug(err)
		return
	}
	ct, err := inits.GetDB().Exec(context.Background(),
		fmt.Sprintf(`insert into %s(name, code, description) values($1, $2, $3)`, table), form.Name, form.Code, form.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg:    "添加角色失败",
			Detail: "添加角色失败",
		})
		inits.Log.WithFields(log.Fields{"action": "【role.AddRole】"}).Error(err)
		return
	}
	if ct.RowsAffected() != 1 {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg:    "添加角色失败",
			Detail: "添加角色失败",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{Msg: "success"})
}

func DelRoute(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg: "id 不能为空",
		})
		inits.Log.WithFields(log.Fields{"action": "【role.delRole】"}).Debug(id)
		return
	}
	ct, err := inits.GetDB().Exec(context.Background(),
		fmt.Sprintf(`delect from %s where id = $1`, table), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg: "删除角色失败",
		})
		inits.Log.WithFields(log.Fields{"action": "【role.DelRole】"}).Error(err)
		return
	}
	if ct.RowsAffected() != 1 {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg:    "删除角色失败",
			Detail: "删除角色失败",
		})
		return
	}
	// TODO 删除用户角色表
	c.JSON(http.StatusOK, utils.Response{Msg: "success"})
}
