package client

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCredentials(t *testing.T) {
	u := url.URL{}
	setCredentials(&u, "user", "password")
	assert.Equal(t, u.User, url.UserPassword("user", "password"))
}
