package events

import (
	"context"
	"fmt"
	logrus "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"testing"
	"time"
)

func TestEvents(t *testing.T) {

	ctx := context.Background()

	//SettingUp Simulator
	model := simulator.VPX()
	model.Machine = 5
	defer model.Remove()
	err := model.Create()
	if err != nil {
		logrus.Fatal(err)
	}
	s := model.Service.NewServer()

	c, _ := govmomi.NewClient(ctx, s.URL, true)

	var datacenters []mo.Datacenter

	m := view.NewManager(c.Client)

	cv, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create Datacenter container view")
	}
	err = cv.Retrieve(ctx, []string{"Datacenter"}, []string{"name", "overallStatus"}, &datacenters)
	assert.NoError(t, err)

	ca := NewCacheMock{}

	//https://pubs.vmware.com/vsphere-51/index.jsp?topic=%2Fcom.vmware.wssdk.apiref.doc%2Fvim.HistoryCollector.html
	ed, err := NewEventDispacher(c.Client, datacenters[0].Reference(), logrus.New(), &ca)
	assert.NoError(t, err)

	ed.CollectEvents("5")
	assert.Equal(t, 10, len(ed.Events), "We were expecting 10 events")
	ed.Cancel()

	ed2, err := NewEventDispacher(c.Client, datacenters[0].Reference(), logrus.New(), &ca)
	assert.NoError(t, err)
	ed2.CollectEvents("noParsable")
	assert.Equal(t, 10, len(ed2.Events), "We were expecting 10 events")
	ed2.Cancel()

	ed3, err := NewEventDispacher(c.Client, datacenters[0].Reference(), logrus.New(), &ca)
	assert.NoError(t, err)
	ed3.CollectEvents("3")
	assert.Equal(t, 10, len(ed.Events), "We were expecting 10 events")
	ed3.Cancel()

}

type NewCacheMock struct{}

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

func (*NewCacheMock) ReadTimestampCache() (time.Time, error) {
	return time.Now(), nil
}

func (*NewCacheMock) WriteTimestampCache(t time.Time) error {
	return nil
}
