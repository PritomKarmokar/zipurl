package response

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

type CodeObject struct {
	httpStatus         int
	stateCode          string
	stateMessage       map[string]string
	stateMessageParams map[string]interface{}
}

type codeResponse struct {
	StateCode    string      `json:"code"`
	StateMessage string      `json:"message"`
	Lang         string      `json:"lang"`
	Data         interface{} `json:"data,omitempty"`
}

func NewCodeObject(httpStatus int, stateCode string, stateMessage map[string]string, stateMessageParams map[string]interface{}) *CodeObject {
	if stateMessage == nil {
		stateMessage = map[string]string{"en": "No message defined"}
	}
	if stateMessageParams == nil {
		stateMessageParams = map[string]interface{}{}
	}
	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}
	if stateCode == "" {
		stateCode = http.StatusText(httpStatus)
	}
	return &CodeObject{
		httpStatus:         httpStatus,
		stateCode:          stateCode,
		stateMessage:       stateMessage,
		stateMessageParams: stateMessageParams,
	}
}

// Getters
func (c *CodeObject) HttpStatus() int {
	return c.httpStatus
}

func (c *CodeObject) StateCode() string {
	return c.stateCode
}

func (c *CodeObject) StateMessage() map[string]string {
	return c.stateMessage
}

func (c *CodeObject) StateMessageParams() map[string]interface{} {
	return c.stateMessageParams
}

// Setters
func (c *CodeObject) SetHttpStatus(status int) {
	c.httpStatus = status
}

func (c *CodeObject) SetStateCode(code string) {
	c.stateCode = code
}

func (c *CodeObject) SetStateMessage(message map[string]string) {
	c.stateMessage = message
}

func (c *CodeObject) SetStateMessageParams(params map[string]interface{}) {
	c.stateMessageParams = params
}

func (c *CodeObject) formatResponse(data interface{}, lang string) codeResponse {
	return codeResponse{
		StateCode:    c.stateCode,
		StateMessage: c.stateMessage[lang],
		Lang:         lang,
		Data:         data,
	}
}

func (c *CodeObject) getLanguage(req *http.Request) string {
	lang := req.Header.Get("Accept-Language")
	lang = strings.ToLower(lang)
	if lang == "bn" || lang == "ban" || lang == "bang" || lang == "bangla" || lang == "bengali" {
		return "bn"
	}
	return "en"
}

func (c *CodeObject) ReturnResponse(ctx *echo.Context, data interface{}) error {
	lang := c.getLanguage(ctx.Request())
	return ctx.JSON(c.httpStatus, c.formatResponse(data, lang))
}

func (c *CodeObject) IsHttpErrorStatus() bool {
	return c.httpStatus >= http.StatusInternalServerError && c.httpStatus <= 599
}

func (c *CodeObject) IsHttpSuccessStatus() bool {
	return c.httpStatus >= http.StatusOK && c.httpStatus <= 299
}
