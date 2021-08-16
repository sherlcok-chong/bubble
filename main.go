package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
)

var (
	DB *gorm.DB
)

//Todo model
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func initMySql() (err error) {
	DB, err = gorm.Open("mysql", "root:lll2002.11.22@tcp(127.0.0.1:3306)/bubble?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println("kkk")
		return err
	}
	return DB.DB().Ping()
}
func main() {
	//链接数据库
	err := initMySql()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	DB.AutoMigrate(&Todo{})
	r := gin.Default()

	r.Static("./static", "static")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})
	v1Group := r.Group("v1")
	{
		//待办事项

		//添加
		v1Group.POST("/todo", func(context *gin.Context) {
			//前端填写代办事项点击提交 请求到这里
			//1.取出请求中的数据
			todo := Todo{}
			context.BindJSON(&todo)
			//2.存入数据库
			if err := DB.Create(&todo).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"err": err,
				})
			} else {
				context.JSON(http.StatusOK, todo)
			}
			//3.返回响应
		})
		// 查看所有待办事项
		v1Group.GET("/todo", func(context *gin.Context) {
			todoList := make([]Todo, 0)
			if err := DB.Find(&todoList).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"err": err,
				})
			} else {
				context.JSON(http.StatusOK, todoList)
			}
		})
		//
		v1Group.GET("/todo/:id", func(context *gin.Context) {

		})
		//
		v1Group.PUT("/todo/:id", func(context *gin.Context) {
			id, _ := context.Params.Get("id")
			todo := Todo{}
			if err := DB.Where("id=?", id).First(&todo); err != nil {
				context.JSON(http.StatusOK, gin.H{"err": err})
			}
			context.BindJSON(&todo)
			if err := DB.Save(&todo).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{
					"err": err,
				})
			} else {
				context.JSON(http.StatusOK,todo)
			}

		})
		//
		v1Group.DELETE("todo/:id", func(context *gin.Context) {
			id, _ := context.Params.Get("id")
			todo := Todo{}
			if err := DB.Where("id=?", id).Delete(todo).Error; err != nil {
				context.JSON(http.StatusOK, gin.H{"err": err})
			} else {
				context.JSON(http.StatusOK,gin.H{id : "deleted"})
			}

		})
	}
	r.Run(":9090")
}
