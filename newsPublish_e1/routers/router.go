package routers

import (
	"newsPublish_e1/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*",beego.BeforeExec,filterFunc)

    beego.Router("/", &controllers.MainController{})
    //注册
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    //登陆
	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	//首页
	beego.Router("/article/index",&controllers.ArticleController{},"get:ShowIndex")
	//添加文章
	beego.Router("/article/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	//查看详情页
	beego.Router("/article/content",&controllers.ArticleController{},"get:ShowContent")
	//编辑文章
	beego.Router("/article/edit",&controllers.ArticleController{},"get:ShowEdit;post:HandleEdit")
	//删除文章
	beego.Router("/article/delete",&controllers.ArticleController{},"get:DeleteArticle")
	//添加类型
	beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAdd;post:HandleAdd")
	//退出
	beego.Router("/article/logout",&controllers.ArticleController{},"get:LogOut")

	beego.Router("/article/redis",&controllers.RedisGit{},"get:ShowRedis")

	beego.Router("/article/deleteType",&controllers.ArticleController{},"get:DeleteType")
}
func filterFunc(ctx*context.Context){
	userName:=ctx.Input.Session("userName")
	if userName==nil{
		ctx.Redirect(302,"/login")
		return
	}
}