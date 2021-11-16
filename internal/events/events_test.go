package events

import (
	"context"
	"fmt"
	"testing"
	"time"

	logrus "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vim25"
)

func TestEvents(t *testing.T) {
	model := simulator.VPX()
	simulator.Test(func(ctx context.Context, vc *vim25.Client) {
		ca := NewCacheMock{
			TimestampCache: time.Now().Add(-15 * time.Second),
		}

		// A virtual machine contains 6 events on the simulator.
		ref := simulator.Map.Any("VirtualMachine").Reference()

		// https://pubs.vmware.com/vsphere-51/index.jsp?topic=%2Fcom.vmware.wssdk.apiref.doc%2Fvim.HistoryCollector.html
		ed, err := NewEventDispacher(vc, ref, logrus.New(), &ca)
		assert.NoError(t, err)

		ed.CollectEvents("5")
		assert.Equal(t, 6, len(ed.Events), "We were expecting 6 events")
		ed.Cancel()

		ed.CollectEvents("noParsable")
		assert.Equal(t, 6, len(ed.Events), "We were expecting 6 events")
		ed.Cancel()
	}, model)
}

type NewCacheMock struct {
	TimestampCache time.Time
}

func TestSanitizeTimestamp(t *testing.T) {
	now := time.Now()
	log := logrus.New()
	last := time.Now().Add(time.Duration(-30) * time.Minute)
	lastTooOld := time.Now().Add(time.Duration(-30) * time.Hour)
	lastTooNew := time.Now().Add(time.Duration(30) * time.Minute)

	err := fmt.Errorf("random Error")

	s := sanitizeTimestamp(err, log, time.Time{}, now)
	assert.Equal(t, now, s)

	s = sanitizeTimestamp(nil, log, lastTooNew, now)
	assert.Equal(t, now, s)

	s = sanitizeTimestamp(nil, log, lastTooOld, now)
	assert.Equal(t, now.Add(time.Duration(-1)*time.Hour), s)

	s = sanitizeTimestamp(nil, log, last, now)
	assert.Equal(t, last.Add(time.Duration(1)*time.Second), s)
}

func (c *NewCacheMock) ReadTimestampCache() (time.Time, error) {
	return c.TimestampCache, nil
}

func (c *NewCacheMock) WriteTimestampCache(t time.Time) error {
	return nil
}
