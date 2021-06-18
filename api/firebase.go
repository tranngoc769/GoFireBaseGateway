package api

import (
	"encoding/json"
	"go-firebase-gateway/common/model"
	"go-firebase-gateway/common/response"
	mdw "go-firebase-gateway/internal/middleware"
	"go-firebase-gateway/service"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type FireBaseHandler struct {
	FireBaseService service.FireBaseService
}

func NewFireBaseHandler(r *gin.Engine, FireBaseService service.FireBaseService) {
	handler := &FireBaseHandler{
		FireBaseService: FireBaseService,
	}
	v1 := r.Group("v1/firebase")
	{
		v1.GET("", mdw.AuthMiddleware(), handler.GetFireBases)
		v1.POST("/event", mdw.AuthMiddleware(), handler.CreateEventFirebase)
		v1.POST("/truncate", mdw.AuthMiddleware(), handler.Truncate)
	}
}
func (handler *FireBaseHandler) GetFireBases(c *gin.Context) {
	result, err := handler.FireBaseService.GetRefData()
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	c.JSON(200, result)
}
func (handler *FireBaseHandler) Truncate(c *gin.Context) {
	code, result := handler.FireBaseService.Truncate()
	c.JSON(code, result)
}
func (handler *FireBaseHandler) CreateEventFirebase(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(response.BadRequestErroMsg(map[string]interface{}{"result": "error", "message": err}))
		return
	}
	defer c.Request.Body.Close()
	var request model.EventBody
	err = json.Unmarshal(bodyBytes, &request)
	if err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	code, result := handler.FireBaseService.CreateEvent(request)

	c.JSON(code, result)
}
