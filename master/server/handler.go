package server

import (
	"corntab/common"
	"corntab/master/model"
	"corntab/master/service"
	"corntab/master/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//处理函数

func pong(c *gin.Context) {

	resp := &common.Response{}
	resp.Code = 0
	resp.Msg = "pong"
	c.JSON(200, resp)
}

func saveTask(c *gin.Context) {

	req := model.SaveJobReq{}
	resp := common.NewResponse()

	logrus.Infof("client request saveTask")
	var err error
	//绑定参数
	if err = c.ShouldBind(&req); err != nil {
		logrus.Errorf("[saveTask] param is invalid,err_msg = %v", err)
		resp.Code, resp.Msg = errors.GetErr(errors.NewCTErr(errors.PARAMPARSEERROR))
		c.JSON(200, resp)
		return
	}

	resp.Data, err = service.SaveJob(&req)
	if err != nil {
		resp.Code, resp.Msg = errors.GetErr(err)
	}

	c.JSON(200, resp)
}

func deleteTask(c *gin.Context) {
	req := model.DeleteJobReq{}
	resp := common.NewResponse()

	logrus.Infof("client request deleteTask")

	var err error
	if err = c.ShouldBind(&req); err != nil {
		logrus.Errorf("[saveTask] param is invalid,err_msg = %v", err)
		resp.Code, resp.Msg = errors.GetErr(errors.NewCTErr(errors.PARAMPARSEERROR))
		c.JSON(200, resp)
		return
	}

	resp.Data, err = service.DeleteJob(&req)
	if err != nil {
		resp.Code, resp.Msg = errors.GetErr(err)
	}
	c.JSON(200, resp)
}

func queryTask(c *gin.Context) {
	req := model.QueryJobReq{}
	resp := common.NewResponse()

	logrus.Infof("client request queryTask")

	var err error
	if err = c.ShouldBind(&req); err != nil {
		logrus.Errorf("[queryTask] param is invalid,err_msg = %v", err)
		resp.Code, resp.Msg = errors.GetErr(errors.NewCTErr(errors.PARAMPARSEERROR))
		c.JSON(200, resp)
		return
	}

	resp.Data, err = service.ListJobs(&req)
	if err != nil {
		resp.Code, resp.Msg = errors.GetErr(err)
	}

	c.JSON(200, resp)
}
