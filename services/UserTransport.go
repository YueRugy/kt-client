package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func GetUserInfoRequest(_ context.Context, request *http.Request, r interface{}) error {
	ur := r.(UserRequest)
	request.URL.Path += "/user/" + strconv.Itoa(ur.Uid)
	return nil
}

func GetUserInfoResponse(_ context.Context, response *http.Response) (interface{}, error) {
	if response.StatusCode > 400 {
		return nil, errors.New("no data")
	}
	var ur UserResponse
	err := json.NewDecoder(response.Body).Decode(&ur)
	if err != nil {
		return nil, err
	}
	return ur, nil
}
