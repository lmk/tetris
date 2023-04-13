package main

var (
	SHAPES = [...][BLOCK_ROW][BLOCK_COLUMN]int{
		{
			{0, 0, 0, 0},
			{1, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		}, {
			{0, 0, 0, 0},
			{1, 1, 1, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		}, {
			{0, 0, 0, 0},
			{0, 1, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
		}, {
			{0, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 0},
		}, {
			{0, 1, 0, 0},
			{0, 1, 0, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
		}, {
			{0, 0, 1, 0},
			{0, 0, 1, 0},
			{0, 1, 1, 0},
			{0, 0, 0, 0},
		}, {
			{0, 0, 1, 0},
			{0, 1, 1, 0},
			{0, 1, 0, 0},
			{0, 0, 0, 0},
		},
	}
)

type Block struct {
	row        int
	col        int
	shape      [BLOCK_ROW][BLOCK_COLUMN]int
	shapeIndex int
}

func (b *Block) cloneShape(shape [BLOCK_ROW][BLOCK_COLUMN]int) {
	for i := 0; i < BLOCK_ROW; i++ {
		for j := 0; j < BLOCK_COLUMN; j++ {
			b.shape[i][j] = shape[i][j]
		}
	}
}

func (b *Block) Clone(from Block) {
	b.row = from.row
	b.col = from.col
	b.shapeIndex = from.shapeIndex
	b.cloneShape(from.shape)
}

func NewBlock(shapeIndex int) Block {
	block := Block{
		row:        0,
		col:        BOARD_CENTER,
		shapeIndex: shapeIndex,
	}
	block.SetShape(shapeIndex)

	return block
}

func (b *Block) SetShape(shapeIndex int) {
	b.shapeIndex = shapeIndex
	b.cloneShape(SHAPES[shapeIndex])
}

func (b *Block) Rotate() {
	var newShape [BLOCK_ROW][BLOCK_COLUMN]int
	for i := 0; i < BLOCK_ROW; i++ {
		for j := 0; j < BLOCK_COLUMN; j++ {
			newShape[i][j] = b.shape[BLOCK_COLUMN-j-1][i]
		}
	}
	b.cloneShape(newShape)
}
