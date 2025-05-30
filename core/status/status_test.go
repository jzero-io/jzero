package status

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	GetUserListError = Code(28001)
)

func TestError(t *testing.T) {
	err := Error(GetUserListError)
	status := FromError(err)
	assert.Equal(t, GetUserListError, status.Code())
	assert.Equal(t, "get user list error", status.Message())
}

func TestUnknownError(t *testing.T) {
	err := Error(28000)

	status := FromError(err)
	assert.Equal(t, http.StatusInternalServerError, int(status.Code()))
}

func TestWrap(t *testing.T) {
	err := Wrap(GetUserListError, errors.New("connect to db error"))
	status := FromError(err)
	assert.Equal(t, GetUserListError, status.Code())
	assert.Equal(t, "get user list error: connect to db error", status.Error())
	assert.Equal(t, "connect to db error", status.Unwrap().Error())
}

func init() {
	RegisterWithMessage(GetUserListError, "get user list error")
}
