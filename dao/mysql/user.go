package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/RaymondCode/simple-demo/models"
)
const secret = "guanhaoliang"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
)

func Register()(err error){
	return err
}
func CheckUserExist(username string)(err error){
	sqlStr := `select count(user_id) from user where username = ?`
	var count int64
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

func InsertUser(user *models.User)(err error){
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	// 执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	//fmt.Println("Testests",user.UserID,user.Username,user.Password)
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	if err!=nil{
		return err
	}
	return
}
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (userID int64,err error) {
	oPassword := user.Password // 用户登录的密码
	sqlStr := `select user_id, username, password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return 0,ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		return 0,err
	}
	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return 0,ErrorInvalidPassword
	}
	return user.UserID,nil
}

func FindUserIdByUserId(userid int64)(user models.User,err error){
	sqlStr := `select count(user_id) from user where user_id = ?`
	var count int64
	if err := db.Get(&count, sqlStr, userid); err != nil {
		return user,err
	}
	if count == 0 {
		return user,ErrorUserNotExist
	}
	sqlStr = `select username,follow_count,follower_count,is_follow from user where user_id = ?`
	if err := db.Get(&user, sqlStr, userid); err != nil {
		return user,err
	}
	return user,nil
}