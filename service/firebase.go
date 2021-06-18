package service

import (
	"go-firebase-gateway/common/log"
	"go-firebase-gateway/common/model"
	"go-firebase-gateway/common/response"
	"go-firebase-gateway/repository"
)

type FireBaseService struct {
}

func NewFireBaseService() FireBaseService {
	return FireBaseService{}
}

func (service *FireBaseService) GetRefData() (int, interface{}) {
	FireBaseRes, err := repository.FireBaseRepo.GetRefData()
	if err != nil {
		log.Error("FireBaseService", "GetRefData", err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.NewOKResponse(FireBaseRes)
}
func (service *FireBaseService) CreateEvent(event model.EventBody) (int, interface{}) {
	FireBaseRes, err := repository.FireBaseRepo.CreateEvent(event)
	if err != nil {
		log.Error("FireBaseService", "CreateEvent", err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.NewOKResponse(FireBaseRes)
}
