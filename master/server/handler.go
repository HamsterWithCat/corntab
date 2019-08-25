package server

import (
	"github.com/gin-gonic/gin"
)
//处理函数

func saveTask(c *gin.Context){
	c.JSON(200,gin.H{
		"status":"",
		"msg":"",
	})
}