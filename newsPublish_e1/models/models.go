package models

import (
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"

)

type User struct {
	Id int
	Name string 	`orm:"unique"`
	Psd string
	Articles []*Article `orm:"rel(m2m)"`
}
type Article struct {
	Id int `orm:"pk;auto"`
	Title string	`orm:"size(40)"`
	Content string	`orm:"size(400)"`
	ReadCount int	`orm:"default(0)"`
	Time time.Time	`orm:"type(datetime);auto_now_add"`
	Img string	`orm:"null"`
	Users []*User	`orm:"reverse(many)"`
	ArticleType *ArticleType `orm:"rel(fk);on_delete(set_null);null"`
}
type ArticleType struct {
	Id int
	TypeName string
	Articles []*Article `orm:"reverse(many)"`
}
func init()  {
	//注册数据库
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/mytest1?charset=utf8")

	//注册表
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	//跑起来
	orm.RunSyncdb("default",false,true)
}