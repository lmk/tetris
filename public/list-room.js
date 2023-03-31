// 화면이 로딩되면, websocket을 연결해서, room list를 조회한다.

// room list를 조회하면, room list를 화면에 표시한다.

var ws;
var print = function(message) {
    console.log(message);
};

var send = function(action, msg) {
  let data = {
    action: action,
    sender: $('#nick').val,
    msg: msg
  };

  print(ws)
  print(JSON.stringify(data))

  ws.send(JSON.stringify(data));
}

// 화면 로딩
window.onload = function() {

  // myNick 생성
  let myNick = 'user' + Math.floor(Math.random() * 1000);
  $('#nick').val(myNick);

  // websocket 연결
  ws = new WebSocket("ws://localhost:8080/ws/list");
  ws.onopen = function(evt) {
    print("OPEN");

    // sleep 1초 room list 조회
   // setTimeout(() => send('list-room', ""), 1000)
  }

  ws.onclose = function(evt) {
    print("CLOSE");
    ws = null;
  }

  ws.onmessage = function(evt) {
    print("RESPONSE: " + evt.data);
  }
  ws.onerror = function(evt) {
      print("ERROR: " + evt.data);
  }

  // room list 조회
  //ws.send('list-room');
  //socket.emit('list-room');
/*
  // room list 화면 표시
  socket.on('list-room', function(data) {
    var roomList = data.roomList;
    var roomListElement = document.getElementById('roomList');
    for (var i = 0; i < roomList.length; i++) {
      var room = roomList[i];
      $('#roomList').append('<div class="room" id="'+room.name+'">' + room.name + '<button class="joinButton">JOIN</button></div>');
    }
  });
*/
  // newGame 버튼 을 클릭하면 새로운 방을 만든다.
  var newGameButton = document.getElementById('newGame');
  newGameButton.addEventListener('click', function() {
    socket.emit('new-room');
  });

  // joinButton 을 클릭하면 방에 입장한다.
  $('.joinButton').click(function() {
    var roomName = $(this).parent().attr('id');
    socket.emit('join-room', {name: roomName});
  });

}


