package cache

import (
	"testing"
	"time"

	"github.com/newrelic/infra-integrations-sdk/v3/persist"
	"github.com/stretchr/testify/assert"
)

func Test_Cache_SavesCorrectTimestamp(t *testing.T) {

	datacenter := "my-datacenter"
	store := persist.NewInMemoryStore()

	//given
	expected := time.Now()
	c := NewCache(datacenter, store)
	assert.NotNil(t, c)

	//when
	err := c.WriteTimestampCache(expected)
	assert.NoError(t, err)

	//then
	actual, err := c.ReadTimestampCache()
	assert.NoError(t, err)

	assert.True(t, expected.Equal(actual))
	assert.Equal(t, expected.UnixNano(), actual.UnixNano())
}
