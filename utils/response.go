package utils

import (
	"fmt"
)

type Response struct {
	Head map[string]string `json:"Head"`
	Body interface{}       `json:"Body"`
}

type ArrayBodyStruct struct {
	CurrentPage  int         `json:"CurrentPage"`
	TotalPages   int         `json:"TotalPages"`
	PerPage      int         `json:"PerPage"`
	NextPage     int         `json:"NextPage"`
	PreviousPage int         `json:"PreviousPage"`
	Data         interface{} `json:"Data"`
}

type ArrayDataResponse struct {
	Head map[string]string `json:"Head"`
	Body interface{}       `json:"Body"`
}

var (
	SuccessResponse = Response{Head: map[string]string{"Code": "1000", "Msg": "Success."}}
	ArrayResponse   = ArrayDataResponse{Head: map[string]string{"Code": "1000", "Msg": "Success."}}
)

func BuildError(code string) Response {
	return Response{Head: map[string]string{"Code": code}}
}

func (errorResponse Response) Error() string {
	return fmt.Sprintf("Code: %s; Msg: %s", errorResponse.Head["Code"], errorResponse.Head["Msg"])
}

func (arrayResponse *ArrayDataResponse) Init(data interface{}, page, count, perPage int) {
	totalPage := count / perPage
	if (count % perPage) != 0 {
		totalPage += 1
	}

	nextPage := page + 1
	if nextPage > totalPage {
		nextPage = totalPage
	}
	previousPage := page - 1
	if previousPage < 1 {
		previousPage = 1
	}

	body := ArrayBodyStruct{}

	body.Data = data
	body.CurrentPage = page
	body.TotalPages = totalPage
	body.PerPage = perPage
	body.NextPage = nextPage
	body.PreviousPage = previousPage

	arrayResponse.Body = body
}
