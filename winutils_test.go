package winutils_test

import (
	"fmt"
	"testing"

	"github.com/hexani-4/go-winutils"
)

const(
	ErrorMessageBox = false
	LoadErrIcon = false
	ErrorMessageNotification = true
)



var text_tc = map[string]string {
	"title": "message",
	"" : "",
	"\\0" : "\\0",
	"\x00" : "\x00",
}


func TestErrorMessageBox(t *testing.T) {
	if ErrorMessageBox {
		for title, message := range text_tc {
			err := winutils.ErrorMessageBox(title, message)
			if err != nil { t.Fatalf(err.Error()) }
		}
	}
}

func TestLoadErrIcon(* testing.T) {
	if LoadErrIcon {
		hIcon, err := winutils.LoadErrIcon()
		fmt.Println(hIcon, err)
	}
}

func TestErrorMessageNotification(t *testing.T) {
	if ErrorMessageNotification{
	
		i := uint32(2000)
		for _, message := range text_tc {
			err := winutils.ErrorMessageNotification(message, i)
			if err != nil { t.Fatalf(err.Error()) }
			i -= 250
		}
    }
}