package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsPublish_e1/models"
)

type UserController struct {
	beego.Controller
}

func (this *UserController)ShowRegister(){
	this.TplName="register.html"
}
func (this *UserController)HandleRegister(){
	username:=this.GetString("userName")
	password:=this.GetString("password")
	o:=orm.NewOrm()
	var user models.User
	user.Name=username
	user.Psd=password
	o.Insert(&user)
	this.Redirect("/login",302)
}
func (this *UserController)ShowLogin(){
	username:=this.Ctx.GetCookie("userName")
	if username!=""{
		this.Data["userName"]=username
		this.Data["checked"]="checked"
	}else {
		this.Data["userName"]=""
		this.Data["checked"]=""
	}
	this.TplName="login.html"
}
func (this *UserController)HandleLogin(){
	username:=this.GetString("userName")
	password:=this.GetString("password")
	if username==""||password==""{
		beego.Error("不能为空")
		this.TplName="login.html"
		return
	}
	o:=orm.NewOrm()
	var user models.User
	user.Name=username
	err:=o.Read(&user,"name")
	if err!=nil{
		beego.Error("未读取到",err)
		this.TplName="login.html"
		return
	}
	remember:=this.GetString("remember")
	if remember=="on"{
		this.Ctx.SetCookie("userName",username,1000)
	}else{
		this.Ctx.SetCookie("userName",username,-1)
	}
	this.SetSession("userName",username)
	this.Redirect("/article/index",302)

}