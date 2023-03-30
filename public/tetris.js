// 게임 화면 안에 시작 버튼 생성
function createStartButton() {
    var startButton = document.createElement('button');
    startButton.id = 'startButton';
    startButton.innerHTML = 'START';
    startButton.addEventListener('click', startGame);
    document.getElementById('game').appendChild(startButton);
}

// 화면이 처음 로딩 되면 시작 버튼을 생성
window.onload = function() {
    createStartButton();
}

// 게임 시작
function startGame() {
    // 게임 화면에 있는 시작 버튼을 삭제
    document.getElementById('game').removeChild(document.getElementById('startButton'));
    // 게임 화면에 블록을 생성
    createBlock();
    // 게임 화면에 블록을 이동
    moveBlock();
}

// 게임 화면에 블록을 생성
function createBlock() {
    // 게임 화면에 블록을 생성
    var block = document.createElement('div');
    block.id = 'block';
    document.getElementById('game').appendChild(block);
}

// 게임 화면에 블록을 이동
function moveBlock() {
    // 게임 화면에 블록을 이동
    var block = document.getElementById('block');
    block.style.top = '0px';
    block.style.left = '0px';
    var position = 0;
    setInterval(function() {
        position += 1;
        block.style.left = position + 'px';
    }, 10);
}

