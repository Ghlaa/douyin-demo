package main

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dao/mysql"
	"github.com/RaymondCode/simple-demo/dao/redis"
	"github.com/RaymondCode/simple-demo/logger"
	"github.com/RaymondCode/simple-demo/pkg/snowflake"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/settings"
	"os"
)

func main() {
	go service.RunMessageServer()

	if len(os.Args) < 2 {
		fmt.Println("need config file.eg: config/config.yaml")
		return
	}
	// 加载配置
	if err := settings.Init(os.Args[1]); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close() // 程序退出关闭数据库连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()
	if err:=snowflake.Init(settings.Conf.StartTime,settings.Conf.MachineID);err!=nil{
		fmt.Printf("init snowflake failed，err:=%v\n",err)
		return
	}

	r:=initRouter(settings.Conf.Mode)

	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err!=nil{
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
