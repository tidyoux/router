package handler

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

var (
	handlers = map[string]interface{}{}
)

func H(path string, h interface{}) {
	if len(path) == 0 {
		return
	}

	if h == nil {
		return
	}

	handlers[path] = h
}

type Response struct {
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Data: data,
	}
}

func NewFailResponse(err string) *Response {
	return &Response{
		Error: err,
	}
}

func Init(group *gin.RouterGroup) error {
	for path, h := range handlers {
		fv := reflect.ValueOf(h)
		if fv.Kind() != reflect.Func {
			return fmt.Errorf("invalid %s handler type: %s", path, fv.Kind())
		}

		ft := fv.Type()

		// Check input.
		if ft.NumIn() != 1 {
			return fmt.Errorf("invalid %s handler input count: %d, should be %d",
				path, ft.NumIn(), 1)
		}

		if ft.In(0).Kind() != reflect.Ptr {
			return fmt.Errorf("invalid %s handler input type: %s, should be %s",
				path, ft.In(0).Kind(), reflect.Ptr)
		}

		// Check output.
		if ft.NumOut() != 2 {
			return fmt.Errorf("invalid %s handler output count: %d, should be %d",
				path, ft.NumOut(), 2)
		}

		if !ft.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			return fmt.Errorf("invalid %s handler last output type, should be error",
				path)
		}

		// Register handler.
		group.POST(path, func(c *gin.Context) {
			data, err := handle(c, fv, ft)
			if err != nil {
				c.JSON(http.StatusOK, NewFailResponse(err.Error()))
				return
			}

			c.JSON(http.StatusOK, NewSuccessResponse(data))
		})
	}

	return nil
}

func handle(c *gin.Context, fv reflect.Value, ft reflect.Type) (interface{}, error) {
	reqT := ft.In(0).Elem()
	reqV := reflect.New(reqT)
	req := reqV.Interface()
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, fmt.Errorf("bind args failed, %v", err)
	}

	rets := fv.Call([]reflect.Value{reqV})
	if !rets[1].IsNil() {
		return nil, rets[1].Interface().(error)
	}

	return rets[0].Interface(), nil
}
