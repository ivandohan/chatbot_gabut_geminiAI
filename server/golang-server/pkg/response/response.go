package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data       interface{}      `json:"data,omitempty"`
	Metadata   interface{}      `json:"metadata,omitempty"`
	Error      CustomErrorModel `json:"error"`
	StatusCode int              `json:"-"`
}

type CustomErrorModel struct {
	ErrorStatus  bool   `json:"status"`
	ErrorMessage string `json:"message"`
	ErrorCode    int    `json:"code"`
}

func (response *Response) RenderJSONResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if response.StatusCode == 0 {
		response.StatusCode = http.StatusOK
	}

	data, err := json.Marshal(response)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
	}

	w.WriteHeader(response.StatusCode)
	_, _ = w.Write(data)
}
