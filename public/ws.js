// 화면이 로딩되면, websocket을 연결해서, room list를 조회한다.

// room list를 조회하면, room list를 화면에 표시한다.

var ws;
var print = function(message) {
    console.log(message);
};

var send = function(action, msg) {
  let data = {
    action: action,
    sender: $('#my-nick').text(),
    data: msg,
    date: Date.now(),
  };

  let json = JSON.stringify(data)

  print("REQUEST: " + json);

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
     break;

    case 'set-nick':
      if (msg.sender == $('#my-nick').text()) {
        window.Game.SetOwner(msg.data);
      } else {
        // 대기실에 nick을 표시할 경우 변경해야 한다.
      }
      
      break;

    case 'list-room':
      $('#roomList').text("");
      let roomCount = 0;

      for (let i = 0; i < msg.roomList.length; i++) {
        let room = msg.roomList[i];

        if (room.roomId == 0) {
          continue;
        }

        let html = '<div class="room" id="'+room.roomId+'"><div class="roomId">'+room.roomId+'</div>'
                    + '<div class="roomTitle">'+room.title+'</div>'
                    + '<div class="roomUsers">'+Object.keys(room.nicks)+'</div>'
                    + '<button class="joinButton btn btn-primary">JOIN</button></div>';
        $('#roomList').append(html);
        roomCount++;
      }

      if (roomCount > 0) {
        // joinButton 을 클릭하면 방에 입장한다.
        $('.joinButton').click(function() {
          send('join-room', $(this).parent().attr('id'))
        });
      } else {
        $('#roomList').html('<div class="no-room">' + 'No room' + '</div>')
      }
      break;

    case 'new-room':
      break;

    case 'join-room':
      if (msg.sender == $('#my-nick').text()) {
        // 방에 입장하면, 게임 화면으로 이동한다.
        $('#room').hide();
        $('#game').show();
        $('#roomId').text(msg.roomId);        
        $('#roomTitle').text(msg.roomList[0].title);

        if (msg.roomList[0].owner == $('#my-nick').text()) {
          createStartButton();
          createAddBotButton();
        }
      }

      for (let nick of Object.keys(msg.roomList[0].nicks)) {
        let el = $('#enermy-game-'+nick);
        if ($('#my-nick').text() != nick && el.length == 0) {
          let html = '<span id="enermy-game-'+nick+'" class="col-sm-4 column">'
              + '<div id="enermy-name">'+nick+'</div>'
              + '<div id="enermy-score-'+nick+'">0</div>'
              + '<div id="enermy-board-'+nick+'" class="enermy-board"></span></div>';
          $('#enermy-boards').append(html)

          window.Game.enermyBoards.set(nick, new Board());
          window.Game.enermyBoards.get(nick).InitEnermy("#enermy-board-"+nick, "#enermy-score-"+nick);
          window.Game.enermyBoards.get(nick).owner = nick;
        }
      }
      break;

    case 'leave-room':
      if (msg.sender == $('#my-nick').text()) {
        $('.triangle-container').hide();
        $('#game').hide();
        $('#room').show();
        send('list-room', "") 
      } else {
        $('#enermy-game-'+msg.sender).remove();
      }

      if (msg.roomId != 0 && msg.roomList[0].owner == $('#my-nick').text() && msg.roomList[0].state == 'ready') {
        createStartButton();
        createAddBotButton();
      }
      break;

    case 'start-game':
      $('.triangle-container').show();

      window.Game.Init();
      window.Game.Start(msg.sender);

      window.Game.myBoard.cells = msg.cells;
      DrawBoard(msg.cells, msg.block);            

      window.Game.myBoard.SetNextBlock(msg.blockIndexs);
      DrawNextBlocks();
      break;
      
    case 'over-game':

      $('.triangle-container').hide();

      if (msg.sender == $('#my-nick').text() && window.Game.myBoard.IsPlaying()) {

        window.Game.myBoard.cells = msg.cells;
        DrawBoard(msg.cells, msg.block);

        window.Game.GameOver(msg.sender);

        $('#modal-over-game').modal('show');

        if (msg.roomList[0].owner == $('#my-nick').text()) {
          createStartButton();
          createAddBotButton();
        }
      } else {
        window.Game.GameOver(msg.sender);
      }
      break;

    case 'end-game':

      $('.triangle-container').hide();
      
      if (msg.data > 0) {
        $('#winner-nick').text(msg.sender);
        $('#winner-rank').text(msg.data);
        $('#winner-score').text((msg.score==undefined)?0:msg.score);
        $('#modal-over-game').modal('hide');
        $('#modal-winner').modal('show');
      } else {
        $('#modal-over-game').modal('show');
        window.Game.GameOver(msg.sender);
      }

      if (msg.roomList[0].owner == $('#my-nick').text()) {
        createStartButton();
        createAddBotButton();
      }      
      break;

    case 'sync-game':
      if (msg.sender == $('#my-nick').text()) {

        window.Game.myBoard.SetNextBlock(msg.blockIndexs);
        DrawNextBlocks();

        $('#score').text(msg.score);

        if (window.Game.myBoard.IsPlaying()) {
          window.Game.myBoard.cells = msg.cells;
          DrawBoard(msg.cells, msg.block);
        }
      } else {
        $('#enermy-score-'+msg.sender).text(msg.score);
        if (window.Game.enermyBoards.has(msg.sender)) {
          DrawEnermyBoard(msg.sender, msg.cells, msg.block);
        }
      }

      break;

    case 'list-rank':
      $('#topRank').html("");

      if (msg.rankList.length > 0) {

        let html = '<ul class="list-group"><li class="list-group-item active list-group-item-light" aria-current="true">Top ' + msg.rankList.length + '</li>';

        for (let i = 0; i < msg.rankList.length; i++) {
          let rank = msg.rankList[i];
          html += '<li class="list-group-item text-center d-flex .bg-light">'
                   +'<span id="rank-num" class="rank p-2">' + rank.rank + '</span>'
                   +'<span id="rank-nick" class="p-2">' + rank.nick + '</span>'
                   +'<span id="rank-score" class="p-2">' + rank.score + '</span>'
                   +'<span id="rank-date" class="p-2"> '+ rank.date +'</span></li> ';
        }

        html += '</ul>';
        $('#topRank').append(html);
  
      }
      break;

    case 'gift-full-blocks':
      window.Game.PlaySound(msg.action);
      break;

    case 'erase-blocks':
      window.Game.PlaySound(msg.action);
      break;

    default:
      print("Unknown message: " + msg.action)
      break;
  }
}


// 화면 로딩
window.onload = function() {

  if (window.Game == undefined) {
    CreateGame();
  }  

  // websocket 연결
  ws = new WebSocket($('#wsServer').val());

  ws.onopen = function(evt) {
    print("OPEN");

    if ($.cookie('nick')!=undefined) {
      setTimeout(() => send('set-nick', $.cookie('nick')), 500) 
    }

    setTimeout(() => send('list-rank', '5'), 600)
    setTimeout(() => send('list-room', ''), 600)
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
    $.cookie('nick', $('#newNick').val(), { expires: 30 });
  });


}


