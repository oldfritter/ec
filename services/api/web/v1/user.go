package v1

import (
	"net/http"

	"github.com/jinzhu/gorm"
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

	token := Token{UserId: int(user.ID), RemoteIp: c.RealIP()}
	db.Create(&token)
	db.DbCommit()
	user.Tokens = append(user.Tokens, &token)
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", Friend).Where("fs.user_id = ?", user.ID).Find(&user.Friends)
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", PendingFriend).Where("fs.user_id = ?", user.ID).Find(&user.PendingFriends)
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}

func UserInfo(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	var user User
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
	user := c.Get("current_user").(User)
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", Friend).Where("fs.user_id = ?", user.ID).Find(&user.Friends)
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", PendingFriend).Where("fs.user_id = ?", user.ID).Find(&user.PendingFriends)
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

func UserFriendAdd(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	user := c.Get("current_user").(User)
	if params["sn"] == "" {
		return utils.BuildError("2005")
	}
	db := MainDbBegin()
	defer db.DbRollback()
	var friend User
	if db.Where("sn = ?", params["sn"]).First(&friend).RecordNotFound() {
		return utils.BuildError("2001")
	}
	db.Save(&FriendShip{UserId: user.ID, FriendId: friend.ID, State: Friend})
	db.Save(&FriendShip{UserId: friend.ID, FriendId: user.ID})
	db.DbCommit()
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", Friend).Where("fs.user_id = ?", user.ID).Find(&user.Friends)
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", PendingFriend).Where("fs.user_id = ?", user.ID).Find(&user.PendingFriends)
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}

func UserFriendAccept(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	user := c.Get("current_user").(User)
	if params["sn"] == "" {
		return utils.BuildError("2005")
	}
	db := MainDbBegin()
	defer db.DbRollback()
	var friend User
	if db.Where("sn = ?", params["sn"]).First(&friend).RecordNotFound() {
		return utils.BuildError("2001")
	}
	var fs FriendShip
	db.Where("user_id = ?", user.ID).Where("friend_id = ?", friend.ID).First(&fs)
	fs.State = Friend
	db.Save(&fs)
	db.DbCommit()
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", Friend).Where("fs.user_id = ?", user.ID).Find(&user.Friends)
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", PendingFriend).Where("fs.user_id = ?", user.ID).Find(&user.PendingFriends)
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}

func UserFriendBlock(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	user := c.Get("current_user").(User)
	if params["sn"] == "" {
		return utils.BuildError("2005")
	}
	db := MainDbBegin()
	defer db.DbRollback()
	var friend User
	if db.Where("sn = ?", params["sn"]).First(&friend).RecordNotFound() {
		return utils.BuildError("2001")
	}
	var fs, ufs FriendShip
	db.Where("user_id = ?", user.ID).Where("friend_id = ?", friend.ID).First(&fs)
	fs.State = Blocked
	db.Save(&fs)
	db.Where("user_id = ?", friend.ID).Where("friend_id = ?", user.ID).First(&ufs)
	ufs.State = Blocked
	db.Save(&ufs)
	db.DbCommit()
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", Friend).Where("fs.user_id = ?", user.ID).Find(&user.Friends)
	config.MainDb.Joins("INNER JOIN (friend_ships as fs) ON (fs.friend_id = users.id)").
		Preload("PublicKeys", func(db *gorm.DB) *gorm.DB {
			return db.Order("public_keys.index ASC")
		}).Where("fs.state = ?", PendingFriend).Where("fs.user_id = ?", user.ID).Find(&user.PendingFriends)
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}
