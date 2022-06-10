package controller

const (
	_                  = -iota
	ServerError        //服务错误
	PageNotExist       //页面不存在
	MethodNotAllow     //不支持该方式
	ServiceUnavailable //服务未待命
)

type ErrorController struct {
	BaseController
}

func (c *ErrorController) Error404() {
	c.Error(PageNotExist)
}

func (c *ErrorController) Error500() {
	c.Error(ServerError)
}

func (c *ErrorController) Error501() {
	c.Error(MethodNotAllow)
}

func (c *ErrorController) Error502() {
	c.Error(ServiceUnavailable)
}
