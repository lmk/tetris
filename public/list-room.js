// 화면이 로딩되면, websocket을 연결해서, room list를 조회한다.

// room list를 조회하면, room list를 화면에 표시한다.

var socket = null;

// 화면 로딩
window.onload = function() {

  // myNick 생성
  var myNick = 'user' + Math.floor(Math.random() * 1000);
  $('#myNick').html(myNick);

  // websocket 연결
  socket =


  //socket = io.connect('http://localhost:3000');

  // room list 조회
  socket.emit('list-room');

  // room list 화면 표시
  socket.on('list-room', function(data) {
    var roomList = data.roomList;
    var roomListElement = document.getElementById('roomList');
    for (var i = 0; i < roomList.length; i++) {
      var room = roomList[i];
      $('#roomList').append('<div class="room" id="'+room.name+'">' + room.name + '<button class="joinButton">JOIN</button></div>');
    }
  });

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


