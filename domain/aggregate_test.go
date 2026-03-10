package domain

import (
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDomainAggregate(t *testing.T) {
	specs.Describe(t, "domain aggregate", func(s *specs.Spec) {
		s.When("NewBaseAggregate", func(s *specs.Spec) {
			s.It("sets ID and zero version and empty events", func(ctx *specs.Context) {
				a := NewBaseAggregate("agg-1")
				assert.Equal(ctx.T, "agg-1", a.ID)
				assert.Equal(ctx.T, int64(0), a.Version)
				assert.Empty(ctx.T, a.Events)
			})
		})

		s.When("BaseAggregate_GetID", func(s *specs.Spec) {
			s.It("returns aggregate ID", func(ctx *specs.Context) {
				a := NewBaseAggregate("my-id")
				assert.Equal(ctx.T, "my-id", a.GetID())
			})
		})

		s.When("BaseAggregate_GetVersion", func(s *specs.Spec) {
			s.It("returns current version", func(ctx *specs.Context) {
				a := NewBaseAggregate("id")
				assert.Equal(ctx.T, int64(0), a.GetVersion())
				a.Version = 5
				assert.Equal(ctx.T, int64(5), a.GetVersion())
			})
		})

		s.When("BaseAggregate_SetID", func(s *specs.Spec) {
			s.It("updates ID", func(ctx *specs.Context) {
				a := NewBaseAggregate("old")
				a.SetID("new")
				assert.Equal(ctx.T, "new", a.GetID())
			})
		})

		s.When("BaseAggregate_SetVersion", func(s *specs.Spec) {
			s.It("updates version", func(ctx *specs.Context) {
				a := NewBaseAggregate("id")
				a.SetVersion(10)
				assert.Equal(ctx.T, int64(10), a.GetVersion())
			})
		})

		s.When("BaseAggregate_GetUncommittedEvents", func(s *specs.Spec) {
			s.It("returns empty then events after AddEvent", func(ctx *specs.Context) {
				a := NewBaseAggregate("id")
				assert.Empty(ctx.T, a.GetUncommittedEvents())
				e := &emptypb.Empty{}
				a.AddEvent(e)
				events := a.GetUncommittedEvents()
				assert.Len(ctx.T, events, 1)
				assert.Equal(ctx.T, e, events[0])
			})
		})

		s.When("BaseAggregate_MarkEventsAsCommitted", func(s *specs.Spec) {
			s.It("clears events", func(ctx *specs.Context) {
				a := NewBaseAggregate("id")
				a.AddEvent(&emptypb.Empty{})
				assert.Len(ctx.T, a.Events, 1)
				a.MarkEventsAsCommitted()
				assert.Empty(ctx.T, a.Events)
			})
		})

		s.When("BaseAggregate_AddEvent", func(s *specs.Spec) {
			s.It("appends event and increments version", func(ctx *specs.Context) {
				a := NewBaseAggregate("id")
				e1, e2 := &emptypb.Empty{}, &emptypb.Empty{}
				a.AddEvent(e1)
				assert.Len(ctx.T, a.Events, 1)
				assert.Equal(ctx.T, int64(1), a.Version)
				a.AddEvent(e2)
				assert.Len(ctx.T, a.Events, 2)
				assert.Equal(ctx.T, int64(2), a.Version)
			})
		})

		s.When("BaseAggregate_LoadFromHistory", func(s *specs.Spec) {
			s.It("sets version from events", func(ctx *specs.Context) {
				a := NewBaseAggregate("id")
				events := []proto.Message{&emptypb.Empty{}, &emptypb.Empty{}}
				err := a.LoadFromHistory(events)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, int64(2), a.Version)
			})
			s.It("empty or nil leaves version zero", func(ctx *specs.Context) {
				a := NewBaseAggregate("id")
				err := a.LoadFromHistory(nil)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, int64(0), a.Version)
				err = a.LoadFromHistory([]proto.Message{})
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, int64(0), a.Version)
			})
			s.It("single event sets version 1", func(ctx *specs.Context) {
				a := NewBaseAggregate("id")
				err := a.LoadFromHistory([]proto.Message{&emptypb.Empty{}})
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, int64(1), a.Version)
			})
			s.It("many events set version to count", func(ctx *specs.Context) {
				a := NewBaseAggregate("id")
				events := []proto.Message{&emptypb.Empty{}, &emptypb.Empty{}, &emptypb.Empty{}}
				err := a.LoadFromHistory(events)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, int64(3), a.Version)
			})
			s.It("five events increments version five times", func(ctx *specs.Context) {
				a := NewBaseAggregate("id")
				events := []proto.Message{
					&emptypb.Empty{}, &emptypb.Empty{}, &emptypb.Empty{},
					&emptypb.Empty{}, &emptypb.Empty{},
				}
				err := a.LoadFromHistory(events)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, int64(5), a.Version)
			})
		})
	})
}
