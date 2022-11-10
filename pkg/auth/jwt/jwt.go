package jwt

import (
	"github.com/shjting0510/sa_user/utils"
	"time"
)

func GenToken(userID int64) (string, error) {
	cc := utils.NewCustomClaims(userID, utils.WithExpire(time.Hour*24*15))
	return cc.GenToken()
}

func ParseToken(token string) (int64, bool) {
	cc, err := utils.ParseToken(token)
	if err != nil {
		// todo
		return 0, false
	}
	return cc.UserID, cc.UserID > 0
}
