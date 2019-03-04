package redisHandler

import (
	"github.com/stretchr/testify/assert"
	"meetingrooms/protocol"
	"testing"
)

var storage *RedisHandler

func TestMain(m *testing.M) {
	storage = CreateRedisHandler()
	m.Run()
	storage.Close()
}

func TestAddMeetingRoom(t *testing.T) {
	assert.Equal(t, storage.conn.Err(), nil)
	_, err := storage.conn.Do("SREM", "meetingRooms", "testRoomA")
	assert.Equal(t, err, nil)
	_, err = storage.conn.Do("SREM", "meetingRooms", "testRoomB")
	storage.conn.Do("SREM", "meetingRooms", "testRoomB")
	assert.Equal(t, storage.AddMeetingRoom("testRoomA"), nil)
	assert.Equal(t, storage.AddMeetingRoom("testRoomB"), nil)
	rooms, err := storage.GetMeetingRooms()
	assert.Equal(t, err, nil)
	cnt := 0
	for _, roomName := range rooms {
		switch roomName {
		case "testRoomA":
			cnt++
		case "testRoomB":
			cnt++
		}
	}
	assert.Equal(t, cnt, 2)
}

func TestAddBooking(t *testing.T) {
	assert.Equal(t, storage.conn.Err(), nil)
	_, err := storage.conn.Do("DEL", makeBookingInfoKey("testRoomA", "20190304"))
	assert.Equal(t, err, nil)
	_, err = storage.conn.Do("DEL", makeBookingTimesKey("testRoomA", "20190304"))
	assert.Equal(t, err, nil)

	assert.Equal(t, storage.AddMeetingRoom("testRoomA"), nil)
	reqBooking := protocol.ReqBooking{RoomName:"testRoomA", Date:"20190304", StartTime:"1600", EndTime:"2030", RepeatCount:1, UserName:"prostars"}
	assert.Equal(t, storage.AddBooking(reqBooking), nil)
	bookingInfo, err := storage.GetBookingInfo("testRoomA", "20190304", "1600:prostars")
	assert.Equal(t, err, nil)
	assert.Equal(t, bookingInfo.StartTime, reqBooking.StartTime)
	assert.Equal(t, bookingInfo.EndTime, reqBooking.EndTime)
	assert.Equal(t, bookingInfo.RepeatCount, reqBooking.RepeatCount)
	assert.Equal(t, bookingInfo.UserName, reqBooking.UserName)
}