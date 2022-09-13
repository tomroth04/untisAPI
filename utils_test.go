package untisAPI

import (
	"testing"
)

func TestBase64(t *testing.T) {
	res := ToBase64("lmrl")
	if res != "bG1ybA==" {
		t.Error("Error base64 doesn't match")
	}
}

func TestGetDateUntisFormat(t *testing.T) {
	// TODO: implement test

}
