package models

import (
	"errors"
	"fmt"
	"login_golang/tables"
	"net/http"
	"strconv"

	// m "snapin-user-api/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"gorm.io/gorm"
)

func SchemaUsr(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("public" + "." + tableName)
	}
}

func SchemaFrm(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("frm" + "." + tableName)
	}
}

type UserModels interface {
	InsertUser(user tables.Users) (tables.Users, error)
	// InsertRoleUser(user tables.RoleUser) (tables.RoleUser, error)
	GetUserRows(user tables.Users) ([]tables.Users, error)
	GetUserRow(user tables.Users) (tables.UserData, error)
	GetUserWhereRow(user tables.Users, whreStr string) (tables.UserData, error)
	UpdateUser(userId int, user tables.UserMember) (tables.UserMember, error)
	DeleteUser(id int, authorID int) (bool, error)
	GetGenderRows(user tables.Genders) ([]tables.Genders, error)
	InsertUserInvites(user tables.UserInvites) (tables.UserInvites, error)
	GetUserInvites(user tables.UserInvites) (tables.UserInvitesOrg, error)
	DeleteUserInvites(id int) (bool, error)
}

type connection struct {
	db *gorm.DB
	r  *http.Request
	mw *jwt.GinJWTMiddleware
}

func NewUserModels(dbg *gorm.DB) UserModels {
	return &connection{
		db: dbg,
	}
}

func (con *connection) InsertUser(data tables.Users) (tables.Users, error) {
	err := con.db.Scopes(SchemaUsr("users")).Create(&data).Error
	if err != nil {
		fmt.Println(err)
		return tables.Users{}, err
	}
	return data, nil
}

func Paginate(r *http.Request) func(con *connection) *gorm.DB {
	return func(con *connection) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return con.db.Offset(offset).Limit(pageSize)
	}
}

func (con *connection) GetUserRows(fields tables.Users) ([]tables.Users, error) {
	var data []tables.Users

	// pageSize := 3
	// offset := (1 - 1) * pageSize
	// err := con.db.Scopes(SchemaUsr("users")).Where(fields).Offset(offset).Limit(pageSize).Find(&data).Error
	// if !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return data, err
	// }

	// offset := (page - 1) * pageSize
	// return db.Offset(offset).Limit(pageSize)

	// sql := con.db.Scopes(helpers.Paginate(con.r)).Find(&data)

	err := con.db.Scopes(SchemaUsr("users")).Where(fields).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *connection) GetUserRow(fields tables.Users) (tables.UserData, error) {
	var data tables.UserData
	// err := con.db.Scopes(SchemaUsr("users")).Where(fields).First(&data).Error

	err := con.db.Table("usr.users").Select("users.id, users.name, users.phone, users.email, users.avatar, users.date_of_birth, users.gender_id, t.translation as gender_name, users.remember_token,users.password, upper(users.encrypt_code) as encrypt_code, users.created_at, users.role_id").Joins("left join mstr.genders g on g.id = users.gender_id").Joins("left join mstr.translations t on t.textcontent_id = g.name_textcontent_id").Where(fields).First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tables.UserData{}, nil
	}

	if err != nil {
		fmt.Println("errrr no sintax:::: -------", err)
		return tables.UserData{}, err
	}

	return data, nil
}

func (con *connection) GetUserWhereRow(fields tables.Users, whreStr string) (tables.UserData, error) {
	var data tables.UserData
	// err := con.db.Scopes(SchemaUsr("users")).Where(fields).First(&data).Error

	err := con.db.Table("public.users").Select("users.id, users.name, users.phone, users.email, users.created_at").Where(fields).Where(whreStr).First(&data).Error
	if err != nil {
		return tables.UserData{}, err
	}
	return data, nil
}

func (con *connection) UpdateUser(userID int, fields tables.UserMember) (tables.UserMember, error) {
	var data tables.UserMember
	err := con.db.Scopes(SchemaUsr("users")).Where("id = ?", userID).Updates(fields).Error

	err = con.db.Scopes(SchemaUsr("users")).Where("id = ?", userID).Update("status", fields.Status).Error

	return data, err
}

func (con *connection) DeleteUser(id int, authorID int) (bool, error) {
	var data tables.Users
	err := con.db.Scopes(SchemaUsr("users")).Where("id = ?", id).Delete(&data).Error
	if err != nil {
		return false, err
	}

	err = con.db.Scopes(SchemaUsr("users")).Where("id = ?", id).Update("deleted_by", authorID).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *connection) GetGenderRows(fields tables.Genders) ([]tables.Genders, error) {
	var data []tables.Genders

	sql := con.db.Table("mstr.genders").Select("genders.id, genders.code, t.translation").Joins("join mstr.translations t on t.textcontent_id = genders.name_textcontent_id").Where("t.language_id = ?", 1).Where(fields).Find(&data)
	if !errors.Is(sql.Error, gorm.ErrRecordNotFound) {
		return data, sql.Error
	}
	return data, nil
}

func (con *connection) InsertUserInvites(data tables.UserInvites) (tables.UserInvites, error) {

	err := con.db.Scopes(SchemaUsr("user_invites")).Create(&data).Error
	if err != nil {
		fmt.Println(err)
		return tables.UserInvites{}, err
	}
	return data, nil
}

func (con *connection) GetUserInvites(fields tables.UserInvites) (tables.UserInvitesOrg, error) {
	var data tables.UserInvitesOrg

	err := con.db.Table("usr.user_invites").Select("user_invites.id,user_invites.user_id, user_invites.phone, o.id as organization_id").Joins("left join usr.users u on u.id = user_invites.user_id").Joins("left join mstr.organizations o on o.created_by = u.id").Where(fields).Find(&data).Error

	// err := con.db.Scopes(SchemaUsr("user_invites")).Where(fields).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

// func (con *connection) GenerateToken(fields tables.Users) (tables.UserInvitesOrg, error) {

// 	authMiddleware := m.SetupMiddleware(con.db)

// 	var data map[string]interface{}
// 	data = map[string]interface{}{
// 		"id": fields.ID,
// 	}
// 	token, time, err := authMiddleware.TokenGenerator(data)
// 	fmt.Println(time, err)

// 	return token, nil
// }

func (con *connection) DeleteUserInvites(id int) (bool, error) {
	var data tables.Users
	err := con.db.Scopes(SchemaUsr("user_invites")).Where("id = ?", id).Delete(&data).Error
	if err != nil {
		return false, err
	}

	// err = con.db.Scopes(SchemaUsr("user_invites")).Where("id = ?", id).Update("deleted_by", authorID).Error
	// if err != nil {
	// 	return false, err
	// }

	return true, err
}
