package redisHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"meetingrooms/protocol"
	"meetingrooms/util"
)

type RedisHandler struct {
	conn redis.Conn
}

type BookingInfo struct {
	StartTime string
	EndTime string
	RepeatCount int
	UserName string
}

func CreateRedisHandler() *RedisHandler {
	return &RedisHandler{newPool().Get()}
}

func (h RedisHandler) Close() {
	h.conn.Close()
}

func (h RedisHandler) AddMeetingRoom(name string) error {
	_, err := h.conn.Do("SADD", "meetingRooms", name)
	if err != nil {
		return err
	}
	return nil
}

func (h RedisHandler) GetMeetingRooms() ([]string, error) {
	names, err := redis.Strings(h.conn.Do("SMEMBERS", "meetingRooms"))
	if err != nil {
		return nil, err
	}
	return names, nil
}

func (h RedisHandler) CanBooking(r protocol.ReqBooking) error {
	if bookedTimes, err := h.GetBookingTimes(r.RoomName, r.Date, r.StartTime, r.EndTime); err == nil {
		if len(bookedTimes) == 0 {
			return nil
		} else {
			return errors.New("already booked")
		}
	} else {
		return err
	}
}

func (h RedisHandler) AddBooking(r protocol.ReqBooking) error {
	if r.RepeatCount < 1 {
		return errors.New("invalid repeatCount")
	}

	if err := h.CanBooking(r); err != nil {
		return err
	}

	bookingTimes := util.BuildBookingTimes(r.StartTime, r.EndTime)
	bookingDate := r.Date
	for cnt := 0; cnt < r.RepeatCount; cnt++ {
		if err := h.addBookingInfo(r.RoomName, bookingDate, r.StartTime, r.EndTime, r.RepeatCount, r.UserName); err != nil {
			return err
		}
		if err := h.addBookingTimes(r.RoomName, bookingDate, r.UserName, bookingTimes...); err != nil {
			return err
		}
		bookingDate = util.GetNextWeekDate(bookingDate)
		if len(bookingDate) == 0 {
			return errors.New("failed date calculation")
		}
	}
	return nil
}

func (h RedisHandler) GetBookingTimes(meetingRoomName string, date string, startTime string, endTime string) ([]string, error) {
	times, err := redis.Strings(h.conn.Do("ZRANGEBYSCORE", makeBookingTimesKey(meetingRoomName, date), startTime, endTime))
	if err != nil {
		return nil, err
	}
	return times, nil
}

func (h RedisHandler) GetStartTimes(meetingRoomName string, date string) ([]string, error) {
	startTimes, err := redis.Strings(h.conn.Do("HKEYS", makeBookingInfoKey(meetingRoomName, date)))
	if err != nil {
		return nil, err
	}
	return startTimes, nil
}

func (h RedisHandler) GetBookingInfo(meetingRoomName string, date string, startTime string) (*BookingInfo, error) {
	value, err := redis.Strings(h.conn.Do("HMGET", makeBookingInfoKey(meetingRoomName, date), startTime))
	if err != nil {
		return nil, err
	}
	if len(value[0]) == 0 {
		return nil, errors.New("can not found booking info")
	}
	bookingInfo := BookingInfo{}
	if err := json.Unmarshal([]byte(value[0]), &bookingInfo); err != nil {
		return nil, err
	}
	return &bookingInfo, nil
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 80,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func (h RedisHandler) addBookingTimes(meetingRoomName string, date string, userName string, times ...string) error {
	for _, time := range times {
		_, err := h.conn.Do("ZADD", makeBookingTimesKey(meetingRoomName, date), time, makeBookingTimeInfo(time, userName))
		if err != nil {
			return err
		}
	}
	return nil
}

func (h RedisHandler) addBookingInfo(meetingRoomName string, date string, startTime string, endTime string, repeatCount int, userName string) error {
	if value, err := json.Marshal(BookingInfo{startTime,endTime, repeatCount, userName}); err == nil {
		key := makeBookingTimeInfo(startTime, userName)
		_, err := h.conn.Do("HMSET", makeBookingInfoKey(meetingRoomName, date), key, value)
		if err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

func makeBookingTimesKey(meetingRoomName string, date string) string {
	return fmt.Sprintf("bookingTimes:%s:%s", meetingRoomName, date)
}

func makeBookingInfoKey(meetingRoomName string, date string) string {
	return fmt.Sprintf("bookingInfo:%s:%s", meetingRoomName, date)
}

func makeBookingTimeInfo(startTime string, userName string) string {
	return fmt.Sprintf("%s:%s", startTime, userName)
}