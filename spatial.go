package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*===============================================================
 * Functions to manipulate a "field" of cells --- the main data
 * that must be managed by this program.
 *==============================================================*/

// The data stored in a single cell of a field
type Cell struct {
	kind  string
	score float64
	pre_kind string
}

// createField should create a new field of the ysize rows and xsize columns,
// so that field[r][c] gives the Cell at position (r,c).
func createField(rsize, csize int) [][]Cell {
	f := make([][]Cell, rsize)
	for i := range f {
		f[i] = make([]Cell, csize)
	}
	return f
}

// inField returns true iff (row,col) is a valid cell in the field
func inField(field [][]Cell, row, col int) bool {
	return row >= 0 && row < len(field) && col >= 0 && col < len(field[0])
}

// readFieldFromFile should open the given file and read the initial
// values for the field. The first line of the file will contain
// two space-separated integers saying how many rows and columns
// the field should have:
//    10 15
// each subsequent line will consist of a string of Cs and Ds, which
// are the initial strategies for the cells:
//    CCCCCCDDDCCCCCC
//
// If there is ever an error reading, this function should cause the
// program to quit immediately.
func readFieldFromFile(filename string) [][]Cell {

    //read from the input file
	in, err := os.Open(filename)
	if err != nil{
		fmt.Println("Error: couldn't open the file")
		os.Exit(3)
	}
	//store them into slices
	var lines []string = make([]string, 0)
	scanner := bufio.NewScanner(in)
	for scanner.Scan(){
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		fmt.Println("Sorry: there was some kind of error during the file reading")
	}

    in.Close()

    //get row and column number from first line
    //space split
    var items []string = strings.Split(lines[0]," ")
    row, err1 := strconv.Atoi(items[0])
	if err1 != nil || row < 0 {
		fmt.Println("Error: bad input row, should be a positivie integer")
		return nil
	}
    col, err2 := strconv.Atoi(items[1])
    if err2 != nil || col < 0 {
		fmt.Println("Error: bad input col, should be a positivie integer")
		return nil
	}

	//store the input file Cs and Ds into the nested slices 
	var c [][]Cell = createField(row, col)
	for i := 0; i < row; i++{
		for j := 0; j < col; j++{
			c[i][j].kind = string(lines[i+1][j])
		}
	}
	return c
}
/*
 * compare two strings
 * just to make it neat and simple in drawfield function.
 */
func compare(a, b string) bool{
	return strings.EqualFold(a, b)
}

// drawField should draw a representation of the field on a canvas and save the
// canvas to a PNG file with a name given by the parameter filename.  Each cell
// in the field should be a 5-by-5 square, and cells of the "D" kind should be
// drawn red and cells of the "C" kind should be drawn blue.
func drawField(field [][]Cell, filename string) {
    //create a canvas with 5 times of the row and col
	pic := CreateNewCanvas(len(field) * 5, len(field[0]) * 5)
	pic.SetLineWidth(1)
	for i := 0; i < len(field); i++{
		for j := 0; j < len(field[0]); j++{
			//pattern match between kind and color
			if compare(field[i][j].kind, "C") && compare(field[i][j].pre_kind, "C"){
				pic.SetFillColor(MakeColor(0,0,255))//blue
				drawSquare(pic, i, j)
			}
			if compare(field[i][j].kind, "D") && compare(field[i][j].pre_kind, "D"){
				pic.SetFillColor(MakeColor(255,0,0))//red
				drawSquare(pic, i, j)
			}
			if compare(field[i][j].kind, "D") && compare(field[i][j].pre_kind, "C"){
				pic.SetFillColor(MakeColor(255,255,0))//yellow
				drawSquare(pic, i, j)
			}
			if compare(field[i][j].kind, "C") && compare(field[i][j].pre_kind, "D"){
				pic.SetFillColor(MakeColor(0,255,0))//green
				drawSquare(pic, i, j)
			}
		}
	}
	pic.SaveToPNG(filename)
}

//draw a single square at (r,c) on the board b of width of 5 pixel
func drawSquare(b Canvas, r, c int) {
    x1, y1 := float64(r * 5), float64(c * 5)
    x2, y2 := float64((r+1) * 5), float64((c+1) * 5)
    b.MoveTo(x1, y1)
    b.LineTo(x1, y2)
    b.LineTo(x2, y2)
    b.LineTo(x2, y1)
    b.LineTo(x1, y1)
    //instead of FillStroke()
    b.Fill()
}

/*===============================================================
 * Functions to simulate the spatial games
 *==============================================================*/

// play a game between a cell of type "me" and a cell of type "them" (both me
// and them should be either "C" or "D"). This returns the reward that "me"
// gets when playing against them.
func gameBetween(me, them string, b float64) float64 {
	if me == "C" && them == "C" {
		return 1
	} else if me == "C" && them == "D" {
		return 0
	} else if me == "D" && them == "C" {
		return b
	} else if me == "D" && them == "D" {
		return 0
	} else {
		fmt.Println("type ==", me, them)
		panic("This shouldn't happen")
	}
}
/*
 * game is the "game" process from the current cell, plays the Prinsoner's 
 * dilema game with every legal neighbor. Uses inField function to eliminate 
 * the illegal neighbors of the boundary cells. Uses gameBetween to play the
 * single game. Returns the cell's score after all th games.
 */
func game(field [][]Cell, i, j int, b float64) float64{
	var res float64 = 0
	for x := i - 1; x <= i + 1; x++{
		for y:= j - 1; y <= j + 1; y++{
			if inField(field, x, y){
				res += gameBetween(field[i][j].kind, field[x][y].kind, b)
			}
		}
	}
	return res
}

// updateScores goes through every cell, and plays the Prisoner's dilema game
// with each of it's in-field nieghbors (including itself). It updates the
// score of each cell to be the sum of that cell's winnings from the game.
func updateScores(field [][]Cell, b float64) {
	row := len(field)
	col := len(field[0])
    for i := 0; i < row; i++{
    	for j := 0; j < col; j++{
    		field[i][j].score = game(field, i, j, b)
    	}
    }
}

// updateStrategies create a new field by going through every cell (r,c), and
// looking at each of the cells in its neighborhood (including itself) and the
// setting the kind of cell (r,c) in the new field to be the kind of the
// neighbor with the largest score
func updateStrategies(field [][]Cell) [][]Cell {
	row := len(field)
	col := len(field[0])
	f := createField(row, col)
	for i := 0; i < row; i++{
		for j:= 0; j < col; j++{
			f[i][j].score = field[i][j].score
			f[i][j].kind = field[i][j].kind
		}
	}
	for i := 0; i < row; i++{
		for j := 0; j < col; j++{
			field[i][j].pre_kind = field[i][j].kind
			field[i][j].kind = decide(f, i, j)
		}
	} 
	return field
}


/* 
 * decide is the strategy decision made process of the current cell. Updates
 * its strategy if its neighbor gets greater score than it.
 */
func decide(field [][]Cell, i, j int) string{
	var res string = field[i][j].kind
	var temp float64 = field[i][j].score
	for x := i - 1; x <= i + 1; x++{
		for y:= j - 1; y <= j + 1; y++{
			if inField(field, x, y){
				if temp < field[x][y].score {
					temp = field[x][y].score
					res = field[x][y].kind
				} 
			}
		}
	}
	return res
}

// evolve takes an intial field and evolves it for nsteps according to the game
// rule. At each step, it should call "updateScores()" and the updateStrategies
func evolve(field [][]Cell, nsteps int, b float64) [][]Cell {
	for i := 0; i < nsteps; i++ {
		updateScores(field, b)	
		//print(field)	
 		field = updateStrategies(field)
 		//print(field)
		//fmt.Println()
	}

	return field
}
/*
 * print function for testing.
 */
func print(field [][]Cell){
	for x := 0; x < len(field); x++{
		for y := 0 ; y < len(field[0]); y++{
			fmt.Print(field[x][y])
		}
		fmt.Println()
	}	
}
// Implements a Spatial Games version of prisoner's dilemma. The command-line
// usage is:
//     ./spatial field_file b nsteps
// where 'field_file' is the file continaing the initial arrangment of cells, b
// is the reward for defecting against a cooperator, and nsteps is the number
// of rounds to update stategies.
//
func main() {
	// parse the command line
	if len(os.Args) != 4 {
		fmt.Println("Error: should spatial field_file b nsteps")
		return
	}

	fieldFile := os.Args[1]

	b, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil || b <= 0 {
		fmt.Println("Error: bad b parameter.")
		return
	}

	nsteps, err := strconv.Atoi(os.Args[3])
	if err != nil || nsteps < 0 {
		fmt.Println("Error: bad number of steps.")
		return
	}

    // read the field
	field := readFieldFromFile(fieldFile)
    fmt.Println("Field dimensions are:", len(field), "by", len(field[0]))

    // evolve the field for nsteps and write it as a PNG
	field = evolve(field, nsteps, b)
	drawField(field, "Prisoners.png")
}
