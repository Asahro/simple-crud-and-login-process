package controllers

import (
	"lemonilo/libs"
	"lemonilo/models"
	"lemonilo/utils"
	"math"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	jwt "github.com/dgrijalva/jwt-go"
)

type ResponUsers struct {
	UserId    int    `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
	UpdateAt  string `json:"update_at"`
	UpdateBy  string `json:"update_by"`
}

type ResponUserList struct {
	UserId  int    `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type ResponLogin struct {
	SessionId string `json:"session_id"`
}

type UsersController struct {
	beego.Controller
}

func (c *UsersController) ReadUsers() {
	var resp utils.ResponseSchema
	resp.Code = 400
	resp.Data = nil
	var responUserListArray []ResponUserList
	var responUserList ResponUserList

	userID, _ := c.GetInt("user_id")
	limit, _ := c.GetInt("limit")
	page, _ := c.GetInt("offset")
	if limit < 1 {
		limit = 10
	}
	offset := page * limit

	userID = int(math.Abs(float64(userID)))
	limit = int(math.Abs(float64(limit)))
	offset = int(math.Abs(float64(offset)))

	if userID == 0 {
		data, err := models.ReadUsers(limit, offset)
		if err != nil {
			resp.Message = err.Error()
		} else {
			for i := 0; i < len(data); i++ {
				responUserList.UserId = data[i].UserId
				responUserList.Name = data[i].Name
				responUserList.Email = data[i].Email
				responUserList.Address = data[i].Address
				responUserListArray = append(responUserListArray, responUserList)
			}
			resp.Data = responUserListArray
			resp.Message = "success"
			resp.Code = 200
		}
	} else {
		if !models.IsUserExist(userID) {
			resp.Message = "user's data not found"
		} else {
			if data, err := models.ReadUserById(userID); err != nil {
				resp.Message = "user's data not found"
			} else {
				resp.Message = "success"
				resp.Code = 200
				resp.Data = data
			}
		}
	}

	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	c.Ctx.Output.Header("Access-Control-Allow-Method", "GET, OPTIONS")
	c.Ctx.Output.Header("Content-Type", "application/json, application/x-www-form-urlencoded")
	c.Ctx.Output.SetStatus(resp.Code)
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *UsersController) CreateUser() {
	var resp utils.ResponseSchema
	resp.Code = 400
	resp.Data = nil

	name := c.GetString("name")
	email := c.GetString("email")
	address := c.GetString("address")
	password := c.GetString("password")

	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.Required(email, "email")
	valid.Required(address, "address")
	valid.Required(password, "password")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			resp.Message = err.Key + " " + strings.ToLower(err.Message)
		}
	} else {
		if err := models.IsUserExistByEmail(email); err == false {
			var user models.Users
			user.Name = name
			user.Email = email
			user.Address = address
			user.Password, _ = libs.HashPassword(password)

			err := models.CreateUser(user)

			if err != nil {
				resp.Code = 500
				resp.Message = "failed"
				resp.Data = err
			} else {
				resp.Code = 200
				resp.Message = "success"
			}
		} else {
			resp.Message = "email sudah digunakan"
		}
	}

	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	c.Ctx.Output.Header("Access-Control-Allow-Method", "POST, OPTIONS")
	c.Ctx.Output.Header("Content-Type", "application/json, application/x-www-form-urlencoded")
	c.Ctx.Output.SetStatus(resp.Code)
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *UsersController) UpdateUser() {
	var resp utils.ResponseSchema
	resp.Code = 400
	resp.Data = nil

	userId, _ := c.GetInt("user_id")
	userId = int(math.Abs(float64(userId)))
	name := c.GetString("name")
	email := c.GetString("email")
	password := c.GetString("password")
	address := c.GetString("address")

	valid := validation.Validation{}
	valid.Required(userId, "user id")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			resp.Message = err.Key + " " + strings.ToLower(err.Message)
		}
	} else {
		if !models.IsUserExist(userId) {
			resp.Message = "user's data not found"
		} else {
			user := orm.Params{
				"user_id":  userId,
				"name":     name,
				"email":    email,
				"password": password,
				"address":  address,
			}

			for i, val := range user {
				if val == "" {
					delete(user, i)
				}
			}

			err := models.UpdateUser(user)

			if err != nil {
				resp.Code = 500
				resp.Message = err.Error()
			} else {
				resp.Code = 200
				resp.Message = "success"
			}
		}
	}

	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	c.Ctx.Output.Header("Access-Control-Allow-Method", "POST, OPTIONS")
	c.Ctx.Output.Header("Content-Type", "application/json, application/x-www-form-urlencoded")
	c.Ctx.Output.SetStatus(resp.Code)
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *UsersController) DeleteUser() {
	var resp utils.ResponseSchema
	resp.Code = 400
	resp.Data = nil

	userId, _ := c.GetInt("user_id")
	userId = int(math.Abs(float64(userId)))

	valid := validation.Validation{}
	valid.Required(userId, "user id")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			resp.Message = err.Key + " " + strings.ToLower(err.Message)
		}
	} else {
		if !models.IsUserExist(userId) {
			resp.Message = "user's data not found"
		} else {
			if err := models.DeleteUser(userId); err != nil {
				resp.Code = 500
				resp.Message = "failed"
			} else {
				resp.Code = 200
				resp.Message = "success"
			}
		}
	}

	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	c.Ctx.Output.Header("Access-Control-Allow-Method", "POST, OPTIONS")
	c.Ctx.Output.Header("Content-Type", "application/json, application/x-www-form-urlencoded")
	c.Ctx.Output.SetStatus(resp.Code)
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *UsersController) LoginUser() {
	var resp utils.ResponseSchema
	var respon ResponLogin
	resp.Code = 400
	resp.Data = nil

	email := c.GetString("email")
	password := c.GetString("password")

	valid := validation.Validation{}
	valid.Required(email, "email")
	valid.Required(password, "password")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			resp.Message = err.Key + " " + strings.ToLower(err.Message)
		}
	} else {
		if !models.IsUserExistByEmail(email) {
			resp.Message = "user's data not found"
		} else {
			if data, err := models.ReadUserByEmail(email); err != nil {
				resp.Message = "user's data not found"
			} else {
				login := libs.CheckPasswordHash(password, data.Password)
				if login == true {
					sign := jwt.New(jwt.GetSigningMethod("HS256"))
					claims := sign.Claims.(jwt.MapClaims)
					claims["user"] = email
					token, err := sign.SignedString([]byte("sahrojuara"))
					if err == nil {
						respon.SessionId = token
						resp.Code = 200
						resp.Message = "login berhasil"
						resp.Data = respon
					} else {
						resp.Message = "generate token gagal"
					}
				} else {
					resp.Message = "login gagal"
				}
			}
		}
	}

	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	c.Ctx.Output.Header("Access-Control-Allow-Method", "POST, OPTIONS")
	c.Ctx.Output.Header("Content-Type", "application/json, application/x-www-form-urlencoded")
	c.Ctx.Output.SetStatus(resp.Code)
	c.Data["json"] = resp
	c.ServeJSON()
}
