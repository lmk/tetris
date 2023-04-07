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

    row = 20;
    column = 20;
    MIDDLE = 9;     // (column-2)/2
    EMPTY = 0;
    FULL = 1;
    CURRENT = 2;    // current moving block

    gameBoard = undefined;

    Init() {

        this.gameBoard = new Board();
        this.gameBoard.Init();

        this.drawGameMap();
    }

    drawGameMap ()
    {
        let html ="";
        for(let i=0;i<this.row;i++)
        {
            html+="<tr>";
            for(let j=0;j<this.column;j++)
            {
                html+="<td id='r"+i+"c"+j+"' class='cell'></td>";
            }
            html+="</tr>";
        }
        $("#board").html(html);

        html="";
        for(let i=0;i<4;i++)
        {
            html+="<tr>";
            for(let j=0;j<4;j++)
            {
                html+="<td id='pr"+i+"pc"+j+"' class='cell'></td>";
            }
            html+="</tr>";
        }
        $("#nextBlock").html(html);
    }

    Start() {
        this.gameBoard.Start();
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

    nextBlock = undefined;
    currentBlock = undefined;

    currentRow = 0;
    currentColumn = window.Game.MIDDLE; //To start at the middle

    Init() {
        for(let row=0;row< window.Game.row;row++)
        {
            let rowObject =[];

            for(let column=0; column< window.Game.column;column++)
            {
                rowObject[column] = window.Game.EMPTY ;// 0 mean empty cell, 1 mean cell occupy a block
            }
            this.cells[row] = rowObject;
        }
    }

    clearGameBoard() {
        for(let row=0;row< window.Game.row;row++)
        {
            for(let column=0; column< window.Game.column;column++)
            {
                $("#r"+row+"c"+column).removeClass("cell");
                $("#r"+row+"c"+column).removeClass("block");
                $("#r"+row+"c"+column).removeClass("animate");
            }
        }
    }

    animateRow(row) {
        for(let column=0; column< window.Game.column;column++)
        {
            $("#r"+row+"c"+column).addClass("animate");
        }
    }

    DrawGameBoard() {
        for(let row=0;row< window.Game.row;row++)
        {
            for(let column=0; column< window.Game.column;column++)
            {
                let className;
                if(this.cells[row][column] == window.Game.EMPTY)
                {
                    className = "cell";
                }
                else
                {
                    className = "block";
                }
                $("#r"+row+"c"+column).removeClass("cell");
                $("#r"+row+"c"+column).removeClass("block");
                $("#r"+row+"c"+column).removeClass("animate");
                $("#r"+row+"c"+column).addClass(className);
            }
        }
    }

    isSafeToRotateBlock() {
        let rotateBlock = this.currentBlock;
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

        this.currentBlock.Draw(this.currentRow, this.currentColumn);
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

        this.currentBlock.Erase(this.currentRow, this.currentColumn);
    }

    addScore(score) {
        this.score += score;
        $("#score").text(this.score);
    }

    processFullRow() {

        for(let row=window.Game.row-1; row>=0; row--)
        {
            let isFull = true;
            let countFull = 0;
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
                for(let r=row; r>0; r--)
                {
                    for(let c=0; c< window.Game.column; c++)
                    {
                        this.cells[r][c] = this.cells[r-1][c];
                    }
                }
                this.animateRow(row);
                countFull++;
                row++;

                addScore(10 * (countFull*countFull));
            }
        }
    }

    createNextBlock() {
        this.nextBlock = Block.Create(Math.floor(Math.random()*8));
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


    NextTern() {
        print("NextTern");
        this.saveCurrentBlock();
        this.currentToFull();
        this.processFullRow();
        this.currentBlock = this.nextBlock;
        this.createNextBlock();
    }

    Start() {
        this.createNextBlock();
        this.currentBlock = this.nextBlock;
        this.createNextBlock();
    }

}

class Block {

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

    static Create(index) {
        let newBlock = new Block();

        newBlock.set(index);
        newBlock.drawNext();

        return newBlock;
    }

    set(index) {
        this.blockCells = this.Shape[index];
    }

    drawNext() {
        if (this.blockCells == undefined) {
            print("blockCells is undefined");
        }

        for(let r=0; r<4; r++)
        {
            for(let c=0; c<4; c++)
            {
                if(this.blockCells[r][c] == window.Game.FULL)
                {
                    $("#pr"+r+"pc"+c).addClass("block");
                }
                else
                {
                    $("#pr"+r+"pc"+c).removeClass("block");
                }
            }
        }
        print("drawNext")
        print(this.blockCells)

    }

    Draw(currentRow, currentColumn) {
        for(let r=0; r<4; r++)
        {
            for(let c=0; c<4; c++)
            {
                if(this.blockCells[r][c] == window.Game.FULL)
                {
                    let y = currentRow + r;
                    let x = currentColumn + c;
                    $("#r"+y+"c"+x).addClass("block");
                }
            }
        }
    }

    Erase(currentRow, currentColumn) {
        for(let r=0; r<4; r++)
        {
            for(let c=0; c<4; c++)
            {
                if(this.blockCells[r][c] == window.Game.EMPTY)
                {
                    let y = currentRow + r;
                    let x = currentColumn + c;
                    $("#r"+y+"c"+x).removeClass("block");
                }
            }
        }
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
                if (window.Game.gameBoard.MoveBottomBlock() == false) {
                    window.Game.GameOver();
                }
            }
            if (e.keyCode == 38) //up
            {
                window.Game.gameBoard.RotateBlock();
            }
            if (e.keyCode == 37) //left
            {
                window.Game.gameBoard.MoveLeftBlock();
            }
            if(e.keyCode == 39) //Right
            {
                window.Game.gameBoard.MoveRightBlock();
            }
            if(e.keyCode == 40) //Down
            {
                if (window.Game.gameBoard.MoveDownBlock() == false) {
                    window.Game.GameOver();
                }
            }
        }
        catch(e)
        {
            print(e);
            window.Game.GameOver();
        }
    });
}

function StartPlay() {
    window.Game.Start();

    window.timer = window.setInterval(function()
    {
        if (window.Game.gameBoard.MoveDownBlock() == false) {
            window.Game.GameOver();
        }
    },1000);
}