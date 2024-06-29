package services

import (
	"fmt"
	"runtime"
	"strings"
)

// 处理定时任务中的panic
func HandlePanic(r interface{}, prefixString string, errType string) {
	var msg string
	threshold := 2
	for skip, stackNum := 0, 1; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			msg = fmt.Sprintf("Unable to retrieve panic information.")
			break
		}

		funcName := runtime.FuncForPC(pc).Name()
		if !strings.Contains(funcName, "runtime.") && !strings.Contains(file, "handle_defender.go") {
			msg = msg + fmt.Sprintf("\nFunction: %s, File: %s, Line: %d, Panic: %v", funcName, file, line, r)

			stackNum = stackNum + 1
			if stackNum > threshold {
				break
			}
		}
	}
	//msgText := fmt.Sprintf("%s, Error: %s", prefixString, msg)

	//my_log.LogError(msgText)
	//SetDBErr(msgText, errType)
	//go tg_bot.SendAdmin(msgText)
}

func HandleError(err error, prefixString string, errType string) {
	//msgText := fmt.Sprintf("%s, Error: %v", prefixString, err)
	//my_log.LogError(msgText)
	//SetDBErr(msgText, errType)
	//go tg_bot.SendAdmin(msgText)
}
