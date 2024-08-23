package base

import (
	"encoding/json"
	"net/http"
)

type ListData struct {
	Rows  interface{} `json:"rows"`
	Total *int64      `json:"total"`
}

type ResponseData struct {
	Code int64       `json:"code"`           //相应状态码
	Msg  string      `json:"msg"`            //提示信息
	Data interface{} `json:"data,omitempty"` //数据
}

type Result struct {
	http.ResponseWriter
}

func R(w http.ResponseWriter) *Result {
	return &Result{w}
}

func (that *Result) Ok(data interface{}) {
	json.NewEncoder(that).Encode(ResponseData{Code: 200, Msg: "操作成功", Data: data})
}

func (that *Result) OkList(rows interface{}, total *int64) {
	json.NewEncoder(that).Encode(ResponseData{Code: 200, Msg: "操作成功", Data: ListData{Rows: rows, Total: total}})
}

func (that *Result) OkMsg(data interface{}, msg string) {
	json.NewEncoder(that).Encode(ResponseData{Code: 200, Msg: msg, Data: data})
}

func (that *Result) Fail() {
	json.NewEncoder(that).Encode(ResponseData{Code: 500, Msg: "操作失败"})
}

func (that *Result) FailMsg(msg string) {
	json.NewEncoder(that).Encode(ResponseData{Code: 500, Msg: msg})
}
