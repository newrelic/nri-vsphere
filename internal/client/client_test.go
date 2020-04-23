package client

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestSetCredentials(t *testing.T) {
	u := url.URL{}
	setCredentials(&u, "user", "password")
	assert.Equal(t, u.User, url.UserPassword("user", "password"))
}
