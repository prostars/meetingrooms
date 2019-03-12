# Meeting Rooms (RESTful API example using golang and redis)

# Install
1. Golang 1.12.x [Download](https://golang.org/dl/)
    * [SettingGOPATH](https://github.com/golang/go/wiki/SettingGOPATH)
2. Redis 5.0.x [Download](https://redis.io/download)
    * You need to build redis. Please refer to the documentation.

# Build
1. Go to the directory you checked out of this project.
2. execute command    
```$ go build```

# Run unit tests
```$ ./go test ./...```

# Run server
```$ ./meetingroom```

# Test server with curl
1. [Download](https://curl.haxx.se/download.html) a curl if there is no that.
2. Request booking list for 2019-03-05
    ```
    $ curl -v http://localhost:30001/meetingrooms/bookinglist/20190305
    ```
3. You can get response for that.  
    ```
    {  
       "MeetingRooms":[  
          {  
             "Name":"MeetingRoom A",
             "StartTimes":[  
                {  
                   "StartTime":"0800",
                   "UserName":"prostars"
                },
                {  
                   "StartTime":"1500",
                   "UserName":"prostars"
                }
             ]
          },
          {  
             "Name":"MeetingRoom B",
             "StartTimes":null
          },
          {  
             "Name":"MeetingRoom C",
             "StartTimes":null
          },
          {  
             "Name":"MeetingRoom D",
             "StartTimes":null
          },
          {  
             "Name":"MeetingRoom E",
             "StartTimes":null
          },
          {  
             "Name":"MeetingRoom F",
             "StartTimes":null
          },
          {  
             "Name":"MeetingRoom G",
             "StartTimes":null
          }
       ]
    }    
    ```  
4. Request new booking for MeetingRoom C at 2019-03-05 13:00 ~ 14:30.
   ```
   curl -v -H 'Content-Type: application/json' -X PUT -d '{"RoomName":"MeetingRoom C","Date":"20190305","StartTime":"1300","EndTime":"1430","RepeatCount":2,"UserName":"guinness"}' http://localhost:30001/meetingrooms/booking
   ```    
5. You can get response for that.
    ```
     {  
        "MeetingRooms":[  
           {  
              "Name":"MeetingRoom A",
              "StartTimes":[  
                 {  
                    "StartTime":"0800",
                    "UserName":"prostars"
                 },
                 {  
                    "StartTime":"1500",
                    "UserName":"prostars"
                 }
              ]
           },
           {  
              "Name":"MeetingRoom B",
              "StartTimes":null
           },
           {  
              "Name":"MeetingRoom C",
              "StartTimes":[  
                 {  
                    "StartTime":"1300",
                    "UserName":"guinness"
                 }
              ]
           },
           ...
           ...
        ]
     }  
    ```
6. You have requested a repeat count of 2, so you can see other booking.
    ```
    {  
       "MeetingRooms":[  
          ...
          ...
          {  
             "Name":"MeetingRoom C",
             "StartTimes":[  
                {  
                   "StartTime":"1300",
                   "UserName":"guinness"
                }
             ]
          },
          ...
          ...
       ]
    }
    ```    
# Use chrome browser
1. input url 'localhost:30001' on your chrome browser    
