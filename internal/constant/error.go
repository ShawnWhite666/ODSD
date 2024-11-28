package constant

import (
	"errors"
	"fmt"
)

var (
	ERR_HANDLE_INPUT = errors.New("handle input error")
)

type ErrCode int // 错误码

const (
	Success           ErrCode = 0
	ErrInternalServer ErrCode = 500
	ErrInputInvalid   ErrCode = 8020
	ErrShouldBind     ErrCode = 8021
	ErrJsonMarshal    ErrCode = 8022
	ErrJwtParse       ErrCode = 8023

	ErrLogin            ErrCode = 10000
	ErrIPLimitInvalid   ErrCode = 10001
	ErrUserLimitInvalid ErrCode = 10002
	ErrBlackedIP        ErrCode = 10003
	ErrBlackedUser      ErrCode = 10004
	ErrPrizeNotEnough   ErrCode = 10005
	ErrNotWon           ErrCode = 100010
)

var errMsgDic = map[ErrCode]string{
	Success:             "ok",
	ErrInternalServer:   "internal server error",
	ErrInputInvalid:     "input invalid",
	ErrShouldBind:       "should bind failed",
	ErrJwtParse:         "json marshal failed",
	ErrLogin:            "login fail",
	ErrIPLimitInvalid:   "ip day num limited",
	ErrUserLimitInvalid: "user day num limited",
	ErrBlackedIP:        "blacked ip",
	ErrBlackedUser:      "blacked user",
	ErrPrizeNotEnough:   "prize not enough",
	ErrNotWon:           "not won,please try again!",
}

// GetErrMsg 获取错误描述
func GetErrMsg(code ErrCode) string {
	if msg, ok := errMsgDic[code]; ok {
		return msg
	}
	return fmt.Sprintf("unknown error code %d", code)
}
