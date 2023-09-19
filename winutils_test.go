package winutils_test

import (
	"testing"

	"github.com/hexani-4/go-winutils"
)

func TestErrorMessageBox(t *testing.T) {
	err := winutils.ErrorMessageBox("title", "message")
	if err != nil { t.Fatalf(err.Error()) }
}