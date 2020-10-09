package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"ec/config"
	. "ec/models"
	"ec/services/api/helpers"
	"ec/utils"
)

func UserLogin(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	db := MainDbBegin()
	defer db.DbRollback()
	var user User
	if db.Joins("INNER JOIN (identities) ON (identities.user_id = users.id)").
		Where("identities.source = ?", params["source"]).
		Where("identities.symbol = ?", params["symbol"]).
		Preload("PublicKeys").
		Preload("Groups").
		First(&user).RecordNotFound() {
		return utils.BuildError("2002")
	} else {
		user.Password = params["password"]
	}
	if b := user.CompareHashAndPassword(); !b {
		return utils.BuildError("2002")
	}
	// TODO:
	// notify user

	db.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys").Where("fs.user_id = ?", user.ID).Find(&user.Friends)
	token := Token{UserId: int(user.ID), RemoteIp: c.RealIP()}
	db.Create(&token)
	db.DbCommit()
	user.Tokens = append(user.Tokens, &token)
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}

func UserInfo(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	db := config.MainDb
	var user User
	if db.Where("sn = ?", params["sn"]).
		Preload("PublicKeys").
		First(&user).RecordNotFound() {
		return utils.BuildError("2001")
	}
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}

func UserRegister(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	db := MainDbBegin()
	defer db.DbRollback()
	var identity Identity
	var user User
	if db.Where("`source` = ?", params["source"]).
		Where("symbol = ?", params["symbol"]).
		First(&identity).RecordNotFound() {
		user.Nickname = params["nickname"]
		user.Password = params["password"]
		identity.Source = params["source"]
		identity.Symbol = params["symbol"]
		db.Save(&user)
		identity.User = user
		token := Token{UserId: int(user.ID), RemoteIp: c.RealIP()}
		db.Create(&token)
		user.Tokens = append(user.Tokens, &token)
		db.Save(&identity)
		db.DbCommit()

		response := utils.SuccessResponse
		response.Body = user
		return c.JSON(http.StatusOK, response)
	}
	return utils.BuildError("2003")
}
