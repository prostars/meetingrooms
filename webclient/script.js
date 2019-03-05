function pad(n, width) {
    n = n + '';
    return n.length >= width ? n : new Array(width - n.length + 1).join('0') + n;
}

function prepareDatePicker() {
    document.getElementById('selectedDate').valueAsDate = new Date();
    document.getElementById('bookingDate').valueAsDate = new Date();
}

function prepareTimeSelector() {
    for (let hour = 8; hour < 24; hour++)
        $('#startTimeHour').append('<option>' + hour + '</option>');

    $('#startTimeMinute').append('<option>00</option>');
    $('#startTimeMinute').append('<option>30</option>');

    const durationTime = $('#durationTime');
    let duration = 0;
    for (let i = 0; i < 6; i++) {
        if (duration % 100 === 0)
            duration += 30;
        else
            duration += 70;
        const time = pad(duration, 3);
        const hourAndMinute = [time.slice(0, 1), ':', time.slice(1)].join('');
        durationTime.append('<option>' + hourAndMinute + '</option>');
    }

    const repeatCount = $('#repeatCount');
    for (let i = 1; i < 11; i++)
        repeatCount.append('<option>' + i + '</option>');
}

function getDateString(datePicker) {
    const date = new Date(datePicker.val());
    const year = date.getFullYear();
    const month = pad(date.getMonth() + 1, 2);
    const day = pad(date.getDate(), 2);
    return [year, month, day].join('');
}

function requestBookingList(date) {
    const reqUrl = 'http://localhost:30001/meetingrooms/bookinglist/' + date;
    $.get(reqUrl, function(data, status) {
        console.log('received response.' + status);
        console.log(data);
        buildSelectBoxForMeetingRooms(data);
        buildTableForMeetingRooms(data);
    });
}

function requestBooking() {
    console.log('called requestBooking');
    let bookingInfo = new Object();

    const startTime = pad($('#startTimeHour').val(), 2) + $('#startTimeMinute').val();
    const time = $('#durationTime').val();
    const duration = [time.slice(0, 1), time.slice(2, 4)].join('');
    let endTime = Number(startTime) + Number(duration);
    endTime = endTime % 100 === 60 ? endTime + 40 : endTime;

    bookingInfo.RoomName = $('#selectedMeetingRoom').val();
    bookingInfo.Date = getDateString($('#bookingDate'));
    bookingInfo.StartTime = startTime;
    bookingInfo.EndTime = pad(endTime, 4);
    bookingInfo.RepeatCount = Number($('#repeatCount').val());
    bookingInfo.UserName = $('#userName').val();

    if (bookingInfo.UserName.length === 0) {
        alert("UserName is required.");
        return;
    }

    let json = JSON.stringify(bookingInfo);
    console.log(bookingInfo);
    console.log(json);
    $.ajax('http://localhost:30001/meetingrooms/booking', {
        type: 'PUT',
        data: json,
        contentType: 'application/json',
        success: function(data, textStatus, xhr) {
            console.log('success booking request.');
            const bookingDate = [bookingInfo.Date.slice(0, 4), bookingInfo.Date.slice(4, 6), bookingInfo.Date.slice(6, 8)].join('-');
            document.getElementById('selectedDate').valueAsDate = new Date(bookingDate);
            requestBookingList(bookingInfo.Date);
        },
        error: function(xhr, textStatus, errorThrown) {
            console.log('failed booking request.');
            console.log(xhr);
            console.log(textStatus);
            console.log(errorThrown);
            alert('invalid request');
        }
    });
}

function buildSelectBoxForMeetingRooms(data) {
    const selectedMeetingRoom = $('#selectedMeetingRoom');
    selectedMeetingRoom.empty();
    const jsonData = JSON.parse(data);
    const roomCount = Object.keys(jsonData['MeetingRooms']).length;
    console.log(roomCount);
    for (let i = 0; i < roomCount; i++) {
        const roomName = jsonData['MeetingRooms'][i]['Name'];
        console.log(roomName);
        selectedMeetingRoom.append('<option>' + roomName + '</option>');
    }
}

function buildTableForMeetingRooms(data) {
    console.log("called buildTableForMeetingRooms.");
    const jsonData = JSON.parse(data);
    console.log(jsonData['MeetingRooms']);

    $('#meetingRooms tr').remove();

    let titleRow = '<tr>';
    let dataRow = '<tr>';
    const roomCount = Object.keys(jsonData['MeetingRooms']).length;
    for (let i = 0; i < roomCount; i++) {
        const roomName = jsonData['MeetingRooms'][i]['Name'].replace(/\s/g, "");
        titleRow += '<td>' + roomName + '</td>';
        dataRow += '<td><div id="' + roomName + '"></div></td>';
    }
    titleRow += '</tr>';
    const meetingRooms = $('#meetingRooms');
    meetingRooms.append(titleRow);
    meetingRooms.append(dataRow);


    for (let i = 0; i < roomCount; i++) {
        const roomName = '#' + jsonData['MeetingRooms'][i]['Name'].replace(/\s/g, "");
        if (jsonData['MeetingRooms'][i]['StartTimes'] == null)
            continue;
        const startTimes = Object.keys(jsonData['MeetingRooms'][i]['StartTimes']).length;
        for (let j = 0; j < startTimes; j++) {
            const startTime = jsonData['MeetingRooms'][i]['StartTimes'][j]['StartTime'];
            const userName = jsonData['MeetingRooms'][i]['StartTimes'][j]['UserName'];
            console.log(roomName);
            console.log(startTime + ':' + userName);
            $(roomName).html(startTime + ':' + userName);
        }
    }

}