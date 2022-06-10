package controller

import (
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"sync"
)

type BaseController struct {
	sync.RWMutex
	Code  int
	Msg   string
	Resp  interface{}
	Model int
	Lang  string
	Ip    string
	web.Controller
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//json 格式化json方式的返回
func (b *BaseController) json() {
	b.Data["json"] = &Response{
		Code: b.Code,
		Msg:  b.Msg,
		Data: b.Resp,
	}
	err := b.ServeJSON()
	if err != nil {
		panic(err)
	}
}

// Success 作为controller成功统一返回
//TODO 之后要兼容tpl,暂定为全局变量控制
func (b *BaseController) Success(o interface{}, msg ...string) {
	if len(msg) > 0 {
		b.Msg = msg[0]
	}
	b.Resp = o
	b.json()
	b.StopRun()
}

// Error 作为controller失败统一返回，错误码由上层进行国际化转换
func (b *BaseController) Error(code int, msg ...string) {
	b.Code = code
	b.ChkLang()
	if len(msg) > 0 {
		b.Msg = msg[0]
	} else {
		b.Msg = i18n.Tr(b.Lang, "error."+fmt.Sprint(code))
	}
	b.json()
	b.StopRun()
}

// InitLang 初始化区域语言包
func (b *BaseController) InitLang(langTypes ...string) error {
	//web.AppConfig.DefaultStrings("LangTypes")
	for _, lt := range langTypes {
		if lt == "" {
			continue
		}
		err := i18n.SetMessage(lt, "conf/locale_"+lt+".ini")
		if err != nil {
			return err
		}
	}
	return nil
}

// ChkLang 语言判断
// conf路径下需存在对应的语言包locale_en-US.ini或locale_zh-CN.ini等
// 需要先行调用InitLang
func (b *BaseController) ChkLang() {
	cl := b.Ctx.Request.Header.Get("Accept-Language")
	if len(cl) > 4 && i18n.IsExist(cl) {
		b.Lang = cl
	} else {
		b.Lang = "zh-CN"
	}
}
