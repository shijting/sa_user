package roles_permissions

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/shjting0510/sa_user/inits"
	"github.com/shjting0510/sa_user/models"
	"github.com/shjting0510/sa_user/pkg/permission"
	"github.com/shjting0510/sa_user/pkg/role"
	"github.com/shjting0510/sa_user/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const table = "roles_permissions"

func AddRolePermission(c *gin.Context) {
	var form models.RolePermission
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg:    "绑定权限失败",
			Detail: err.Error(),
		})
		inits.Log.WithFields(log.Fields{"action": "【role_permission.AddRolePermission】"}).Debug(err)
		return
	}

	_role, err := role.GetRoleByCode(form.RoleCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg: "绑定权限失败",
		})
		inits.Log.WithFields(log.Fields{"action": "【role_permission.AddRolePermission】"}).Error(err)
		return
	}
	if _role.Id == 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg: "角色不存在",
		})
		inits.Log.WithFields(log.Fields{"action": "【role_permission.AddRolePermission】"}).Debug(err)
		return
	}

	_permission, err := permission.GetPermissionByCode(form.PermissionCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg: "绑定权限失败",
		})
		inits.Log.WithFields(log.Fields{"action": "【role_permission.AddRolePermission】"}).Error(err)
		return
	}
	if _permission.Id == 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Msg: "权限不存在",
		})
		inits.Log.WithFields(log.Fields{"action": "【role_permission.AddRolePermission】"}).Debug(err)
		return
	}

	ct, err := inits.GetDB().Exec(context.Background(),
		fmt.Sprintf(`insert into %s(role_code, permission_code) values($1, $2)`, table), form.RoleCode, form.PermissionCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg: "绑定权限失败",
		})
		inits.Log.WithFields(log.Fields{"action": "【role_permission.AddRolePermission】"}).Error(err)
		return
	}
	if ct.RowsAffected() != 1 {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Msg: "绑定权限失败",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{Msg: "绑定权限成功"})
}

type RolePermissions struct {
	models.Role
	Permissions []models.Permission
}

func GetPermissionsByRoleCodes(roleCodes []string) ([]RolePermissions, error) {
	rolePermissions := make([]RolePermissions, 0)
	for _, roleCode := range roleCodes {
		var rps []models.RolePermission
		queryUser := fmt.Sprintf(`select id, role_code, permission_code from %s where role_code=$1`, table)
		err := pgxscan.Select(context.Background(), inits.GetDB(), &rps, queryUser, roleCode)
		if err != nil && !pgxscan.NotFound(err) {
			return nil, err
		}

		// 得到角色的所有权限
		permissions := make([]models.Permission, 0)
		for _, rolePermission := range rps {
			_permission := models.Permission{}
			queryUser := fmt.Sprintf(`select id, name, code, description, action, parent_id from %s where code=$1`, "permissions")
			err := pgxscan.Get(context.Background(), inits.GetDB(), &_permission, queryUser, rolePermission.PermissionCode)
			if err != nil && !pgxscan.NotFound(err) {
				fmt.Println("_permission:", err)
				return nil, err
			}
			permissions = append(permissions, _permission)
		}

		// 获取角色
		_role := models.Role{}
		err = pgxscan.Get(context.Background(), inits.GetDB(), &_role, fmt.Sprintf(`select id, name, code, description from %s where code=$1`, "roles"), roleCode)
		if err != nil && !pgxscan.NotFound(err) {
			return nil, err
		}
		rp := RolePermissions{
			Role:        _role,
			Permissions: permissions,
		}
		rolePermissions = append(rolePermissions, rp)
	}

	return rolePermissions, nil
}

func CheckPermission(roles []string, action string) (bool, error) {
	flag := false

	res, err := GetPermissionsByRoleCodes(roles)
	if err != nil {
		log.Error(err)
		return flag, err
	}
	rules := make([]string, 0)
	for _, role := range res {
		if len(role.Permissions) > 0 {
			for _, p := range role.Permissions {
				rules = append(rules, p.Action)
			}
		}
	}

	for _, rule := range rules {
		if rule == action {
			flag = true
			break
		}
	}
	return flag, nil
}
