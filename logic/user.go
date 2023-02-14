package logic

import (
	"github.com/RaymondCode/simple-demo/dao/mysql"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/RaymondCode/simple-demo/pkg/jwt"
	"github.com/RaymondCode/simple-demo/pkg/snowflake"
)

func Register(p *models.ParamSignUp)( int64,string, error){
	//判断用户是否存在
	if err:=mysql.CheckUserExist(p.Username);err!=nil{
		return 0,"",err
	}
	//2.生成uid 然后构造user示例
	userid:=snowflake.GenID()
	user:=&models.User{
		UserID: userid,
		Username: p.Username,
		Password: p.Password,
	}
	//3.把用户插入到数据库中,需要返回用户的id
	if err:=mysql.InsertUser(user);err!=nil{
		return 0,"",nil
	}
	token,err:=jwt.GenToken(user.UserID, user.Username)
	if err!=nil{
		return 0,"",err
	}
	return userid,token,nil
}



func Login(p *models.ParamLogin) ( int64, string,  error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是指针，就能拿到user.UserID
	 userID,err := mysql.Login(user)
	 if err != nil {
		return 0,"", err
	}

	// 生成JWT
	 token,err:=jwt.GenToken(user.UserID, user.Username)
	 if err!=nil{
		return 0,"",err
	}
	return userID,token,nil
}

//主要是根据用户的id查询该用户的详细信息
func UserInfo(userid int64) ( models.User, error){
	user,err:=mysql.FindUserIdByUserId(userid)
	if err!=nil{
		return user,err
	}
	return user,nil
}
