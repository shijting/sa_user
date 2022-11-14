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

type permissionsForm struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type permission []models.Permission
type Res struct {
	List  permission `json:"list"`
	Count int        `json:"count"`
}

// f=
func GetPermissions(c *gin.Context) {
	var form permissionsForm
	if err := c.ShouldBind(&form); err != nil {
		inits.Log.WithFields(log.Fields{"action": "【permission.GetPermissions】", "data": form}).Debug(err)
		return
	}
	if form.Page == 0 {
		form.Page = 1
	}
	if form.PageSize == 0 || form.PageSize >= 500 {
		form.PageSize = 10
	}
	var res Res
	queryStatement := fmt.Sprintf(`select  id, name, code, description, action, parent_id from %s offset $1 limit $2`, table)
	err := pgxscan.Select(context.Background(), inits.GetDB(), &res.List, queryStatement, (form.Page-1)*form.PageSize, form.PageSize)
	if err != nil {
		inits.Log.Debug(err)
		inits.Log.WithFields(log.Fields{"action": "【permission.GetPermissions】"}).Error(err)
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg: "获取数据失败",
		})
		return
	}

	err = pgxscan.Get(context.Background(), inits.GetDB(), &res.Count, fmt.Sprintf(`select count(*) as count from %s`, table))
	if err != nil {
		inits.Log.Debug("err:", err)
		inits.Log.WithFields(log.Fields{"action": "【permission.GetPermissions】"}).Error(err)
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg: "获取数据失败",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Msg:  "ok",
		Data: res,
	})
}

func GetPermissionByCode(code string) (models.Permission, error) {
	var permission models.Permission
	queryStatement := fmt.Sprintf(`select id, name, code, description, action, parent_id from %s where code = $1 limit 1`, table)
	err := pgxscan.Get(context.Background(), inits.GetDB(), &permission, queryStatement, code)
	if err != nil && !pgxscan.NotFound(err) {
		return models.Permission{}, err
	}

	return permission, nil
}

// AddPermission 添加权限规则
func AddPermission(c *gin.Context) {
	var form models.Permission
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg:    "添加权限失败",
			Detail: err.Error(),
		})
		inits.Log.WithFields(log.Fields{"action": "【permission.AddPermission】"}).Debug(err)
		return
	}
	ct, err := inits.GetDB().Exec(context.Background(),
		fmt.Sprintf(`insert into %s(name, code, description, action, parent_id) values($1, $2, $3, $4, $5)`, table),
		form.Name, form.Code, form.Description, form.Action, form.ParentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg:    "添加权限失败",
			Detail: "添加权限失败",
		})
		inits.Log.WithFields(log.Fields{"action": "【permission.AddPermission】"}).Error(err)
		return
	}
	if ct.RowsAffected() != 1 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg: "添加权限失败",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{Msg: "添加权限成功"})
}

type EditPermissionForm struct {
	Id int `json:"id" form:"id" binding:"required"`
	models.Permission
	//Name        string `json:"name" form:"name"`
	//Code        string `json:"code" form:"code" binding:"required,min=1,max=30"`
	//Description string `json:"description"`
	//Action      string `json:"action" form:"action" binding:"required,min=1,max=50"`
	//ParentID    int    `json:"parent_id" form:"parent_id"`
}

func EditPermission(c *gin.Context) {
	var form EditPermissionForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg:    "修改权限失败",
			Detail: err.Error(),
		})
		inits.Log.WithFields(log.Fields{"action": "【permission.EditPermission】"}).Debug(err)
		return
	}
	parameters := make([]any, 0)
	sqlStatement := fmt.Sprintf(`UPDATE %s set name=$1, code=$2, action=$3, description = $4, parent_id = $5 where id = $6`, table)
	parameters = append(parameters, form.Name, form.Code, form.Action, form.Description, form.ParentID, form.Id)
	inits.Log.Info(sqlStatement, parameters)
	ct, err := inits.GetDB().Exec(context.Background(),
		sqlStatement,
		parameters...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg:    "修改权限失败",
			Detail: "修改权限失败",
		})
		inits.Log.WithFields(log.Fields{"action": "【permission.EditPermission】"}).Error(err)
		return
	}
	if ct.RowsAffected() != 1 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg: "修改权限失败",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{Msg: "修改权限成功"})
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
		fmt.Sprintf(`delete from %s where id = $1`, table), id)
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
