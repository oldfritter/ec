package v1

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	"ec/config"
	"ec/models"
	"ec/services/api/helpers"
	"ec/utils"
)

func UserLogin(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	db := models.MainDbBegin()
	defer db.DbRollback()
	var user models.User
	if db.Joins("INNER JOIN (identities) ON (identities.user_id = users.id)").
		Where("identities.source = ?", params["source"]).
		Where("identities.symbol = ?", params["symbol"]).
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
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", 1).
		Where("fs.user_id = ?", user.ID).Find(&user.Friends)
	token := models.Token{UserId: int(user.ID), RemoteIp: c.RealIP()}
	db.Create(&token)
	db.DbCommit()
	user.Tokens = append(user.Tokens, &token)
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}

func UserInfo(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	var user models.User
	if config.MainDb.Where("sn = ?", params["sn"]).
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).First(&user).RecordNotFound() {
		return utils.BuildError("2001")
	}
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}

func UserMe(c echo.Context) (err error) {
	user := c.Get("current_user").(models.User)
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys").Where("fs.state = ?", 1).Where("fs.user_id = ?", user.ID).Find(&user.Friends)
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}

func UserRegister(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	db := models.MainDbBegin()
	defer db.DbRollback()
	var identity models.Identity
	var user models.User
	if db.Where("`source` = ?", params["source"]).
		Where("symbol = ?", params["symbol"]).
		First(&identity).RecordNotFound() {
		user.Nickname = params["nickname"]
		user.Password = params["password"]
		identity.Source = params["source"]
		identity.Symbol = params["symbol"]
		db.Save(&user)
		identity.User = user
		token := models.Token{UserId: int(user.ID), RemoteIp: c.RealIP()}
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

func UserFriendAdd(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	user := c.Get("current_user").(models.User)
	if params["friend"] == "" {
		return utils.BuildError("2005")
	}
	db := models.MainDbBegin()
	defer db.DbRollback()
	var friend models.User
	if db.Where("sn = ?", params["friend"]).First(&friend).RecordNotFound() {
		return utils.BuildError("2001")
	}
	db.Save(&models.FriendShip{UserId: user.ID, FriendId: friend.ID, State: 1})
	db.Save(&models.FriendShip{UserId: friend.ID, FriendId: user.ID})
	db.DbCommit()
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", 1).
		Where("fs.user_id = ?", user.ID).Find(&user.Friends)
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}

func UserFriendAccept(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	user := c.Get("current_user").(models.User)
	if params["friend"] == "" {
		return utils.BuildError("2005")
	}
	db := models.MainDbBegin()
	defer db.DbRollback()
	var friend models.User
	if db.Where("sn = ?", params["friend"]).First(&friend).RecordNotFound() {
		return utils.BuildError("2001")
	}
	var fs models.FriendShip
	db.Where("user_id = ?", user.ID).Where("friend_id = ?", friend.ID).First(&fs)
	fs.State = 1
	db.Save(&fs)
	db.DbCommit()
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", 1).
		Where("fs.user_id = ?", user.ID).Find(&user.Friends)
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}
