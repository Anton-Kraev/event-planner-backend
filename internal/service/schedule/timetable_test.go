package schedule

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/Anton-Kraev/event-planner-backend/internal/domain/schedule"
	"github.com/Anton-Kraev/event-planner-backend/internal/domain/schedule/timetable"
	"github.com/Anton-Kraev/event-planner-backend/internal/service/schedule/mocks"
)

func TestService_GetTimetableSchedule(t *testing.T) {
	t.Parallel()

	type fields struct {
		client *mocks.MocktimetableClient
		cache  *mocks.MocktimetableCache
	}

	var (
		someCtx             = context.Background()
		someErr             = errors.New("")
		someTimetableEvents = []timetable.Event{{}, {}, {}}
		someCalendarOwner   = timetable.CalendarOwner{
			Class: timetable.Classroom,
			Name:  "classroom",
		}
		someCalendar = schedule.Calendar{
			Source: schedule.Timetable,
			Owner:  someCalendarOwner.String(),
			Events: []schedule.Event{{}, {}, {}},
		}
	)

	tests := []struct {
		name    string
		fields  fields
		mockFn  func(f fields)
		want    schedule.Calendar
		wantErr bool
	}{
		{
			name: "events_from_timetable",
			mockFn: func(f fields) {
				f.cache.EXPECT().
					GetEvents(gomock.Any(), someCalendarOwner).
					Return(nil, timetable.ErrNotCachedYet)
				f.client.EXPECT().
					GetClassroomEvents(gomock.Any(), someCalendarOwner.Name).
					Return(someTimetableEvents, nil)
				f.cache.EXPECT().
					SetEvents(gomock.Any(), someCalendarOwner, someTimetableEvents).
					Return(nil)
			},
			want:    someCalendar,
			wantErr: false,
		},
		{
			name: "events_from_cache",
			mockFn: func(f fields) {
				f.cache.EXPECT().
					GetEvents(gomock.Any(), someCalendarOwner).
					Return(someTimetableEvents, nil)
				f.cache.EXPECT().
					SetEvents(gomock.Any(), someCalendarOwner, someTimetableEvents).
					Return(nil)
			},
			want:    someCalendar,
			wantErr: false,
		},
		{
			name: "err_cache_get",
			mockFn: func(f fields) {
				f.cache.EXPECT().
					GetEvents(gomock.Any(), someCalendarOwner).
					Return(nil, someErr)
			},
			want:    schedule.Calendar{},
			wantErr: true,
		},
		{
			name: "err_timetable_get",
			mockFn: func(f fields) {
				f.cache.EXPECT().
					GetEvents(gomock.Any(), someCalendarOwner).
					Return(nil, timetable.ErrNotCachedYet)
				f.client.EXPECT().
					GetClassroomEvents(gomock.Any(), someCalendarOwner.Name).
					Return(nil, someErr)
			},
			want:    schedule.Calendar{},
			wantErr: true,
		},
		{
			name: "err_cache_set",
			mockFn: func(f fields) {
				f.cache.EXPECT().
					GetEvents(gomock.Any(), someCalendarOwner).
					Return(someTimetableEvents, nil)
				f.cache.EXPECT().
					SetEvents(gomock.Any(), someCalendarOwner, someTimetableEvents).
					Return(someErr)
			},
			want:    schedule.Calendar{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			tt.fields = fields{
				client: mocks.NewMocktimetableClient(ctrl),
				cache:  mocks.NewMocktimetableCache(ctrl),
			}
			tt.mockFn(tt.fields)

			res, err := NewService(tt.fields.client, tt.fields.cache).
				GetTimetableSchedule(someCtx, someCalendarOwner)

			require.True(t, (err != nil) == tt.wantErr)
			require.Equal(t, tt.want, res)
		})
	}
}

func TestService_getEventsFromTimetable(t *testing.T) {
	t.Parallel()

	type (
		args struct {
			owner timetable.CalendarOwner
		}

		fields struct {
			client *mocks.MocktimetableClient
		}
	)

	var (
		someCtx           = context.Background()
		someErr           = errors.New("")
		someEvents        = []timetable.Event{{}}
		someID     uint64 = 1
	)

	tests := []struct {
		name    string
		args    args
		fields  fields
		mockFn  func(f fields)
		want    []timetable.Event
		wantErr bool
	}{
		{
			name:    "err_educator_not_parsed",
			args:    args{owner: timetable.CalendarOwner{Class: timetable.Educator}},
			mockFn:  func(f fields) {},
			want:    nil,
			wantErr: true,
		},
		{
			name: "err_educator_not_found",
			args: args{owner: timetable.CalendarOwner{
				Class: timetable.Educator,
				Name:  "Иванов Иван Иванович",
			}},
			mockFn: func(f fields) {
				f.client.EXPECT().
					FindEducator(gomock.Any(), "Иванов", "Иван", "Иванович").
					Return(uint64(0), someErr)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "educator_events",
			args: args{owner: timetable.CalendarOwner{
				Class: timetable.Educator,
				Name:  "Петров Петр Петрович",
			}},
			mockFn: func(f fields) {
				f.client.EXPECT().
					FindEducator(gomock.Any(), "Петров", "Петр", "Петрович").
					Return(someID, nil)
				f.client.EXPECT().
					GetEducatorEvents(gomock.Any(), someID).
					Return(someEvents, nil)
			},
			want:    someEvents,
			wantErr: false,
		},
		{
			name: "err_group_not_found",
			args: args{owner: timetable.CalendarOwner{
				Class: timetable.Group,
				Name:  "неизвестная группа",
			}},
			mockFn: func(f fields) {
				f.client.EXPECT().
					FindGroup(gomock.Any(), "неизвестная группа").
					Return(uint64(0), someErr)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "group_events",
			args: args{owner: timetable.CalendarOwner{
				Class: timetable.Group,
				Name:  "21.Б15-мм",
			}},
			mockFn: func(f fields) {
				f.client.EXPECT().
					FindGroup(gomock.Any(), "21.Б15-мм").
					Return(someID, nil)
				f.client.EXPECT().
					GetGroupEvents(gomock.Any(), someID).
					Return(someEvents, nil)
			},
			want:    someEvents,
			wantErr: false,
		},
		{
			name: "classroom_events",
			args: args{owner: timetable.CalendarOwner{
				Class: timetable.Classroom,
				Name:  "мат-мех 3248",
			}},
			mockFn: func(f fields) {
				f.client.EXPECT().
					GetClassroomEvents(gomock.Any(), "мат-мех 3248").
					Return(someEvents, nil)
			},
			want:    someEvents,
			wantErr: false,
		},
		{
			name:    "err_bad_owner_class",
			args:    args{owner: timetable.CalendarOwner{}},
			mockFn:  func(f fields) {},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tt.fields = fields{
				client: mocks.NewMocktimetableClient(ctrl),
			}
			tt.mockFn(tt.fields)

			res, err := Service{
				ttClient: tt.fields.client,
			}.getEventsFromTimetable(someCtx, tt.args.owner)

			require.True(t, (err != nil) == tt.wantErr)
			require.Equal(t, tt.want, res)
		})
	}
}
