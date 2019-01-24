package controllers

import (
	"github.com/astaxie/beego"
	_"newsPublish_e1/models"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"newsPublish_e1/models"
	"math"
	"github.com/gomodule/redigo/redis"
	"encoding/gob"
	"bytes"
)

type ArticleController struct{
	beego.Controller
}
//展示首页
func (this *ArticleController)ShowIndex(){
	username:=this.GetSession("userName")
	if username==nil{
		this.Redirect("/article/login",302)
		return
	}
	//获取数据库中的数据
	o:=orm.NewOrm()
	var articels []models.Article
	qs:=o.QueryTable("Article")
	//qs.All(&articels)

	//每页显示数
	pageSize:=10
	//处理首末页码
	pageIndex,err:=this.GetInt("pageIndex")
	if err!=nil{
		pageIndex=1
	}
	start:=pageSize*(pageIndex-1)


	var articleTypes []models.ArticleType

	//存入redis中
	conn,err:=redis.Dial("tcp",":6379")
	if err!=nil{
		beego.Error("redis连接错误",err)
		return
	}
	defer conn.Close()
	//conn.Do("set","articleTypes",articleTypes)
	data,_:=redis.Bytes(conn.Do("get","articleTypes"))
	beego.Info(data)
	if len(data)==0{
		//从数据库获取类型
		o.QueryTable("ArticleType").All(&articleTypes)
		beego.Info("从数据库获取数据")
		//序列化和反序列化
		//字节容器
		var buffer bytes.Buffer
		//获取一个编码器
		encode:=gob.NewEncoder(&buffer)
		//编码
		encode.Encode(&articleTypes)

		conn.Do("set","articleTypes",buffer.Bytes())
	}else {

		//获取需要解码的数据
		//data,_:=redis.Bytes(conn.Do("get","articleTypes"))
		//获取解码器
		beego.Info("从redis获取数据库")
		dec:=gob.NewDecoder(bytes.NewBuffer(data))
		//解码

		dec.Decode(&articleTypes)
		beego.Info(articleTypes)
	}




	this.Data["articleTypes"]=articleTypes

	typeName:=this.GetString("select")
	beego.Info(typeName)
	var count int64
	if typeName==""{
		//记录数
		count,_=qs.RelatedSel("ArticleType").Count()
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articels)
	}else{
		count,_=qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articels)
	}

	//页数
	pageCount:=math.Ceil(float64(count)/float64(pageSize))

	this.Data["count"]=count
	this.Data["pageCount"]=pageCount


	this.Data["typeName"]=typeName
	this.Data["pageIndex"]=pageIndex
	this.Data["articles"]=articels
	this.Layout="layout.html"
	this.TplName="index.html"
}

//展示增加文章页面
func (this *ArticleController)ShowAddArticle()  {
	//添加文章类型显示
	o:=orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	this.Data["articleTypes"]=articleTypes
	this.Layout="layout.html"
	this.TplName="add.html"
}
//处理增加文章
func (this *ArticleController)HandleAddArticle(){

	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	file,head,err:=this.GetFile("uploadname")
	if err!=nil||articleName==""||content==""{
		beego.Error("获取出错",err)
		this.TplName="add.html"
		return
	}
	defer file.Close()

	//判断大小
	if head.Size>5000000{
		beego.Error("图片太大",err)
		this.TplName="add.html"
		return
	}
	//判断格式
	ext:=path.Ext(head.Filename)
	if ext!=".jgp"&&ext!=".png"{
		beego.Error("格式不正确",err)
		this.TplName="add.html"
		return
	}
	//去重
	filename :=time.Now().Format("2006-01-02 15:04:05")

	//操作数据
	this.SaveToFile("uploadname","./static/img/"+filename+ext)

	//操作数据库
	o:=orm.NewOrm()
	var article models.Article
	article.Title=articleName
	article.Content=content
	article.Img="./static/img/"+filename+ext
	//获取类型
	typeName:=this.GetString("select")
	beego.Info(typeName)
	var articleType models.ArticleType
	articleType.TypeName=typeName
	err=o.Read(&articleType,"TypeName")
	if err!=nil{
		beego.Error("读取出错")
		return
	}
	article.ArticleType=&articleType

	o.Insert(&article)

	//返回数据
	this.Redirect("/article/index",302)

}

//展示详情页
func  (this *ArticleController)ShowContent()  {
	id,err:=this.GetInt("Id")
	if err!=nil{
		beego.Error("未获取到ID")
		this.TplName="index.html"
		return
	}

	o:=orm.NewOrm()
	/*var article models.Article
	article.Id=id*/
	article:= models.Article{Id:id}
	o.Read(&article)

	//插入浏览信息
	userName:=this.GetSession("userName")
	m2m:=o.QueryM2M(&article,"Users")
	user:= models.User{Name:userName.(string)}
	o.Read(&user,"Name")
	m2m.Add(user)
	//显示浏览信息方法1
	//o.LoadRelated(&article,"Users")
	//显示浏览信息方法2
	var users []models.User

	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)

	/*var articles []models.Article
	o.QueryTable("Article").Filter("Users__User__Name",userName).Distinct().All(&articles)
	this.Data["articles"]=articles
	beego.Info(articles)*/
	this.Data["users"]=users
	//阅读数增加
	article.ReadCount+=1
	o.Update(&article)

	this.Data["article"]=article
	this.Layout="layout.html"
	this.TplName="content.html"
}

//展示编辑页面
func (this *ArticleController)ShowEdit(){
	id,err:=this.GetInt("Id")
	if err!=nil{
		beego.Error("未获取到Id")
		this.Layout="layout.html"
		this.TplName="index.html"
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	o.Read(&article)
	this.Data["article"]=article
	this.Layout="layout.html"
	this.TplName="update.html"
}
//封装接口：上传文件
func UploadFunc(this *ArticleController,uploadname string)string{
	file,head,err:=this.GetFile(uploadname)
	if err!=nil{
		beego.Error("有错误")
		return ""
	}
	defer file.Close()
	//文件大小校验
	if head.Size>5000000{
		beego.Error("超出大小")
		return ""
	}
	//文件类型校验
	ext:=path.Ext(head.Filename)
	if ext!=".png" && ext!=".jpg"{
		beego.Error("文件类型有错误")
		return ""
	}
	//去重名校验
	filename:=time.Now().Format("2006-01-02 15:04:05")
	this.SaveToFile(uploadname,"./static/img/"+filename+ext)

	return "/static/img/"+filename+ext

}
//编辑页面
func (this *ArticleController)HandleEdit(){
	//获取数据
	id,err:=this.GetInt("Id")
	if err!=nil{
		beego.Error("获取id出错")
		return
	}
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	filepath:=UploadFunc(this,"uploadname")
	if articleName==""||content==""||filepath==""{
		beego.Error("不能为空！")
		this.Layout="layout.html"
		this.TplName="update.html"
		return
	}
	//更新
	o:=orm.NewOrm()
	var article models.Article
	//指定查询条件
	article.Id=id
	err=o.Read(&article)
	if err!=nil{
		beego.Error("没有读取到")
		return
	}
	article.Title=articleName
	article.Content=content
	article.Img=filepath

	_,err=o.Update(&article)
	if err!=nil{
		beego.Error("更新出错")
		return
	}
	this.Redirect("/article/index",302)
}
//删除文章
func (this *ArticleController)DeleteArticle(){
	id,err:=this.GetInt("Id")
	if err!=nil{
		beego.Error("获取ID出错")
		this.TplName="index.html"
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	o.Delete(&article)
	this.Redirect("/index",302)
}

//展示添加类型
func (this *ArticleController)ShowAdd(){
	o:=orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	this.Data["articleTypes"]=articleTypes
	this.Layout="layout.html"
	this.TplName="addType.html"
}
//处理添加类型
func (this *ArticleController)HandleAdd(){
	typeName:=this.GetString("typeName")
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName=typeName
	o.Insert(&articleType)
	this.Redirect("/article/addType",302)

}
//退出登陆
func(this *ArticleController)LogOut(){
	this.DelSession("userName")
	this.Redirect("/article/login",302)
}
//删除类型
func (this *ArticleController)DeleteType(){
	id,err:=this.GetInt("id")
	if err!=nil{
		beego.Error("获取id出错")
		return
	}
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id=id
	o.Delete(&articleType)
	this.Redirect("/article/addType",302)
}

