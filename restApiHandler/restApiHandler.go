package restApiHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"meetingrooms/protocol"
	"meetingrooms/redisHandler"
	"meetingrooms/util"
	"net/http"
	"sort"
)

const (
	PathPreFixForBookingList = "/meetingrooms/bookingList/"
	PathForBooking = "/meetingrooms/booking"
)

type HttpHandler struct {
	storage *redisHandler.RedisHandler
}

func CreateHttpHandler(storage *redisHandler.RedisHandler) HttpHandler {
	return HttpHandler{storage}
}

func (h HttpHandler) BookingListHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("called BookingListHandler")

	date := req.URL.Path[len(PathPreFixForBookingList):]
	if len(date) != 8 {
		sendError(w, 400, "invalid date param")
		return
	}

	resBookingList, err := h.buildResBookingList(date)
	if err == nil {
		sendResponse(w, resBookingList)
	} else {
		sendError(w, 500, err.Error())
	}
}

func (h HttpHandler) BookingHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("called BookingHandler")

	reqBooking := protocol.ReqBooking{}
	err := json.NewDecoder(req.Body).Decode(&reqBooking)
	if err != nil {
		sendError(w, 400, "invalid request")
		return
	}

	if err := h.storage.CanBooking(reqBooking); err != nil {
		sendError(w, 409, err.Error())
		return
	}

	err = h.storage.AddBooking(reqBooking)
	if err != nil {
		sendError(w, 500, err.Error())
		return
	}

	resBookingList, err := h.buildResBookingList(reqBooking.Date)
	if err == nil {
		sendResponse(w, resBookingList)
	} else {
		sendError(w, 500, err.Error())
	}
}

func (h HttpHandler) buildResBookingList(date string) (*protocol.ResBookingList, error) {
	if rooms,err := h.storage.GetMeetingRooms(); err == nil {
		if len(rooms) == 0 {
			return nil, errors.New("invalid date param")
		}

		sort.Strings(rooms)
		resBookingList := protocol.ResBookingList{}
		for _, roomName := range rooms {
			if roomInfo, err := h.buildRoomInfo(roomName, date); err == nil {
				resBookingList.MeetingRooms = append(resBookingList.MeetingRooms, *roomInfo)
			}
		}
		return &resBookingList, nil
	} else {
		return nil, errors.New("db error")
	}
}

func (h HttpHandler) buildRoomInfo(roomName string, date string) (*protocol.RoomInfo, error) {
	startTimes, err := h.storage.GetStartTimes(roomName, date)
	if err != nil {
		return nil, err
	}

	sort.Strings(startTimes)
	roomInfo := protocol.RoomInfo {Name: roomName}
	for _, timeInfo := range startTimes {
		startTime, userName := util.SplitStartTimeAndUserName(timeInfo)
		roomInfo.StartTimes = append(roomInfo.StartTimes, protocol.StartTimeInfo {startTime, userName})
	}
	return &roomInfo, nil
}

func sendResponse(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		sendError(w, 400, "failed to prepare json")
	}
}

func sendError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	fmt.Fprintln(w, message)
}
