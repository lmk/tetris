// 게임 화면 안에 시작 버튼 생성
function createStartButton() {
    if ($('#startButton').length != 0) return;

    let button = $('<button id="startButton" class="btn btn-primary"/>')
                    .text('START')
                    .click(function() {
                        send("start-game", "");

                        // delete button
                        $('#startButton').remove();
                    });
    $('#startButtonArea').append(button);
}

class Game {

    row = 15;
    column = 10;
    
    EMPTY = 0;

    SHAPES = [
        [
            [0,0,0,0],
            [1,1,1,1],
            [0,0,0,0],
            [0,0,0,0],
        ], [            
            [0,0,0,0],
            [1,1,1,0],
            [0,1,0,0],
            [0,0,0,0],

        ], [
            [0,0,0,0],
            [0,1,1,0],
            [0,1,1,0],
            [0,0,0,0],
        ], [
            [0,1,0,0],
            [0,1,1,0],
            [0,0,1,0],
            [0,0,0,0]
        ], [
            [0,1,0,0],
            [0,1,0,0],
            [0,1,1,0],
            [0,0,0,0]
        ], [
            [0,0,1,0],
            [0,0,1,0],
            [0,1,1,0],
            [0,0,0,0]
        ], [
            [0,0,1,0],
            [0,1,1,0],
            [0,1,0,0],
            [0,0,0,0]
        ],
    ];    

    myBoard = undefined;
    enermyBoards = new Map();

    Init() {

        if (this.myBoard == undefined) {
            this.myBoard = new Board();
        }

        this.myBoard.Init("#my-board", "#my-score", "#my-nextBlock");
    }

    SetOwner(owner) {
        this.myBoard.owner = owner;
        $('#my-nick').text(owner);
        $('#newNick').val(owner);        
    }    

    Start(sender) {
        if (sender == this.myBoard.owner) {
            this.myBoard.Start();
        } else if (this.enermyBoards.has(sender)) {
            this.enermyBoards.get(sender).Start();
        }
    }

    GameOver(sender) {
        if (sender == this.myBoard.owner) {
            this.myBoard.GameOver();
        } else if (this.enermyBoards.has(sender)) {
            this.enermyBoards.get(sender).GameOver();
        }
    }
}

class Board {

    cells = []; // the game board (Game.row x Game.column)
//    score = 0;

    nextBlocks = [];

    // html element name
    boardID = undefined;
    scoreID = undefined;
    nextBlockID = undefined;

    owner = "";
    state = "";

    isEermy = false;

    Start() {
        this.state = 'playing'
    }    

    GameOver() {
        this.state = 'over';
    }

    IsPlaying() {
        return this.state == 'playing'
    }    

    SetNextBlock(blocks) {

        this.nextBlocks = [];

        blocks.forEach((index) => {
             let block = Block.New(index);
             this.nextBlocks.push(block);
        });
    }

    SetCells(cells) {
        this.cells = cells;

        this.drawBoard();
    }

    // for me
    Init(boardID, scoreID, nextBlockID) {

        this.boardID = boardID;
        this.scoreID = scoreID;
        this.nextBlockID = nextBlockID;

        this.nextBlocks = [];

        for(let row=0;row< window.Game.row;row++)
        {
            let rowObject =[];

            for(let column=0; column< window.Game.column;column++)
            {
                rowObject[column] = window.Game.EMPTY ;// 0 mean empty cell, 1 mean cell occupy a block
            }
            this.cells[row] = rowObject;
        }

        this.initBoard();

        this.score = 0;
        $(scoreID).text(this.score);

        this.state = 'ready';        
    }

    // for enermy
    InitEnermy(boardID, scoreID) {

        this.boardID = boardID;
        this.scoreID = scoreID;
        this.isEermy = true;

        this.nextBlocks = [];

        for(let row=0;row< window.Game.row;row++)
        {
            let rowObject =[];

            for(let column=0; column< window.Game.column;column++)
            {
                rowObject[column] = window.Game.EMPTY ;// 0 mean empty cell, 1 mean cell occupy a block
            }
            this.cells[row] = rowObject;
        }

        this.initBoard();

        this.score = 0;
        $(scoreID).text(this.score);

        this.state = 'ready';        
    }


    
    initBoard ()
    {
        let html ="";

        let cl = "cell";
        if (this.isEermy) {
            cl = "enermy-cell";
        }

        // board
        for(let i=0;i<window.Game.row;i++)
        {
            html+="<tr>";
            for(let j=0;j<window.Game.column;j++)
            {
                html+="<td id='r"+i+"c"+j+"' class='"+cl+"'></td>";
            }
            html+="</tr>";
        }
        $(this.boardID).html(html);

        // next block
        if (!this.isEermy) {
            for (let i = 0; i < 3; i++) {
                html="";
                for(let r=0;r<4;r++)
                {
                    html+="<tr>";
                    for(let c=0;c<4;c++)
                    {
                        html+="<td id='p"+i+"r"+r+"c"+c+"' class='preview-cell'></td>";
                    }
                    html+="</tr>";
                }
                $(this.nextBlockID + i).html(html);
            }
        }
    }

}

class Block {

    Index = -1;

    //block is a two dimensional mtatrix of 4*4
    Shape =[];

    static New(index) {
        let newBlock = new Block();

        newBlock.Index = index;
        newBlock.Shape = window.Game.SHAPES[index-1];

        return newBlock;
    }
}

function CreateGame() {

    window.Game = new Game();
    window.Game.Init();

    $(document).keydown(function(e)
    {
        if(e.keyCode == 32) //space
        {
            PlaySound("block-down");
            send("block-drop", "");
            e.preventDefault();
        }
        if (e.keyCode == 38) //up
        {
            PlaySound("block-down");
            send("block-rotate", "");
            e.preventDefault();
        }
        if (e.keyCode == 37) //left
        {
            PlaySound("block-down");
            send("block-left", "");
            e.preventDefault();
        }
        if(e.keyCode == 39) //Right
        {
            PlaySound("block-down");
            send("block-right", "");
            e.preventDefault();
        }
        if(e.keyCode == 40) //Down
        {
            PlaySound("block-down");
            send("block-down", "");
            e.preventDefault();
        }
    });
/*
    $("#game").on('touchstart',function(e){

        if (!window.Game.myBoard.IsPlaying()) { return; }

        let touch = e.originalEvent.changedTouches[0];

        window.Game.touchX = touch.screenX;
        window.Game.touchY = touch.screenY;
    });

    $("#game").on('touchend',function(e){

        if (!window.Game.myBoard.IsPlaying()) { return; }

        let touch = e.originalEvent.changedTouches[0];

        x = touch.screenX;
        y = touch.screenY;

        if(window.Game.touchY-y>50){
            print("to top");
            e.originalEvent.preventDefault();
        } else if(y-window.Game.touchY>50){
            print("to bottom");
            e.originalEvent.preventDefault();     
        } else if(window.Game.touchX-x>50){
            print("to left");
            e.originalEvent.preventDefault();
        }else if(x-window.Game.touchX>50){
            print("to right");
            e.originalEvent.preventDefault();
        }
    });
*/
    $("#triangle-top").on('click', function(e) {
        PlaySound("block-down");
        send("block-rotate", "");
        e.preventDefault();
    });

    $("#triangle-left").on('click', function(e) {
        PlaySound("block-down");
        send("block-left", "");
        e.preventDefault();
    });

    $("#triangle-bottom").on('click', function(e) {
        PlaySound("block-down");
        send("block-drop", "");
        e.preventDefault();
    });

    $("#triangle-right").on('click', function(e) {
        PlaySound("block-down");
        send("block-right", "");
        e.preventDefault();
    });

    /*
    document.addEventListener('touchmove', function(e) {
        if (window.Game.myBoard.IsPlaying()) {
            let touch =  e.touches[0] || e.changedTouches[0];
            
            window.Game.touchStartX = touch.pageX;
            window.Game.touchStartY = touch.pageY;
        }
    }, false);

    document.addEventListener("touchend", function(e) {
        if (window.Game.myBoard.IsPlaying()) {
            let touch = e.touches[0] || e.changedTouches[0];
            
            let x = touch.pageX - window.Game.touchStartX;
            let y = touch.pageY - window.Game.touchStartY;

            print("touch.pageX:"+ touch.pageX + ", window.Game.touchStartX:" + window.Game.touchStartX);
            print("touch.pageY:"+ touch.pageY + ", window.Game.touchStartY:" + window.Game.touchStartY);
            print("x:" + x + ", y:" + y);

            if (Math.abs(x) > Math.abs(y)) {
                if (x > 0) {
                    PlaySound("block-down");
                    send("block-right", "");
                } else {
                    PlaySound("block-down");
                    send("block-left", "");
                }
            } else {
                if (y > 0) {
                    PlaySound("block-down");
                    send("block-drop", "");
                } else {
                    PlaySound("block-down");
                    send("block-rotate", "");
                }
            }
        }
    }, false);
    */
}

function DrawNextBlock(id, i, block) {

    id = id + i;

    for(let r=0; r<4; r++)
    {
        for(let c=0; c<4; c++)
        {
            let obj = $(id).find("#p"+i+"r"+r+"c"+c);

            if(block.Shape[r][c] != window.Game.EMPTY)
            {
                obj.removeAttr('class').addClass('preview-cell preview-block block' + block.Index);
            }
            else
            {
                obj.removeAttr('class').addClass('preview-cell');
            }
        }
    }
}

function inRect(row, col, block) {

    return (block.row<=row && row<block.row+4 && block.col<=col && col<block.col+4)
}

function DrawBoard (cells, block)
{
    let board = window.Game.myBoard;
    board.cells = cells;

    let html ="";
    for(let r=0;r<window.Game.row;r++)
    {
        html+="<tr>";
        for(let c=0;c<window.Game.column;c++)
        {
            let cl = "cell";
            if (board.cells[r][c] != window.Game.EMPTY) 
            {
                cl += " block block"+board.cells[r][c]
            } 
            if (inRect(r,c,block) && block.shape[r-block.row][c-block.col] != window.Game.EMPTY) {
                cl += " block block" + block.shapeIndex
            } 

            html+="<td id='r"+r+"c"+c+"' class='"+cl+"'></td>";
        }
        html+="</tr>";
    }
    $(board.boardID).html(html);
}

function DrawEnermyBoard (nick, cells, block)
{
    let board = window.Game.enermyBoards.get(nick);
    board.cells = cells;

    let html ="";
    for(let r=0;r<window.Game.row;r++)
    {
        html+="<tr>";
        for(let c=0;c<window.Game.column;c++)
        {
            let cl = "enermy-cell";
            if (board.cells[r][c] != window.Game.EMPTY) 
            {
                cl += " enermy-block block"+board.cells[r][c]
            } 
            if (inRect(r,c,block) && block.shape[r-block.row][c-block.col] != window.Game.EMPTY) {
                cl += " enermy-block block" + block.shapeIndex
            } 

            html+="<td id='r"+r+"c"+c+"' class='"+cl+"'></td>";
        }
        html+="</tr>";
    }
    $(board.boardID).html(html);
}


function DrawNextBlocks() {
    let board = window.Game.myBoard;
    DrawNextBlock(board.nextBlockID, 0, board.nextBlocks[0]);
    DrawNextBlock(board.nextBlockID, 1, board.nextBlocks[1]);
    DrawNextBlock(board.nextBlockID, 2, board.nextBlocks[2]);
}

function PlaySound(action) {
    let audio = new Audio("sound/" + action + ".wav");
    audio.play();
}