package util

import (
	"encoding/json"
	IRedis "go-firebase-gateway/internal/redis"
)

var CallBackListID string
var CallBackRequestHash string
var HookCallStatusListID string
var HookCallStatusHash string

func InterfaceToJson(mapItf interface{}) (string, error) {
	jsonStr, err := json.Marshal(mapItf)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}
func SetRequestIdToList(requestId string) error {
	ctx := IRedis.Redis.GetClient().Context()
	_ = IRedis.Redis.GetClient().Do(ctx, "LREM", CallBackListID, 0, requestId)
	_ = IRedis.Redis.GetClient().Do(ctx, "rpush", CallBackListID, requestId)
	return nil
}

func SetCustomerDataHash(requestId string, requestData interface{}) error {
	dataString, err := InterfaceToJson(requestData)
	if err != nil {
		return err
	}
	data := []interface{}{requestId, dataString}
	_, err = IRedis.Redis.HSet(CallBackRequestHash, data)
	return err
}
func GetListRequestId(amount int) (interface{}, error) {
	ctx := IRedis.Redis.GetClient().Context()
	data, err := IRedis.Redis.GetClient().Do(ctx, "LRANGE", CallBackListID, 0, amount-1).Result()
	return data, err
}
func PopFirstRequestIDList() (interface{}, error) {
	ctx := IRedis.Redis.GetClient().Context()
	data, err := IRedis.Redis.GetClient().Do(ctx, "LPOP", CallBackListID).Result()
	return data, err
}

func LeftPushList(requestId string) error {
	ctx := IRedis.Redis.GetClient().Context()
	_ = IRedis.Redis.GetClient().Do(ctx, "lpush", CallBackListID, requestId)
	return nil
}
func GetDataByRequestId(requestId string) (interface{}, error) {
	data, err := IRedis.Redis.HGet(CallBackRequestHash, requestId)
	return data, err
}
func DeleteRequestIdFromHash(requestId string) error {
	err := IRedis.Redis.HDel(CallBackRequestHash, requestId)
	return err
}

// Set CallBackCallStatus DAta to List
func SetLeadIdToList(leadId int) error {
	ctx := IRedis.Redis.GetClient().Context()
	_ = IRedis.Redis.GetClient().Do(ctx, "LREM", HookCallStatusListID, 0, leadId)
	_ = IRedis.Redis.GetClient().Do(ctx, "rpush", HookCallStatusListID, leadId)
	return nil
}

// Set CallBackCallStatus DAta
func SetCallStatusData(leadId int, callBackData interface{}) error {
	dataString, err := InterfaceToJson(callBackData)
	if err != nil {
		return err
	}
	data := []interface{}{leadId, dataString}
	_, err = IRedis.Redis.HSet(HookCallStatusHash, data)
	return err
}
func GetListLeadId(amount int) (interface{}, error) {
	ctx := IRedis.Redis.GetClient().Context()
	data, err := IRedis.Redis.GetClient().Do(ctx, "LRANGE", HookCallStatusListID, 0, amount-1).Result()
	return data, err
}
func PopFirstLeadIdList() (interface{}, error) {
	ctx := IRedis.Redis.GetClient().Context()
	data, err := IRedis.Redis.GetClient().Do(ctx, "LPOP", HookCallStatusListID).Result()
	return data, err
}

func LeftPushLeadList(requestId string) error {
	ctx := IRedis.Redis.GetClient().Context()
	_ = IRedis.Redis.GetClient().Do(ctx, "lpush", HookCallStatusListID, requestId)
	return nil
}
func GetDataByLeadId(leadId string) (interface{}, error) {
	data, err := IRedis.Redis.HGet(HookCallStatusHash, leadId)
	return data, err
}
func DeleteLeadFromHash(leadId string) error {
	err := IRedis.Redis.HDel(HookCallStatusHash, leadId)
	return err
}
func CheckListExist(listName string) bool {
	ctx := IRedis.Redis.GetClient().Context()
	data, _ := IRedis.Redis.GetClient().Do(ctx, "LRANGE", listName, 0, 0).Result()
	return len(data.([]interface{})) > 0
}
