package meeting_rooms

import (
	"testing"
	"time"
)

func TestMeetingRooms(t *testing.T) {
	t.Parallel()
	schedules := []MeetingSchedule{
		MeetingSchedule{
			start: time.Date(2017, 11, 10, 1, 0, 0, 0, time.UTC),
			end:   time.Date(2017, 11, 10, 2, 30, 0, 0, time.UTC),
		},
		MeetingSchedule{
			start: time.Date(2017, 11, 10, 2, 0, 0, 0, time.UTC),
			end:   time.Date(2017, 11, 10, 3, 0, 0, 0, time.UTC),
		},
		MeetingSchedule{
			start: time.Date(2017, 11, 10, 1, 30, 0, 0, time.UTC),
			end:   time.Date(2017, 11, 10, 2, 15, 0, 0, time.UTC),
		},
		MeetingSchedule{
			start: time.Date(2017, 11, 10, 2, 30, 0, 0, time.UTC),
			end:   time.Date(2017, 11, 10, 3, 30, 0, 0, time.UTC),
		},
	}

	expected := 3
	result := NumRequiredMeetingRooms(schedules)

	if expected != result {
		t.Errorf("%v != %v", expected, result)
	}
}
