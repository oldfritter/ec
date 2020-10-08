package initializers

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v2"

	. "ec/models"
	"ec/utils"
)

type ApiInterface struct {
	Method             string `yaml:"method"`
	Path               string `yaml:"path"`
	Auth               bool   `yaml:"auth"`
	Sign               bool   `yaml:"sign"`
	CheckFormat        bool   `yaml:"check_format"`
	CheckTimestamp     bool   `yaml:"check_timestamp"`
	LimitTrafficWithIp bool   `yaml:"limit_traffic_with_ip"`
}

var GlobalApiInterfaces []ApiInterface

func LoadInterfaces() {
	files, err := ioutil.ReadDir("config/interfaces/")
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, f := range files {
		if matched, err := regexp.MatchString(".yml$", f.Name()); matched && err == nil {
			path_str, _ := filepath.Abs("config/interfaces/" + f.Name())
			content, err := ioutil.ReadFile(path_str)
			if err != nil {
				log.Fatal(err)
				return
			}
			var interfaces []ApiInterface
			err = yaml.Unmarshal(content, &interfaces)
			if err != nil {
				log.Fatal(err)
			}
			GlobalApiInterfaces = append(GlobalApiInterfaces, interfaces...)
		}
	}
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		treatLanguage(context)

		var currentApiInterface ApiInterface
		for _, apiInterface := range GlobalApiInterfaces {
			if context.Path() == apiInterface.Path && context.Request().Method == apiInterface.Method {
				currentApiInterface = apiInterface
				if currentApiInterface.LimitTrafficWithIp && LimitTrafficWithIp(context) != true {
					return utils.BuildError("1027")
				}
				if apiInterface.Auth != true {
					return next(context)
				}
			}
		}

		params := make(map[string]string)
		for k, v := range context.QueryParams() {
			params[k] = v[0]
		}
		values, _ := context.FormParams()
		for k, v := range values {
			params[k] = v[0]
		}
		if currentApiInterface.Path == "" {
			return utils.BuildError("1025")
		}
		if context.Request().Header.Get("Authorization") == "" {
			return utils.BuildError("1026")
		} else {
			params["authorization"] = context.Request().Header.Get("Authorization")
		}
		params["device"] = context.Request().Header.Get("Device")
		if currentApiInterface.CheckTimestamp && !checkTimestamp(context, params) {
			return utils.BuildError("1024")
		}
		if currentApiInterface.Sign && !checkSign(context, params) {
			return utils.BuildError("1023")
		}
		var user User
		var err error
		if params["device"] == "" {
			user, err = normalAuth(params)
		} else {
			user, err = deviceAuth(params)
		}
		if err != nil {
			log.Println(err)
			return err
		}
		context.Set("current_user", user)
		return next(context)
	}
}

func normalAuth(params map[string]string) (user User, err error) {
	db := MainDbBegin()
	defer db.DbRollback()
	if db.Joins("INNER JOIN (tokens) ON (tokens.user_id = users.id)").
		Where("tokens.`type` = ?", "Login::Token").
		Where("tokens.token = ? AND ? < tokens.expire_at", params["authorization"], time.Now()).
		First(&user).RecordNotFound() {
		return user, utils.BuildError("1101")
	}
	return
}

func deviceAuth(params map[string]string) (user User, err error) {
	return
}
