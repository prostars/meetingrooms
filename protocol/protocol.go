package protocol

type StartTimeInfo struct {
	StartTime string
	UserName string
}
type RoomInfo struct {
	Name string
	StartTimes []StartTimeInfo
}
type ResBookingList struct {
	MeetingRooms []RoomInfo
}

type ReqBooking struct {
	RoomName string
	Date string
	StartTime string
	EndTime string
	RepeatCount int
	UserName string
}