// 게임 화면 안에 시작 버튼 생성
var createStartButton = function() {

    let button = $('<button id="startButton" />')
                    .text('START')
                    .click(function() {
                        print("START")
                        print(window.Game)
                        if (window.Game == undefined) {
                            CreateGame();
                            ResetPlay();
                        }
                        ResetPlay();
                        StartPlay();

                        // delete button
                        $('#startButton').remove();
                    });
    $('#game').append(button);
}

var CreateGame = function() {
    window.Game = window.Game || {};
    Game.row= 20;
    Game.column = 20;
    Game.LEFT = "Left";
    Game.RIGHT = "Right";
    Game.UP = "Up";
    Game.DOWN = "Down";
    Game.EMPTY = 0;
    Game.FULL = 1;
    Game.MIDDLE = 9; // (column-2)/2
    Game.drawGameMap=function()
    {
        var html ="";
        for(var i=0;i<Game.row;i++)
        {
            html+="<tr>";
            for(var j=0;j<Game.column;j++)
            {
                html+="<td id='r"+i+"c"+j+"' class='cell'></td>";
            }

            html+="</tr>";
        }
        $("#gameMap").html(html);

        html="";
        for(var i=0;i<4;i++)
        {
            html+="<tr>";
            for(var j=0;j<4;j++)
            {
                html+="<td id='pr"+i+"pc"+j+"' class='cell'></td>";
            }

            html+="</tr>";
        }
        $("#gamePreviewMap").html(html);
    }
    Game.next =0;//Store preview block
    Game.Board = function()
    {
        this.cells =[];
        for(var row=0;row< Game.row;row++)
        {
            var rowObject =[];

            for(var column=0; column< Game.column;column++)
            {
                rowObject[column ] = Game.EMPTY ;// 0 mean empty cell, 1 mean cell occupy a block
            }
            this.cells[row] = rowObject;
        }
        this.clearGameBoard = function()
        {
            for(var row=0;row< Game.row;row++)
                {

                    for(var column=0; column< Game.column;column++)
                    {
                        $("#r"+row+"c"+column).removeClass("cell");
                        $("#r"+row+"c"+column).removeClass("block");
                        $("#r"+row+"c"+column).removeClass("animate");
                    }

                }

        };
        this.animateRow = function(row)
        {
            for(var column=0; column< Game.column;column++)
            {
                $("#r"+row+"c"+column).addClass("animate");

            }

        };
        this.drawGameBoard = function ()
        {
                this.clearGameBoard();
                ///alert("After clear");
                console.log("Redrawing gameboard");
                for(var row=0;row< Game.row;row++)
                {

                    for(var column=0; column< Game.column;column++)
                    {

                        var className;
                        if(this.cells[row][column] == Game.EMPTY)
                        {
                            className = "cell";
                        }
                        else
                        {
                            className = "block";
                        }
                        $("#r"+row+"c"+column).addClass(className);
                    }

                }
        };

    };
    Game.gameBoard  = new Game.Board();
    Game.Block=function()
    {
        this.currentRow = 0;
        this.currentColumn = Game.MIDDLE;//To start at the middle

        //block is a two dimensional mtatrix of 4*4
        this.blockCells =[];

        this.init=function()
        {
            for(var row=0;row< 4;row++)
            {
                var rowObject =[];

                for(var column=0; column< 4;column++)
                {
                    rowObject[column ] = 0 ;// 0 mean empty cell, 1 mean cell occupy a block
                }
                this.blockCells[row] = rowObject;
            }
            this.createARadomBlock();
            this.drawBlock();
        };
        this.createARadomBlock = function()
        {
            var random = parseInt(Math.random() * 8);
            switch(Game.next)
            {
                case 0:
                    this.blockCells = this.getShape1();
                break;
                case 1:
                    this.blockCells = this.getShape2();
                break;
                case 2:
                    this.blockCells = this.getShape3();
                break;
                case 3:
                    this.blockCells = this.getShape4();
                break;
                case 4:
                    this.blockCells = this.getShape5();
                break;
                case 5:
                    this.blockCells = this.getShape6();
                break;
                case 6:
                    this.blockCells = this.getShape7();
                break;
                default:
                    this.blockCells = this.getShape8();
            }
            Game.next = random;
            this.showPreview();
        };
        this.showPreview = function()
        {
            var blockCells;
            switch(Game.next)
            {
                case 0:
                    blockCells = this.getShape1();
                break;
                case 1:
                    blockCells = this.getShape2();
                break;
                case 2:
                    blockCells = this.getShape3();
                break;
                case 3:
                    blockCells = this.getShape4();
                break;
                case 4:
                    blockCells = this.getShape5();
                break;
                case 5:
                    blockCells = this.getShape6();
                break;
                case 6:
                    blockCells = this.getShape7();
                break;
                default:
                    blockCells = this.getShape8();
            }
            for(var r = 0;r< 4; r++)
            {
                for(var c=0;c< 4; c++)
                {
                    //console.log("r "+r+" c "+c + " "+ this.blockCells[r][c]);
                    //$("#pr"+y+"pc"+x).removeClass("block");
                    if(blockCells[r][c]==1)
                    {
                        var y =  r;
                        var x =  c;
                        $("#pr"+y+"pc"+x).addClass("block");
                        console.log("#pr"+y+"pc"+x);
                    }
                    else
                    {
                        var y =  r;
                        var x =  c;
                        $("#pr"+y+"pc"+x).removeClass("block");
                    }

                }
            }
            console.log("Show preview");
        };
        this.getShape1 = function()
        {
            var blockCells =[
                                        [1,1,1,0],
                                        [0,1,0,0],
                                        [0,0,0,0],
                                        [0,0,0,0],
                                    ];
            console.log("Shape 1");
            return blockCells;
        };
        this.getShape2 = function()
        {
            var blockCells =[
                                        [1,1,1,1],
                                        [0,0,0,0],
                                        [0,0,0,0],
                                        [0,0,0,0],
                                    ];
            console.log("Shape 2");
            return blockCells;
        };
        this.getShape3 = function()
        {
            var blockCells =[
                                        [1,1,0,0],
                                        [1,1,0,0],
                                        [0,0,0,0],
                                        [0,0,0,0],
                                    ];

            console.log("Shape 3");
            return blockCells;
        };
        this.getShape4 = function()
        {
            var blockCells =[
                                        [1,0,0,0],
                                        [1,1,0,0],
                                        [0,1,0,0],
                                        [0,0,0,0],
                                    ];
            console.log("Shape 4");
            return blockCells;
        };
        this.getShape5 = function()
        {
            var blockCells =[
                                        [1,0,0,0],
                                        [1,0,0,0],
                                        [1,1,0,0],
                                        [0,0,0,0],
                                    ];
            console.log("Shape 5");
            return blockCells;
        };
        this.getShape6 = function()
        {
            var blockCells =[
                                        [0,0,1,0],
                                        [0,0,1,0],
                                        [0,1,1,0],
                                        [0,0,0,0],
                                    ];
            console.log("Shape 6");
            return blockCells;
        };
        this.getShape7 = function()
        {
            var blockCells =[
                                        [0,1,0,0],
                                        [1,1,0,0],
                                        [1,0,0,0],
                                        [0,0,0,0],
                                    ];
            console.log("Shape 7");
            return blockCells;
        };
        this.getShape8 = function()
        {
            var blockCells =[
                                        [1,0,0,0],
                                        [1,0,0,0],
                                        [1,0,0,0],
                                        [1,0,0,0],
                                    ];
            console.log("Shape 8");
            return blockCells;
        };
        this.isOrigin = function()
        {
                if(this.currentRow == 0 /*&& this.currentColumn == Game.MIDDLE*/)
                {
                    return true;
                }
                else
                {
                    return false;
                }

        };
        /*
        * This method draw the block in the game board
        */
        this.drawBlock = function()
        {
            //current x ,current y is the first corner of block cell[0,0]
            for(var r = 0;r< 4; r++)
            {
                for(var c=0;c< 4; c++)
                {
                    //console.log("r "+r+" c "+c + " "+ this.blockCells[r][c]);
                    if(this.blockCells[r][c]==1)
                    {
                        var y = this.currentRow + r;
                        var x = this.currentColumn + c;
                        $("#r"+y+"c"+x).addClass("block");
                    }

                }
            }
        };
        this.isSafeToRotate = function()
        {
            var newBlock = [];
            for(var r = 0; r < 4 ;r++)
            {
                newBlock[r] = [];
                for(var c =0; c< 4;c++)
                {
                    newBlock[r][c] = 0;
                }
            }
            for(var r = 0; r < 4;r++)
            {
                for(var c =0; c< 4;c++)
                {
                    newBlock[c][r] = this.blockCells[r][3-c];
                }
            }
            // Temporarily rotate the block
            // used game board to check if it is ok to rotate;

            //Game.gameBoard.cells[]
            var ok = true;
            for(var r = 0; r < 4;r++)
            {
                for(var c =0; c< 4;c++)
                {
                    if( newBlock[r][c] == Game.FULL)
                    {
                        // Then game board must be empty
                        var y = this.currentRow + r;
                        var x = this.currentColumn + c;

                        if (Game.gameBoard.cells[ y ] [ x] != Game.EMPTY )
                        {
                            return false;
                        }
                    }
                }
            }
            return ok;
        };
        this.rotate = function()
        {
            if(!this.isSafeToRotate())
            {
                return;
            }
            this.clearOldDrawing();
            /*
                Mathmatically rotate
                change row to columns
            */
            var newBlock = [];
            for(var r = 0; r < 4 ;r++)
            {
                newBlock[r] = [];
                for(var c =0; c< 4;c++)
                {
                    newBlock[r][c] = 0;
                }
            }
            for(var r = 0; r < 4 ;r++)
            {
                for(var c =0; c< 4;c++)
                {
                    newBlock[c][r] = this.blockCells[r][3-c];
                }
            }

            this.blockCells = newBlock;

            this.drawBlock();
        };
        this.clearOldDrawing = function()
        {
                //current x ,current y is the first corner of block cell[0,0]
            for(var r = 0;r< 4; r++)
            {
                for(var c=0;c< 4; c++)
                {
                    //console.log("r "+r+" c "+c + " "+ this.blockCells[r][c]);
                    if(this.blockCells[r][c]==1)
                    {
                        var y = this.currentRow + r;
                        var x = this.currentColumn + c;
                        $("#r"+y+"c"+x).removeClass("block");
                    }

                }
            }
        };
        this.isSafeToMoveDown = function()
        {

            var ok = true;

            for(var r = 0; r < 4 ;r++)
            {
                for(var c =0; c< 4;c++)
                {

                    // Then game board must be empty

                    if( this.blockCells[r][c] == Game.FULL)
                    {
                        var y = this.currentRow + r +1 ;
                        var x = this.currentColumn + c;
                        if(y >= Game.row  )//last row, reach the end stop and return false ,cannot move further down
                        {
                            return false;
                        }
                        if (Game.gameBoard.cells[ y ] [ x] != Game.EMPTY )
                        {
                            return false;
                        }
                    }
                }
            }
            return ok;
        };
        this.moveDown = function()
        {
            if (this.isSafeToMoveDown())
            {
                //Move down

                this.clearOldDrawing();
                this.currentRow ++;
                this.drawBlock();
            }
            else  // Set the game board cell
            {
                Game.nextTern();
            }
        };
        this.storeGameBoardData = function()
        {
            for(var r = 0; r < 4 ;r++)
                {
                    for(var c =0; c< 4;c++)
                    {

                        var y = this.currentRow + r  ;
                        var x = this.currentColumn + c;
                        if( this.blockCells[r][c] == Game.FULL)
                        {
                            Game.gameBoard.cells[ y ] [ x] = Game.FULL;
                            console.log("Set row "+y +" column "+x +" to 1");
                        }
                    }
                }
        };
        this.processGameRow = function()
        {
            //start from the last row
            var rowIndexToRemove= [];
            for(var last = Game.row-1 ; last >=0; last--)
            {
                var ok  = true;
                for(var col =0 ;col < Game.column ;col++)
                {
                    ok = ok && Game.gameBoard.cells[last][col] == Game.FULL;

                }

                if (ok) // This row is full
                {
                    console.log("Checking row "+last+" full "+ok);
                    rowIndexToRemove.unshift(last);
                }
            }
            // For each remove row shit the row from top
            for(var lastIndex = 0 ; lastIndex < rowIndexToRemove.length;  lastIndex ++)
            {
                var rowIndex = rowIndexToRemove[ lastIndex ];
                var animateRow = rowIndex;
                //shift the one row down
                console.log("Shifting down row "+rowIndex);
                for(var c =0; c < Game.column; c++)
                {
                    console.log("remove row "+rowIndex +" column "+c +" with row "+(rowIndex-1)+" column "+c);
                    Game.gameBoard.cells[rowIndex][c] = Game.gameBoard.cells[rowIndex-1][c];

                }

                console.log("Congrulation");
                rowIndex --;
                while(rowIndex > 0 )
                {
                    for(var c =0; c < Game.column; c++)
                    {
                        Game.gameBoard.cells[rowIndex][c] = Game.gameBoard.cells[rowIndex-1][c];

                    }
                    rowIndex --;
                }
                for(var col =0; col < Game.column; col++)
                {
                    //Add the empty row at top
                    Game.gameBoard.cells[0][col] = Game.EMPTY;
                }
                Game.gameBoard.animateRow(animateRow);
                setTimeout(function()
                {
                    Game.gameBoard.drawGameBoard();
                },100);
                Game.score += 1000;
            }
            //Add bonus if more than one row
            if(rowIndexToRemove.length >1)
            {
                Game.score += (rowIndexToRemove.length -1)*500;
            }

            Game.displayScore();
        };

        this.isSafeToMoveLeft = function()
        {

            var ok = true;
            for(var r = 0; r < 4 ;r++)
            {
                for(var c =0; c< 4;c++)
                {

                    // Then game board must be empty
                    var y = this.currentRow + r  ;
                    var x = this.currentColumn + c -1;
                    if( this.blockCells[r][c] == Game.FULL)
                    {
                        if(x < 0)
                        {
                                return false;
                        }
                        if (Game.gameBoard.cells[ y ] [ x] != Game.EMPTY )
                        {
                            return false;
                        }
                    }
                }
            }
            return ok;
        };
        this.moveLeft = function()
        {
            if(this.isSafeToMoveLeft())
            {
                this.clearOldDrawing();
                this.currentColumn --;
                this.drawBlock();

            }
        };

        this.isSafeToMoveRight = function()
        {

            var ok = true;
            for(var r = 0; r < 4 ;r++)
            {
                for(var c =0; c< 4;c++)
                {

                    // Then game board must be empty
                    var y = this.currentRow + r  ;
                    var x = this.currentColumn + c +1;
                    if( this.blockCells[r][c] == Game.FULL)
                    {
                        if(this.x + 2 >= Game.column)
                        {
                            return false;
                        }
                        if (Game.gameBoard.cells[ y ] [ x] != Game.EMPTY )
                        {
                            return false;
                        }
                    }
                }
            }
            return ok;
        };
        this.moveRight = function()
        {
            if(this.isSafeToMoveRight())
            {
                this.clearOldDrawing();
                this.currentColumn ++;
                this.drawBlock();

            }
        };
    };

    Game.displayScore = function()
    {
        $("#gameScore").text(Game.score);
    };

    Game.gameOver = function()
    {
        alert("Game over, please refresh the page to start new game");
        clearInterval(timer);
        createStartButton();
    }

    Game.nextTern = function()
    {
        Game.current.storeGameBoardData();
        Game.current.processGameRow();
        Game.current  = new Game.Block ();
        Game.current.init();
    }

    $(document).keydown(function(e)
    {
        try
        {
            if(e.keyCode == 32) //space
            {
                // to bottom
                while(Game.current.isSafeToMoveDown())
                {
                    Game.current.moveDown();
                }

                if (! Game.current.isOrigin())
                {
                    Game.nextTern();
                }
                else
                {
                    Game.gameOver();
                }

            }
            if (e.keyCode == 38) //up
            {
                Game.current.rotate();
            }
            if (e.keyCode == 37) //left
            {
                Game.current.moveLeft();
            }
            if(e.keyCode == 39) //Right
            {
                Game.current.moveRight();
            }
            if(e.keyCode == 40) //Down
            {
                if(Game.current.isSafeToMoveDown())
                {
                    Game.current.moveDown();
                }
                else if (! Game.current.isOrigin())
                {
                    Game.nextTern();
                }
                else
                {
                    Game.gameOver();
                }

            }
        }
        catch(e)
        {
            print(e);
            Game.gameOver();
        }
    });
}

var ResetPlay = function() {
    clearInterval(window.timer);

    Game.drawGameMap();
    Game.current  = new Game.Block ();
    Game.current.init();
    Game.score = 0;
}

var StartPlay = function() {

    window.timer = window.setInterval(function()
    {
        if(Game.current.isSafeToMoveDown())
        {
            Game.current.moveDown();
        }
        else if (! Game.current.isOrigin())
        {
            Game.nextTern();
        }
        else
        {
            Game.gameOver();
        }
    },1000);
}