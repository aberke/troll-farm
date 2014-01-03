package trolls

import (
        "fmt"
        "log"
        "os"
        "math"
        "encoding/json"
)

const GRID_WIDTH  = 10
const GRID_HEIGHT = 10

// special GridItem ids
const FOODBUTTON_ID = -1
const BANANA_ID     = -2


/* The grid is a 2D array that maps (x,y) positions to the key of a GridItem in the gridItemsMap */
func NewGrid () [][]int {
    // Allocate the top-level slice.
    grid := make([][]int, GRID_HEIGHT)  // One row per unit of y.
    // Loop over the rows, allocating the slice for each row.
    for i := range grid {
        grid[i] = make([]int, GRID_WIDTH)
    }
    return grid
}

// Troll client JSON data
type GridItem struct {
    Name        string  // used names: {"DELETE": indicates to client to delete troll}
    Color       string
    Coordinates map[string]int
    Messages    []string
    Points      int64
}
func NewFoodButton () *GridItem {
    log.Println("*** NewFoodButton *****")

    coordinates     := make(map[string]int)
    coordinates["x"] = 9
    coordinates["y"] = 9
    messages        := make([]string, 0)
    gi := GridItem{"FOODBUTTON", "#FFFFFF", coordinates, messages, 0}

    encodedGi, err := json.MarshalIndent(gi, "", " ")
    if err != nil {
        fmt.Println("****** err *****", err)
    }
    os.Stdout.Write(encodedGi)

    return &gi
}
// Create new TrollData from Troll
func (t *Troll) NewGridItem () *GridItem{
    log.Println("*** Troll.NewGridItem *****")

    coordinates     := make(map[string]int)
    coordinates["x"] = int(math.Mod(float64(t.id), 9))
    coordinates["y"] = 0
    messages        := make([]string, 5)
    gi := GridItem{"no-name", "#FF00FF", coordinates, messages, 0}

    encodedGi, err := json.MarshalIndent(gi, "", " ")
    if err != nil {
        fmt.Println("****** err *****", err)
    }
    os.Stdout.Write(encodedGi)

    return &gi
}

func RemoveGridItem(grid [][]int, gi *GridItem) {
    grid[gi.Coordinates["x"]][gi.Coordinates["y"]] = 0
}