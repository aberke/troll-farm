package trolls

import (
        "fmt"
        "log"
        "os"
        "math"
        "encoding/json"
)

const GRID_WIDTH    = 10
const GRID_HEIGHT   = 10

const GRID_CAPACITY = 5 // if there are this many items, IsFull returns true

// any non-troll GridItem (goes in funCells) has negative id
const FOODBUTTON_ID = -1
const BANANA_ID     = -2

const BANANA_REWARD = 1 // troll gets 1 point per banana


/* The grid is a 2D array that maps (x,y) positions to the key of a GridItem in the gridItemsMap */
func createCells () [][]int {
    // Allocate the top-level slice.
    cells := make([][]int, GRID_HEIGHT)  // One row per unit of y.
    // Loop over the rows, allocating the slice for each row.
    for i := range cells {
        cells[i] = make([]int, GRID_WIDTH)
    }
    return cells
}

var totalGrids int = 0
var minGridId  int = 0
var maxGridId  int = -1 // there are no grids at first

/* The Grid has
    id
    trollCells - a 2d array to map an (x,y) cell to id of troll-GridItems there (will be positive ID)
    funCells   - a 2d array to map an (x,y) cell to id of any non-troll GridItem there (will be negative ID)
    itemsMap   - maps ids of GridItems to GridItem
    update     - like itemsMap but only for recent changes
*/
type Grid struct {
    id              int
    trollCells      [][]int
    funCells        [][]int
    itemsMap        map[int]*GridItem
    updateMap       map[int]*GridItem
}    
func NewGrid () *Grid {
    trollCells  := createCells()
    funCells    := createCells()
    itemsMap    := make(map[int]*GridItem)
    updateMap   := make(map[int]*GridItem)

    maxGridId ++
    totalGrids ++
    g :=  &Grid{ maxGridId, trollCells, funCells, itemsMap, updateMap }

    // add the food button to the grid
    foodButton := NewGridItem("FOODBUTTON", 9, 9)
    g.itemsMap[FOODBUTTON_ID] = foodButton
    g.funCells[9][9] = FOODBUTTON_ID

    return g
}
/* if grid is safe to remove -- gets ready for removal by decrementing maxGridId and returns true 
    otherwise returns false
    only ever called by server right before potential removal (server removes after call iff returns true)
*/
func (g *Grid) SafelyRemove () bool{
    /* if there is a grid after this one, we don't want to leave it as an island */
    if (g.id < maxGridId) {
        return false
    }

    /* if there are any trolls still in the itemsMap, then can't remove */
    for itemId, _ := range g.itemsMap {
        if (itemId > 0) {
            return false
        }
    }
    maxGridId --
    totalGrids --
    return true
}


// Troll client JSON data
type GridItem struct {
    Name        string  // used names: {"DELETE": indicates to client to delete troll}
    Color       string
    Coordinates map[string]int
    Messages    []string
    Points      int64
}

// Create new GridItem
func NewGridItem (name string, x int, y int) *GridItem{
    log.Println("*** NewGridItem *****")

    coordinates     := make(map[string]int)
    coordinates["x"] = x
    coordinates["y"] = y
    messages        := make([]string, 5)
    gi := GridItem{name, "#FF00FF", coordinates, messages, 0}

    encodedGi, err := json.MarshalIndent(gi, "", " ")
    if err != nil {
        fmt.Println("****** err *****", err)
    }
    os.Stdout.Write(encodedGi)

    return &gi
}
func (g *Grid) bananaCollision(trollID int) {
    log.Println("bananaCollision: TODO")
    g.removeBanana()

    // give the troll that collided with the banana points
    gi := g.itemsMap[trollID]
    gi.Points += BANANA_REWARD
    g.updateMap[trollID] = gi
}
func (g *Grid) foodButtonCollision() {
    log.Println("foodButtonPressed")
    if (g.itemsMap[BANANA_ID] == nil) {
        g.generateBanana()    
    }
}
func (g *Grid) removeBanana() {
    gi := g.itemsMap[BANANA_ID]
    g.removeFromCell(g.funCells, gi)
    delete(g.itemsMap, BANANA_ID)
}
/* places a banana at the bottom left corner or next best empty spot */
func (g *Grid) generateBanana() {
    x, y := g.emptySpot(0, GRID_HEIGHT - 1)
    banana := NewGridItem("BANANA", x, y)
    g.itemsMap[BANANA_ID] = banana
    g.funCells[x][y] = BANANA_ID
    g.updateMap[BANANA_ID] = banana
}
/* Takes as parameters the coordinates of the desired space.  Returns next best empty space */
func (g *Grid) emptySpot(x int, y int) (retX int, retY int) {
    count := 0

    for (count < GRID_WIDTH*GRID_HEIGHT) {
        if (x >= GRID_WIDTH) {
            x = 0
            y += 1
        }
        if (y >= GRID_HEIGHT) {
            y = 0
            x = int(math.Mod(float64(x + 1), GRID_WIDTH))
        }

        if ((g.trollCells[x][y] ==0 && g.funCells[x][y] == 0)) {
            return x, y
        }
        x += 1
        count += 1
    }
    panic("No more empty spots on grid")
}

/*********************************************************/
// getter functions
/*********************************************************/
func (g *Grid) UpdateMap() map[int]*GridItem{
    return g.updateMap
}
func (g *Grid) ItemsMap() map[int]*GridItem{
    return g.itemsMap
}
/* returns boolean -- true if full, false otherwise */
func (g *Grid) IsFull() bool {
    if (len(g.itemsMap) >= GRID_CAPACITY) {
        return true
    }
    return false
}
/*********************************************************/
// setter functions
/*********************************************************/

// returns false if move is not valid (collision)
func (g *Grid) MoveTroll(trollID int, moveX int, moveY int) bool {
    log.Println("MoveTroll")
    gi :=g.itemsMap[trollID]
    // retrieve troll client's current position
    currentX := gi.Coordinates["x"]
    currentY := gi.Coordinates["y"]
    // calculate requested new position coordinates
    requestedX := (currentX + moveX)
    requestedY := (currentY + moveY)

    // collision detection with grid boundaries
    if (requestedX < 0 || requestedX >= GRID_WIDTH || requestedY < 0 || requestedY >= GRID_HEIGHT) {
        return false
    }
    // collision detection with other trolls
    if (g.trollCells[requestedX][requestedY] != 0) { 
        return false
    }

    // check if ran into non-troll item
    nonTrollItem := g.funCells[requestedX][requestedY]
    if (nonTrollItem != 0) {
        log.Println("nonTrollItem: ", nonTrollItem)
    }
    switch nonTrollItem {
        case FOODBUTTON_ID: g.foodButtonCollision()
        case BANANA_ID:     g.bananaCollision(trollID)
    }

    // move that troll
    g.trollCells[currentX][currentY] = 0
    g.trollCells[requestedX][requestedY] = trollID

    gi.Coordinates["x"] = requestedX
    gi.Coordinates["y"] = requestedY

    g.updateMap[trollID] = gi
    
    return true
}
func (g *Grid) AddTroll(trollID int) {

    x, y := g.emptySpot(0, 0)
    gi := NewGridItem("no-name", x, y)
    g.trollCells[x][y] = trollID
    g.itemsMap[trollID] = gi
    g.updateMap[trollID] = gi
}
func (g *Grid) ClearUpdateMap() {
    g.updateMap = make(map[int]*GridItem)
}
func (g *Grid) DeleteTroll(trollID int) {
    gi := g.itemsMap[trollID]

    // set troll to be deleted in updateMap
    gi.Name = "DELETE"
    g.updateMap[trollID] = gi

    // delete troll GridItem
    g.removeFromCell(g.trollCells, gi)
    delete(g.itemsMap, trollID)
}

func (g *Grid) removeFromCell(cells [][]int, gi *GridItem) {
    if (cells[gi.Coordinates["x"]][gi.Coordinates["y"]] == 0) {
        panic ("Tried to removeGridItem where no GridItem is located")
    }
    cells[gi.Coordinates["x"]][gi.Coordinates["y"]] = 0
}


