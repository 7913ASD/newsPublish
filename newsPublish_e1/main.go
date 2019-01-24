package main

import (
	_ "newsPublish_e1/routers"
	"github.com/astaxie/beego"
	_"newsPublish_e1/models"
)

func main() {
	beego.AddFuncMap("prePage",preFunc)
	beego.AddFuncMap("nextPage",nextFunc)
	beego.Run()
}
func preFunc(pageIndex int) int {
	if pageIndex<=1{
		return 1
	}
	return pageIndex-1
}
func nextFunc(pageCount float64,pageIndex int)int{
	if pageIndex>=int(pageCount){
		return int(pageCount)
	}
	return pageIndex+1
}


