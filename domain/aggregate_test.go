package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestNewBaseAggregate(t *testing.T) {
	a := NewBaseAggregate("agg-1")
	assert.Equal(t, "agg-1", a.ID)
	assert.Equal(t, int64(0), a.Version)
	assert.Empty(t, a.Events)
}

func TestBaseAggregate_GetID(t *testing.T) {
	a := NewBaseAggregate("my-id")
	assert.Equal(t, "my-id", a.GetID())
}

func TestBaseAggregate_GetVersion(t *testing.T) {
	a := NewBaseAggregate("id")
	assert.Equal(t, int64(0), a.GetVersion())
	a.Version = 5
	assert.Equal(t, int64(5), a.GetVersion())
}

func TestBaseAggregate_SetID(t *testing.T) {
	a := NewBaseAggregate("old")
	a.SetID("new")
	assert.Equal(t, "new", a.GetID())
}

func TestBaseAggregate_SetVersion(t *testing.T) {
	a := NewBaseAggregate("id")
	a.SetVersion(10)
	assert.Equal(t, int64(10), a.GetVersion())
}

func TestBaseAggregate_GetUncommittedEvents(t *testing.T) {
	a := NewBaseAggregate("id")
	assert.Empty(t, a.GetUncommittedEvents())
	e := &emptypb.Empty{}
	a.AddEvent(e)
	events := a.GetUncommittedEvents()
	assert.Len(t, events, 1)
	assert.Equal(t, e, events[0])
}

func TestBaseAggregate_MarkEventsAsCommitted(t *testing.T) {
	a := NewBaseAggregate("id")
	a.AddEvent(&emptypb.Empty{})
	assert.Len(t, a.Events, 1)
	a.MarkEventsAsCommitted()
	assert.Empty(t, a.Events)
}

func TestBaseAggregate_AddEvent(t *testing.T) {
	a := NewBaseAggregate("id")
	e1, e2 := &emptypb.Empty{}, &emptypb.Empty{}
	a.AddEvent(e1)
	assert.Len(t, a.Events, 1)
	assert.Equal(t, int64(1), a.Version)
	a.AddEvent(e2)
	assert.Len(t, a.Events, 2)
	assert.Equal(t, int64(2), a.Version)
}

func TestBaseAggregate_LoadFromHistory(t *testing.T) {
	a := NewBaseAggregate("id")
	events := []proto.Message{&emptypb.Empty{}, &emptypb.Empty{}}
	err := a.LoadFromHistory(events)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), a.Version)
}

func TestBaseAggregate_LoadFromHistory_Empty(t *testing.T) {
	a := NewBaseAggregate("id")
	err := a.LoadFromHistory(nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), a.Version)
	err = a.LoadFromHistory([]proto.Message{})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), a.Version)
}

func TestBaseAggregate_LoadFromHistory_SingleEvent(t *testing.T) {
	a := NewBaseAggregate("id")
	err := a.LoadFromHistory([]proto.Message{&emptypb.Empty{}})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), a.Version)
}
