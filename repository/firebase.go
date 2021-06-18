package repository

import (
	"context"
	"go-firebase-gateway/common/model"

	"firebase.google.com/go/db"
)

var FireBaseContext context.Context
var EventRef *db.Ref

type FireBaseRepository struct {
}

func NewFireBaseRepository() FireBaseRepository {
	repo := FireBaseRepository{}
	return repo
}

var FireBaseRepo FireBaseRepository

func (repo *FireBaseRepository) GetRefData() (interface{}, error) {
	var data interface{}
	if err := EventRef.Get(FireBaseContext, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *FireBaseRepository) CreateEvent(event model.EventBody) (interface{}, error) {
	data, err := EventRef.Push(FireBaseContext, event)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *FireBaseRepository) Truncate() (interface{}, error) {
	err := EventRef.Delete(FireBaseContext)
	if err != nil {
		return false, err
	}
	return true, nil
}
