package main

import (
	"fmt"
	"log"
	"meetingrooms/protocol"
	"meetingrooms/redisHandler"
	"meetingrooms/restApiHandler"
	"net/http"
	"sort"
)

func main() {
	startServer()
}

func prepare(storage *redisHandler.RedisHandler) {
	for ch := 'A'; ch <= 'G'; ch++ {
		if err := storage.AddMeetingRoom(fmt.Sprintf("MeetingRoom %c", ch)); err != nil {
			log.Fatal("failed prepare MeetingRooms.\nServer shutting down.")
			return
		}
	}

	storage.AddBooking(protocol.ReqBooking{"MeetingRoom A", "20190305", "0800", "1130", 3, "prostars"})
	storage.AddBooking(protocol.ReqBooking{"MeetingRoom A", "20190305", "1500", "1930", 1, "prostars"})
	times, _ := storage.GetStartTimes("MeetingRoom A", "20190305")
	sort.Strings(times)
	fmt.Println(times)
	fmt.Println(storage.GetBookingTimes("MeetingRoom A", "20190305", "1200", "1230"))
	fmt.Println(storage.GetBookingTimes("MeetingRoom A", "20190305", "0800", "1230"))
	fmt.Println(storage.GetBookingTimes("MeetingRoom A", "20190305", "1200", "2000"))
	fmt.Println(storage.GetBookingInfo("MeetingRoom A", "20190305", "0800:prostars"))
	fmt.Println(storage.GetBookingInfo("MeetingRoom A", "20190312", "0800:prostars"))
	fmt.Println(storage.GetBookingInfo("MeetingRoom A", "20190314", "0800:prostars"))
}

func startServer() {
	fmt.Printf("start booking server\n")

	storage := redisHandler.CreateRedisHandler()
	defer storage.Close()

	prepare(storage)

	fmt.Println(storage.GetMeetingRooms())

	apiHandler := restApiHandler.CreateHttpHandler(storage)

	http.HandleFunc(restApiHandler.PathPreFixForBookingList, apiHandler.BookingListHandler)
	http.HandleFunc(restApiHandler.PathForBooking, apiHandler.BookingHandler)

	log.Fatal(http.ListenAndServe(":30001", nil))
}
