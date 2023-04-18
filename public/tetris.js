// 게임 화면 안에 시작 버튼 생성
function createStartButton() {
    let button = $('<button id="startButton" />')
                    .text('START')
                    .click(function() {
                        send("start-game", "");

                        // delete button
                        $('#startButton').remove();
                    });
    $('#game').append(button);
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
        $('#nick').val(owner);
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

    Start() {
        this.state = 'start-game'
    }    

    GameOver() {
        this.state = 'game-over';
    }

    IsPlaying() {
        return this.state == 'start-game'
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

    Init(boardID, scoreID, nextBlockID) {

        this.boardID = boardID;
        this.scoreID = scoreID;
        this.nextBlockID = nextBlockID;

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

        this.state = 'ready-game';        
    }

    
    initBoard ()
    {
        let html ="";
        for(let i=0;i<window.Game.row;i++)
        {
            html+="<tr>";
            for(let j=0;j<window.Game.column;j++)
            {
                html+="<td id='r"+i+"c"+j+"' class='cell'></td>";
            }
            html+="</tr>";
        }
        $(this.boardID).html(html);

        html="";
        for(let i=0;i<4;i++)
        {
            html+="<tr>";
            for(let j=0;j<4;j++)
            {
                html+="<td id='pr"+i+"pc"+j+"' class='preview-cell'></td>";
            }
            html+="</tr>";
        }
        $(this.nextBlockID + "0").html(html);
        $(this.nextBlockID + "1").html(html);
        $(this.nextBlockID + "2").html(html);
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
            send("block-drop", "")
            e.preventDefault();
        }
        if (e.keyCode == 38) //up
        {
            send("block-rotate", "")
            e.preventDefault();
        }
        if (e.keyCode == 37) //left
        {
            send("block-left", "")
            e.preventDefault();
        }
        if(e.keyCode == 39) //Right
        {
            send("block-right", "")
            e.preventDefault();
        }
        if(e.keyCode == 40) //Down
        {
            send("block-down", "")
            e.preventDefault();
        }
    });
}

function DrawNextBlock(id, block) {
    for(let r=0; r<4; r++)
    {
        for(let c=0; c<4; c++)
        {
            let obj = $(id).find("#pr"+r+"pc"+c);

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

function DrawNextBlocks() {
    let board = window.Game.myBoard;
    DrawNextBlock(board.nextBlockID + "0", board.nextBlocks[0]);
    DrawNextBlock(board.nextBlockID + "1", board.nextBlocks[1]);
    DrawNextBlock(board.nextBlockID + "2", board.nextBlocks[2]);
}