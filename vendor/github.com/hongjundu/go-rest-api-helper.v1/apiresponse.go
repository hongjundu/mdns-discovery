package apihelper

const (
	ResponseStatusOK    = "ok"
	ResponseStatusError = "error"
)

func NewOKResponse(data interface{}) interface{} {
	return &apiResponse{Status: ResponseStatusOK, Code: 200, Data: data}
}

func NewListResponse(count int, list interface{}) interface{} {
	return NewOKResponse(NewListResponseData(count, list))
}

func NewPageableListResponse(total, count, page, limit int, list interface{}) interface{} {
	return NewOKResponse(NewPageableListResponseData(total, count, page, limit, list))
}

func NewListResponseData(count int, list interface{}) interface{} {
	return &responseListData{List: list, Count: count}
}

func NewPageableListResponseData(total, count, page, limit int, list interface{}) interface{} {
	return &pageableResponseListData{responseListData: responseListData{List: list, Count: count}, Total: total, Page: page, Limit: limit}
}

func NewErrorResponse(err error) interface{} {
	if err == nil {
		return NewOKResponse(nil)
	} else {
		if apiErr, ok := err.(ApiError); ok {
			return &apiResponse{Status: ResponseStatusError, Code: apiErr.Code(), Msg: apiErr.Error()}
		} else {
			return &apiResponse{Status: ResponseStatusError, Code: 500, Msg: err.Error()}
		}
	}
}

// ok response format:
// {"status":"ok","data":{ ... }}
// error response format:
// {"status":"error","code":400,"msg":"error message description"}

type apiResponse struct {
	Status string      `json:"status"`
	Code   int         `json:"code,omitempty"`
	Msg    string      `json:"msg,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

type responseListData struct {
	Count int         `json:"count"`
	List  interface{} `json:"list,omitempty"`
}

type pageableResponseListData struct {
	responseListData
	Total int `json:"total"`
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}
