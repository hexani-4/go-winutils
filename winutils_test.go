package winutils_test

import (
	"testing"

	"github.com/hexani-4/go-winutils"
)

const(
	ErrorMessageBox = false
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

func TestErrorMessageNotification(t *testing.T) {
	if ErrorMessageNotification{
		for _, message := range text_tc {
			err := winutils.ErrorMessageNotification(message)
			if err != nil { t.Fatalf(err.Error()) }
		}
    }
}