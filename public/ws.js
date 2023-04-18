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
    data: msg
  };

  let json = JSON.stringify(data)

  print(json)

  if (ws != null && ws.readyState == 1) {
    ws.send(json);
  } else {
    print("websocket is not connected");
  }
}

var messageHandler = function(msg) {
  switch (msg.action) {

    case 'error':
      alert(msg.data);
      break;

    case 'new-nick':
      window.Game.SetOwner(msg.data);

          // sleep 1초 room list 조회
      setTimeout(() => send('list-room', ""), 1000)
      break;

    case 'set-nick':
      window.Game.SetOwner(msg.data);
      break;

    case 'list-room':
      if (msg.roomList != null) {
        $('#roomList').text("")
        for (var i = 0; i < msg.roomList.length; i++) {
          var room = msg.roomList[i];
          $('#roomList').append('<div class="room" id="'+room.RoomId+'">' + room.RoomId + '<button class="joinButton">JOIN</button></div>');
        }
      } else {
        $('#roomList').html('<div class="no-room">' + 'No room' + '</div>')
      }
      break;

    case 'new-room':
      break;

    case 'join-room':
      if (msg.sender == $('#nick').val()) {
        $('#room').hide();
        $('#game').show();
        createStartButton();
      }
      break;

    case 'leave-room':
      if (msg.sender == $('#nick').val()) {
        $('#game').hide();
        $('#room').show();
        send('list-room', "") 
      }
      break;

    case 'start-game':
      window.Game.Init();
      window.Game.Start(msg.sender);
      DrawNextBlocks();
      break;
      
    case 'over-game':
      if (msg.sender == $('#nick').val() && window.Game.myBoard.IsPlaying()) {
        window.Game.GameOver(msg.sender);
        alert("Game over, please refresh the page to start new game");
        createStartButton();
      } else {
        window.Game.GameOver(msg.sender);
      }
      break;

    case 'sync-game':
      if (msg.sender == $('#nick').val()) {
        $('#score').text(msg.score);
        window.Game.myBoard.cells = msg.cells;
        DrawBoard(msg.cells, msg.block);
      }
      break;

    case 'next-block':
      window.Game.myBoard.SetNextBlock(JSON.parse(msg.data));
      if (window.Game.myBoard.IsPlaying()) {
        DrawNextBlocks();
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

  if (window.Game == undefined) {
    CreateGame();
  }  

  // websocket 연결
  ws = new WebSocket("ws://localhost:8080/ws/list");
  ws.onopen = function(evt) {
    print("OPEN");

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

  // newGame 버튼 을 클릭하면 새로운 방을 만든다.
  $('#newGame').click(function() {
    send('new-room', "")
  });

  $('#setNick').click(function() {
    send('set-nick', $('#newNick').val())
  });

  // joinButton 을 클릭하면 방에 입장한다.
  $('.joinButton').click(function() {
    var roomName = $(this).parent().attr('id');
    send('join-room', roomName)
  });

}


