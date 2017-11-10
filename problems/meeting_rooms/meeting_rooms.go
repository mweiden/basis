package meeting_rooms

import (
	"fmt"
	"sort"
	"time"
)

type MeetingSchedule struct {
	start time.Time
	end   time.Time
}

type scheduleEvent struct {
	eventTime time.Time
	isStart   bool
}

func NumRequiredMeetingRooms(schedules []MeetingSchedule) int {
	var timeQueue []scheduleEvent

	// load the meeting schedules into an array of events
	for _, schedule := range schedules {
		startEvent := scheduleEvent{
			eventTime: schedule.start,
			isStart:   true,
		}
		endEvent := scheduleEvent{
			eventTime: schedule.end,
			isStart:   false,
		}
		timeQueue = append(timeQueue, startEvent)
		timeQueue = append(timeQueue, endEvent)
	}

	// sort the array of events, earlier events and endings first
	sort.Slice(timeQueue, func(i, j int) bool {
		diff := timeQueue[j].eventTime.Unix() - timeQueue[i].eventTime.Unix()
		return diff > 0 || (diff == 0 && !timeQueue[i].isStart && timeQueue[j].isStart)
	})

	fmt.Printf(">> %v\n", timeQueue)

	roomsRequired := 0
	maxRooms := 0

	// iterate through, track max rooms
	for _, event := range timeQueue {
		if event.isStart {
			roomsRequired++
		} else {
			roomsRequired--
		}
		if roomsRequired > maxRooms {
			maxRooms = roomsRequired
		}
	}
	return maxRooms
}
