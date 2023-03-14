package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBase64(t *testing.T) {
	assert.Equal(t, "bG1ybA==", ToBase64("lmrl"))
}

func TestGetDateUntisFormat(t *testing.T) {
	assert.Equal(t, "20221202", GetDateUntisFormat(time.Date(2022, 12, 2, 0, 0, 0, 0, time.Now().Location())))
}
