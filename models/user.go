package models

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/orm"
)

type Users struct {
	UserId    int `orm:"column(user_id);pk"`
	Name      string
	Email     string
	Address   string
	Password  string
	CreatedAt string
	CreatedBy string
	UpdateAt  string
	UpdateBy  string
}

func init() {
	orm.RegisterModel(new(Users))
}

func IsUserExist(userID int) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("users").Filter("user_id", userID).Exist()
	return exist
}

func IsUserExistByEmail(email string) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("users").Filter("email", email).Exist()
	return exist
}

func ReadUsers(limit, offset int) ([]Users, error) {
	var v []Users
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("user_id", "name", "email", "address", "password", "created_by", "created_at", "update_at", "update_by").From("users").Limit(limit).Offset(offset)
	sql := qb.String()
	o := orm.NewOrm()
	fmt.Println(sql)
	if _, err := o.Raw(sql).QueryRows(&v); err != nil {
		fmt.Println(err)
		return nil, err
	}
	if len(v) == 0 {
		return nil, errors.New("empty data")
	}
	return v, nil
}

func ReadUserById(userID int) (Users, error) {
	o := orm.NewOrm()
	var v Users
	if err := o.QueryTable("users").Filter("user_id", userID).One(&v); err != nil {
		fmt.Println()
		return Users{}, err
	}
	return v, nil
}

func ReadUserByEmail(email string) (Users, error) {
	o := orm.NewOrm()
	var v Users
	if err := o.QueryTable("users").Filter("email", email).One(&v); err != nil {
		fmt.Println()
		return Users{}, err
	}
	return v, nil
}

func CreateUser(data Users) error {
	o := orm.NewOrm()
	_, err := o.Insert(&data)
	return err
}

func UpdateUser(user orm.Params) error {
	o := orm.NewOrm()
	if _, err := o.QueryTable("users").Filter("user_id", user["user_id"]).Update(user); err != nil {
		return err
	}
	return nil
}

func DeleteUser(userID int) error {
	o := orm.NewOrm()
	if _, err := o.QueryTable("users").Filter("user_id", userID).Delete(); err != nil {
		return err
	}
	return nil
}
