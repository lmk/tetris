// 화면이 로딩되면, websocket을 연결해서, room list를 조회한다.

// room list를 조회하면, room list를 화면에 표시한다.

var ws;
var print = function(message) {
    console.log(message);
};

var send = function(action, msg) {
  let data = {
    action: action,
    sender: $('#nick').val(),
    msg: msg
  };

  print(ws)
  print(JSON.stringify(data))

  if (ws != null && ws.readyState == 1) {
    ws.send(JSON.stringify(data));
  } else {
    print("websocket is not connected");
  }
}

var messageHandler = function(data) {
  print("RECV: " + data.action);
  print(data);
  switch (data.action) {
    case 'new-nick':
      $('#nick').val(data.msg);
      break;

    case 'list-room':
      if (data.roomList != null) {
        $('#roomList').text("")
        for (var i = 0; i < data.roomList.length; i++) {
          var room = data.roomList[i];
          $('#roomList').append('<div class="room" id="'+room.RoomId+'">' + room.RoomId + '<button class="joinButton">JOIN</button></div>');
        }
      } else {
        $('#roomList').html('<div class="no-room">' + 'No room' + '</div>')
      }
      break;

    case 'new-room':
      break;

    case 'join-room':
      if (data.sender == $('#nick').val()) {
        $('#room').hide();
        $('#game').show();
        createStartButton();
      }
      break;

    case 'leave-room':
      if (data.sender == $('#nick').val()) {
        $('#game').hide();
        $('#room').show();
        send('list-room', "") 
      }
      break;

    case 'start-game':
      window.Game.Init();
      window.Game.myBoard.Start();
  
      window.timer = window.setInterval(function()
      {
          if (window.Game.myBoard.MoveDownBlock() == false) {
              if (window.Game.myBoard.NextTern() == false) {
                send("over-game", "next-tern-fail")
              }
          }
      }, 1000);
      break;
      
    case 'over-game':
      print("over-game")
      print(data)
      if (data.sender == $('#nick').val()) {
           window.Game.GameOver();
      }
      break;

    case 'new-block':
      if (data.sender == $('#nick').val()) {
        window.Game.myBoard.AppendNextBlock(JSON.parse(data.msg));
      }
      break;

    default:
      break;
  }
}


// 화면 로딩
window.onload = function() {

  // myNick 생성
  // let myNick = 'user' + Math.floor(Math.random() * 1000);
  // $('#nick').val(myNick);

  // websocket 연결
  ws = new WebSocket("ws://localhost:8080/ws/list");
  ws.onopen = function(evt) {
    print("OPEN");

    // sleep 1초 room list 조회
    setTimeout(() => send('list-room', ""), 1000)
  }

  ws.onclose = function(evt) {
    print("CLOSE");
    ws = null;
  }

  ws.onmessage = function(evt) {
    print("RESPONSE: " + evt.data);
    var data = JSON.parse(evt.data);
    messageHandler(data);

  }
  ws.onerror = function(evt) {
      print("ERROR: " + evt.data);
  }

  // room list 조회
  //ws.send('list-room');
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
    send('new-room', "")
  });

  // joinButton 을 클릭하면 방에 입장한다.
  $('.joinButton').click(function() {
    var roomName = $(this).parent().attr('id');
    send('join-room', roomName)
  });

}


