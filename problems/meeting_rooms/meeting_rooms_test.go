package meeting_rooms

import (
	"testing"
	"time"
)

func TestMeetingRooms(t *testing.T) {
	t.Parallel()
	schedules := []MeetingSchedule{
		MeetingSchedule{
			start: time.Date(2017, 11, 10, 8, 0, 0, 0, time.UTC),
			end:   time.Date(2017, 11, 10, 9, 0, 0, 0, time.UTC),
		},
		MeetingSchedule{
			start: time.Date(2017, 11, 10, 8, 30, 0, 0, time.UTC),
			end:   time.Date(2017, 11, 10, 9, 30, 0, 0, time.UTC),
		},
		MeetingSchedule{
			start: time.Date(2017, 11, 10, 9, 30, 0, 0, time.UTC),
			end:   time.Date(2017, 11, 10, 10, 30, 0, 0, time.UTC),
		},
		MeetingSchedule{
			start: time.Date(2017, 11, 10, 9, 30, 0, 0, time.UTC),
			end:   time.Date(2017, 11, 10, 11, 00, 0, 0, time.UTC),
		},
	}

	expected := 2
	result := NumRequiredMeetingRooms(schedules)

	if expected != result {
		t.Errorf("%v != %v", expected, result)
	}
}
