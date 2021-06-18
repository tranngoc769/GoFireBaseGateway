package service

import (
	"go-firebase-gateway/common/auth"
	"go-firebase-gateway/common/log"
	"go-firebase-gateway/common/response"
)

type UserService struct {
}

func NewUserService() UserService {
	return UserService{}
}
func (service *UserService) GenerateTokenByApiKey(apiKey string, isRefresh bool) (int, interface{}) {
	log.Debug("UserService", "GenerateTokenByApiKey", apiKey)
	clientAuth := auth.AuthClient{
		ClientID:     apiKey,
		ClientSecret: apiKey,
	}
	token, err := auth.ClientCredential(clientAuth, isRefresh)
	if err != nil {
		log.Error("UserService", "GenerateTokenByApiKey", err.Error())
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.NewOKResponse(token)
}
