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
		c.JSON(200, resp)
		return
	}

	c.JSON(200, resp)
}
