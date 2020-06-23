package events

import (
	"context"
	"fmt"
	"github.com/newrelic/nri-vsphere/internal/cache"
	logrus "github.com/sirupsen/Logrus"
	"github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/types"
	"strconv"
	"time"
)

type EventDispacher struct {
	cancelCtx context.CancelFunc
	collector *event.HistoryCollector
	ctx       *context.Context

	LastTimestamp *time.Time
	Events        []types.BaseEvent
	log           *logrus.Logger
	cacheResource string
	cachePath     string
}

const (
	pageSizeDefault = 200
)

func NewEventDispacher(client *vim25.Client, mo types.ManagedObjectReference, log *logrus.Logger, resourceName string, cachePath string) (*EventDispacher, error) {

	manager := event.NewManager(client)
	ctx, cancel := context.WithCancel(context.Background())

	now := time.Now()
	lastTimestamp, err := cache.ReadTimestampCache(cachePath, resourceName)
	lastTimestamp = sanitizeTimestamp(err, log, lastTimestamp, now)

	log.WithField("lastTimestamp", lastTimestamp.String()).Debug("Creating collector for events")
	collector, err := manager.CreateCollectorForEvents(ctx,
		types.EventFilterSpec{
			Time: &types.EventFilterSpecByTime{
				BeginTime: &lastTimestamp,
				EndTime:   &now,
			},
			Entity: &types.EventFilterSpecByEntity{
				Recursion: types.EventFilterSpecRecursionOptionAll,
				Entity:    mo,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error while creating historyCollector: %s ", err.Error())
	}

	ed := EventDispacher{
		cancelCtx:     cancel,
		LastTimestamp: &lastTimestamp,
		collector:     collector,
		ctx:           &ctx,
		Events:        []types.BaseEvent{},
		log:           log,
		cachePath:     cachePath,
		cacheResource: resourceName,
	}
	return &ed, nil
}

func sanitizeTimestamp(err error, log *logrus.Logger, lastTimestamp time.Time, now time.Time) time.Time {

	if err != nil {
		log.WithError(err).Error("Error reading cache, setting default timestamp to current time")
		lastTimestamp = now
		return lastTimestamp
	}

	//we are interested into the events logged since 1 second after the last one retrieved
	lastTimestamp = lastTimestamp.Add(time.Duration(1) * time.Second)

	limitTimestamp := now.Add(time.Duration(-1) * time.Hour)
	if lastTimestamp.Before(limitTimestamp) {
		//we try to avoid a deadlock where tue to a really old timestamp the integration try to fetch too many events timing out
		log.WithField("timestamp", lastTimestamp.String()).Warn("Timestamp is too old, setting lastTimestamp to 1 hour ago to fetch events")
		lastTimestamp = limitTimestamp
		return lastTimestamp
	}

	//we make sure the last timestamp is smaller or equal then now, otherwise the API call fails
	if lastTimestamp.After(now) {
		log.WithField("timestamp", lastTimestamp.String()).Warn("Timestamp after the time.Now(), setting lastTimestamp to time.Now()")
		lastTimestamp = now
		return lastTimestamp
	}

	//we are interested into the events logged since 1 second after the last one retrieved
	return lastTimestamp
}

func (ed *EventDispacher) Cancel() {
	ed.cancelCtx()
	ed.log.WithField("date", ed.LastTimestamp).Debug("saving in cache last read message ")

	for _, e := range ed.Events {
		if ed.LastTimestamp.Before(e.GetEvent().CreatedTime) {
			ed.LastTimestamp = &e.GetEvent().CreatedTime
		}
	}

	err := cache.WriteTimestampCache(ed.cachePath, ed.cacheResource, *ed.LastTimestamp)
	if err != nil {
		ed.log.WithError(err).Error("error while saving cache")
	}
}

func (ed *EventDispacher) CollectEvents(eventsPageSize string) {
	ed.log.WithField("timestamp", ed.LastTimestamp.String()).Debug("using as starting event")

	pageSize, err := strconv.Atoi(eventsPageSize)
	if err != nil {
		ed.log.WithError(err).Error("error while parsing EventsPageSize, using default value")
		pageSize = pageSizeDefault
	}

	for {
		eventsCollected, err := ed.collector.ReadNextEvents(*ed.ctx, int32(pageSize))
		if err != nil {
			ed.log.WithError(err).Error("error while fetching events")
			break
		}
		ed.Events = append(ed.Events, eventsCollected...)
		ed.log.WithField("number", len(eventsCollected)).Debug("readNextEventsExecuted")

		//There are no events left if: no events has been collected or if the number of events is smaller than the pagSize
		if len(eventsCollected) == 0 || len(eventsCollected) != pageSize {
			break
		}
	}
}
