package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)
type RedisGit struct {
	beego.Controller
}

func (this *RedisGit)ShowRedis(){
	//连接数据库
	conn,err:=redis.Dial("tcp",":6379")
	if err!=nil{
		beego.Error("redis连接错误",err)
		return
	}
	defer conn.Close()
	//操作数据库
	//conn.Send("set","kk","vv")
	//conn.Flush()
	//conn.Receive()
	resp,_:=conn.Do("mget","kk","k2","ll")
	//回复助手函数
	result,_:=redis.Values(resp,err)
	var kk,k2 string
	var ll int
	redis.Scan(result,&kk,&k2,&ll)
	beego.Info(kk,k2,ll)

}
