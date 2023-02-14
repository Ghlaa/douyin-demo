package controller

import (
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/dao/mysql"
	"github.com/RaymondCode/simple-demo/logic"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync/atomic"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register_cop(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		atomic.AddInt64(&userIdSequence, 1)
		newUser := User{
			Id:   userIdSequence,
			Name: username,
		}
		usersLoginInfo[token] = newUser
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    username + password,
		})
	}
}

func Register(c *gin.Context) {
	//1.获取请求参数
	p:=new(models.ParamSignUp)
	username:=c.Query("username")
	password:=c.Query("password")
	if len(username)==0||len(password)==0{
		//请求参数出现错误
		zap.L().Error("register with invalid param")

		//ResponseError(c, CodeInvalidParam)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: CodeInvalidParam, StatusMsg: CodeInvalidParam.Msg()},
		})

		return
	}

	//2.业务逻辑处理
	p.Username=username
	p.Password=password
	userID,token,err:=logic.Register(p)
	if err!=nil{
		zap.L().Error("logic.Register failed",zap.Error(err))
		if errors.Is(err,mysql.ErrorUserExist){
			//ResponseError(c,CodeUserExist)
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: CodeUserExist, StatusMsg: CodeUserExist.Msg()},
			})
			return
		}
		//ResponseError(c,CodeServerBusy)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: CodeServerBusy, StatusMsg: CodeServerBusy.Msg()},
		})
		return
	}
	//返回响应
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode:CodeSuccess,StatusMsg: "register success"},
		UserId:   userID,
		Token:    token,
	})

	//c.Redirect(http.StatusMovedPermanently,"/douyin/user/login/")

}


func Login_copy(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
func Login(c *gin.Context) {
	// 1.获取请求参数及参数校验
	p:=new(models.ParamLogin)
	username:=c.Query("username")
	password:=c.Query("password")
	if len(username)==0||len(password)==0{
		//请求参数出现错误
		zap.L().Error("login with invalid param")

		//ResponseError(c, CodeInvalidParam)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: CodeInvalidParam, StatusMsg: CodeInvalidParam.Msg()},
		})
		return
	}
	// 2.业务逻辑处理
	p.Username=username
	p.Password=password
	userID,token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			//ResponseError(c, CodeUserNotExist)

			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: CodeUserNotExist, StatusMsg: CodeUserNotExist.Msg()},
				UserId: userID,
				Token: token,
			})
			return
		}
		//如果用户名存在还出现出错,那么肯定就是密码错误了
		c.JSON(http.StatusOK,UserLoginResponse{
			Response: Response{StatusCode: CodeInvalidPassword, StatusMsg: CodeInvalidPassword.Msg()},
			UserId: userID,
			Token: token,
		})
		return
	}

	// 3.返回响应
	c.JSON(http.StatusOK,UserLoginResponse{
		Response:Response{StatusCode: CodeSuccess,StatusMsg: "login success"},
		UserId: userID,
		Token: token,
	})

}

func UserInfo_copy(c *gin.Context) {
	fmt.Println("testsetttttsttttttttt")
	token := c.Query("token")
	fmt.Printf("token;lens",len(token))
	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {

	//1.参数校验：获取token和用户id，看用户是否已经登录
	token := c.Query("token")
	userId_str:=c.Query("user_id")
	userId_int,err:=strconv.Atoi(userId_str)
	userId:=int64(userId_int)
	if err!=nil{
		zap.L().Error("error with  strconving string to int64",zap.Error(err))
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: CodeOtherErr, StatusMsg: "error with  strconving string to int64"},
		})
		return
	}

	if len(token)==0||len(userId_str)==0{
		zap.L().Error("please login first before you query your information")
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: CodeNeedLogin, StatusMsg: CodeNeedLogin.Msg()},
		})
		return
	}
	if len(token)!=193{
		zap.L().Error("bad token with a wrong length")
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: CodeInvalidToken, StatusMsg: CodeInvalidToken.Msg()},
		})
		return
	}

	//2.业务逻辑处理
	user,err:=logic.UserInfo(userId)
	if err!=nil{
		zap.L().Error("logic.UserInfo failed",zap.Error(err))
		if errors.Is(err,mysql.ErrorUserNotExist){
			//ResponseError(c,CodeUserExist)
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: CodeUserNotExist, StatusMsg: CodeUserNotExist.Msg()},
			})
			return
		}
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: CodeServerBusy, StatusMsg: CodeServerBusy.Msg()},
		})
		return
	}
	//3.返回响应
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: CodeSuccess,StatusMsg: "query personal information success!!"},
		User:     User{
			Id: user.UserID,
			Name: user.Username,
			FollowCount: user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow: user.IsFollow,
		},
	})
}

