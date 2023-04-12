// 게임 화면 안에 시작 버튼 생성
function createStartButton() {
    let button = $('<button id="startButton" />')
                    .text('START')
                    .click(function() {
                        if (window.Game == undefined) {
                            CreateGame();
                        }
                        StartPlay();

                        // delete button
                        $('#startButton').remove();
                    });
    $('#game').append(button);
}

class Game {

    row = 15;
    column = 10;
    MIDDLE = parseInt((this.column-3)/2);
    EMPTY = 0;
    FULL = 1;
    CURRENT = 2;    // current moving block

    myBoard = undefined;

    Init() {

        if (this.myBoard == undefined) {
            this.myBoard = new Board();
        }

        this.myBoard.Init("#my-board", "#my-score", "#my-nextBlock");
    }

    GameOver()
    {
        alert("Game over, please refresh the page to start new game");
        clearInterval(timer);
        createStartButton();
    }
}

class Board {

    cells = []; // the game board (Game.row x Game.column)
    score = 0;

    nextBlocks = [];
    currentBlock = undefined;

    currentRow = 0;
    currentColumn = window.Game.MIDDLE; //To start at the middle

    // html element name
    boardID = undefined;
    scoreID = undefined;
    nextBlockID = undefined;

    AppendNextBlock(blocks) {

        blocks.forEach((index) => {
             let block = Block.New(index);
             this.nextBlocks.push(block);
        });
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


    isSafeToRotateBlock() {
        let rotateBlock = Block.Clone(this.currentBlock);
        rotateBlock.RotateCell();

        for(let r=0; r<4; r++)
        {
            for(let c=0; c<4; c++)
            {
                if(rotateBlock.blockCells[r][c] == window.Game.FULL)
                {
                    let y = this.currentRow + r;
                    let x = this.currentColumn + c;
                    if(y<0 || y>=window.Game.row || x<0 || x>=window.Game.column)
                    {
                        return false;
                    }
                    if(this.cells[y][x] == window.Game.FULL)
                    {
                        return false;
                    }
                }
            }
        }

        return true;
    }

    isSafeNewBlock() {

        let newBlock = this.nextBlocks[0];

        for(let r=0; r<4; r++)
        {
            for(let c=0; c<4; c++)
            {
                if(newBlock.blockCells[r][c] == window.Game.FULL)
                {
                    let y = this.currentRow + r;
                    let x = this.currentColumn + c;
                    if(y<0 || y>=window.Game.row || x<0 || x>=window.Game.column)
                    {
                        return false;
                    }
                    if(this.cells[y][x] == window.Game.FULL)
                    {
                        return false;
                    }
                }
            }
        }
        return true;
    }


    isSafeToMoveDownBlock() {
        for(let r=0; r<4; r++)
        {
            for(let c=0; c<4; c++)
            {
                if(this.currentBlock.blockCells[r][c] == window.Game.FULL)
                {
                    let y = this.currentRow + r +1;
                    let x = this.currentColumn + c;
                    if(y<0 || y>=window.Game.row || x<0 || x>=window.Game.column)
                    {
                        return false;
                    }
                    if(this.cells[y][x] == window.Game.FULL)
                    {
                        return false;
                    }
                }
            }
        }
        return true;
    }

    isSafeToMoveLeftBlock() {
        for(let r=0; r<4; r++)
        {
            for(let c=0; c<4; c++)
            {
                if(this.currentBlock.blockCells[r][c] == window.Game.FULL)
                {
                    let y = this.currentRow + r;
                    let x = this.currentColumn + c -1;
                    if(y<0 || y>=window.Game.row || x<0 || x>=window.Game.column)
                    {
                        return false;
                    }
                    if(this.cells[y][x] == window.Game.FULL)
                    {
                        return false;
                    }
                }
            }
        }
        return true;
    }

    isSafeToMoveRightBlock() {
        for(let r=0; r<4; r++)
        {
            for(let c=0; c<4; c++)
            {
                if(this.currentBlock.blockCells[r][c] == window.Game.FULL)
                {
                    let y = this.currentRow + r;
                    let x = this.currentColumn + c +1;
                    if(y<0 || y>=window.Game.row || x<0 || x>=window.Game.column)
                    {
                        return false;
                    }
                    if(this.cells[y][x] == window.Game.FULL)
                    {
                        return false;
                    }
                }
            }
        }
        return true;
    }


    RotateBlock() {
        if(! this.isSafeToRotateBlock())
        {
            return;
        }

        this.removeCurrentBlock();
        this.currentBlock.RotateCell();
        this.saveCurrentBlock();
    }

    MoveDownBlock() {
        if(! this.isSafeToMoveDownBlock())
        {
            return false;
        } else {
            this.removeCurrentBlock();
            this.currentRow ++;
            this.saveCurrentBlock();
        }
        return true;
    }

    MoveLeftBlock() {
        if(! this.isSafeToMoveLeftBlock())
        {
            return;
        }

        this.removeCurrentBlock();
        this.currentColumn --;
        this.saveCurrentBlock();
    }

    MoveRightBlock() {
        if(! this.isSafeToMoveRightBlock())
        {
            return;
        }

        this.removeCurrentBlock();
        this.currentColumn ++;
        this.saveCurrentBlock();
    }

    MoveBottomBlock() {
        while(this.isSafeToMoveDownBlock())
        {
            if (this.MoveDownBlock() == false)
            {
                return false;
            }
        }

        this.NextTern();

        return true;
    }

    saveCurrentBlock() {
        for(let r=0; r<4; r++)
        {
            for(let c =0; c<4; c++)
            {
                if( this.currentBlock.blockCells[r][c] == window.Game.FULL)
                {
                    let y = this.currentRow + r;
                    let x = this.currentColumn + c;

                    this.cells[y][x] = window.Game.CURRENT;
                }
            }
        }

        let row = this.currentRow;
        let column = this.currentColumn;

        DrawBlock(this.boardID, this.currentBlock, row, column);
    }

    removeCurrentBlock() {
        for(let r=0; r<4; r++)
        {
            for(let c =0; c<4; c++)
            {
                if( this.currentBlock.blockCells[r][c] == window.Game.FULL)
                {
                    let y = this.currentRow + r;
                    let x = this.currentColumn + c;

                    this.cells[y][x] = window.Game.EMPTY;
                }
            }
        }

        let row = this.currentRow;
        let column = this.currentColumn;

        EraseBlock(this.boardID, this.currentBlock, row, column);
    }

    addScore(score) {
        this.score += score;
        $(this.scoreID).text(this.score);
    }

    processFullRow() {

        let countFull = 0;

        for(let row=window.Game.row-1; row>=0; row--)
        {
            let isFull = true;
            for(let column=0; column< window.Game.column; column++)
            {
                if(this.cells[row][column] == window.Game.EMPTY)
                {
                    isFull = false;
                    break;
                }
            }

            if(isFull)
            {
                // hide row - element
                for(let c=0; c< window.Game.column;c++)
                {
                    $(this.boardID).find("#r"+row+"c"+c).addClass("bingo")
                }
        
                // move row - cells
                for(let r=row; r>0; r--)
                {
                    for(let c=0; c< window.Game.column; c++)
                    {
                        this.cells[r][c] = this.cells[r-1][c];
                    }
                }

                // move row - element attribute
                for(let r=row; r>0; r--)
                {
                    for(let c=0; c< window.Game.column; c++)
                    {
                        let hClass = $(this.boardID).find("#r"+(r-1)+"c"+c).attr("class");
                        $(this.boardID).find("#r"+r+"c"+c).removeAttr("class").addClass(hClass);
                    }
                }

                // clear first row - element
                for(let c=0; c< window.Game.column; c++)
                {
                    $(this.boardID).find("#r0c"+c).removeAttr("class").addClass("cell");
                }                

                countFull++;
                row++;

                this.addScore(10 * countFull);
            }
        }
    }

    createNextBlock() {
        if (this.nextBlocks.length < 10) {
            send("new-block", (10 - this.nextBlocks.length).toString());
        }
    }

    currentToFull() {
        for(let r=0; r<4; r++)
        {
            for(let c=0; c<4; c++)
            {
                if(this.currentBlock.blockCells[r][c] == window.Game.FULL)
                {
                    let y = this.currentRow + r;
                    let x = this.currentColumn + c;
                    this.cells[y][x] = window.Game.FULL;
                }
            }
        }
    }


    ChangeCurrnetBlock() {
        this.currentBlock = this.nextBlocks.shift();
        DrawNextBlock(this.nextBlockID + "0", this.nextBlocks[0]);
        DrawNextBlock(this.nextBlockID + "1", this.nextBlocks[1]);
        DrawNextBlock(this.nextBlockID + "2", this.nextBlocks[2]);
    }

    NextTern() {
        this.currentToFull();
        this.processFullRow();
   
        this.currentRow = 0;
        this.currentColumn = window.Game.MIDDLE;
        this.ChangeCurrnetBlock();

        this.createNextBlock();
        if (this.isSafeNewBlock() == false)
        {
            return false;
        }


        this.saveCurrentBlock();

        return true;
    }

    Start() {
        this.createNextBlock();
        let timer = window.setInterval(function()
        {
            clearTimeout(timer);
            window.Game.myBoard.ChangeCurrnetBlock();

        },100);        
    }

}

class Block {

    Index = -1;

    //block is a two dimensional mtatrix of 4*4
    blockCells =[];

    Shape = [
        [
            [0,0,0,0],
            [1,1,1,0],
            [0,1,0,0],
            [0,0,0,0],
        ], [
            [0,0,0,0],
            [1,1,1,1],
            [0,0,0,0],
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
        ], [
            [0,1,0,0],
            [0,1,0,0],
            [0,1,0,0],
            [0,1,0,0]
        ],
    ];

    static New(index) {
        let newBlock = new Block();

        newBlock.Index = index;
        newBlock.blockCells = newBlock.Shape[index];

        return newBlock;
    }

    static Clone(block) {
        let newBlock = new Block();

        newBlock.Index = block.Index;

        for(let r=0; r<4; r++)
        {
            newBlock.blockCells[r] = [];
            for(let c=0; c<4; c++)
            {
                newBlock.blockCells[r][c] = block.blockCells[r][c];
            }
        }

        return newBlock;
    }

    NewEmptyCell() {
        let cell = [];
        for(let r=0; r<4; r++)
        {
            cell[r] = [];
            for(let c=0; c<4; c++)
            {
                cell[r][c] = window.Game.EMPTY;
            }
        }
        return cell;
    }

    RotateCell() {
        let rotateCell = this.NewEmptyCell();
        for(let r=0; r<4; r++)
        {
            for(let c=0; c<4; c++)
            {
                rotateCell[c][r] = this.blockCells[r][3-c];
            }
        }
        this.blockCells = rotateCell;
    }
}

function CreateGame() {

    window.Game = new Game();
    window.Game.Init();

    $(document).keydown(function(e)
    {
        try
        {
            if(e.keyCode == 32) //space
            {
                if (window.Game.myBoard.MoveBottomBlock() == false) {
                    send("over-game", "to-bottom")
                }
                e.preventDefault();
            }
            if (e.keyCode == 38) //up
            {
                window.Game.myBoard.RotateBlock();
                e.preventDefault();
            }
            if (e.keyCode == 37) //left
            {
                window.Game.myBoard.MoveLeftBlock();
                e.preventDefault();
            }
            if(e.keyCode == 39) //Right
            {
                window.Game.myBoard.MoveRightBlock();
                e.preventDefault();
            }
            if(e.keyCode == 40) //Down
            {
                if (window.Game.myBoard.MoveDownBlock() == false) {
                    send("over-game", "to-down")
                }
                e.preventDefault();
            }
        }
        catch(e)
        {
            print(e);
            send("over-game", "error: " + e)
        }
    });
}

function StartPlay() {
    send("start-game", "");
}

function DrawBlock(id, block, currentRow, currentColumn) {
    for(let r=0; r<4; r++)
    {
        for(let c=0; c<4; c++)
        {
            if(block.blockCells[r][c] == window.Game.FULL)
            {
                let y = currentRow + r;
                let x = currentColumn + c;
                $(id).find("#r"+y+"c"+x).addClass("block block" + block.Index);
            }
        }
    }
}

function EraseBlock(id, block, currentRow, currentColumn) {
    for(let r=0; r<4; r++)
    {
        for(let c=0; c<4; c++)
        {
            if(block.blockCells[r][c] == window.Game.FULL)
            {
                let y = currentRow + r;
                let x = currentColumn + c;
                $(id).find("#r"+y+"c"+x).removeAttr('class').addClass('cell');
            }
        }
    }
}

function DrawNextBlock(id, block) {
    for(let r=0; r<4; r++)
    {
        for(let c=0; c<4; c++)
        {
            let obj = $(id).find("#pr"+r+"pc"+c);

            if(block.blockCells[r][c] == window.Game.FULL)
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
