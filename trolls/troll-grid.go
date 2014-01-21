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

const DOOR_Y        = 5 // the Y value that doors live on

const GRID_CAPACITY = 5 // if there are this many items, IsFull returns true

// any non-troll GridItem (goes in funCells) has negative id
const FOODBUTTON_ID     = -1
const BANANA_ID         = -2
const DOOR_PREVIOUS_ID  = -3 // door to previous Grid
const DOOR_NEXT_ID      = -4 // door to next Grid
    
const BANANA_REWARD = 1 // troll gets 1 point per banana



/* The Grid has
    id         - corresponds to where Grid resides in the GridMap's list of Grids -- determined by GridMap 
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
func NewGrid (gId int) *Grid {
    trollCells  := createCells()
    funCells    := createCells()
    itemsMap    := make(map[int]*GridItem)
    updateMap   := make(map[int]*GridItem)

    g :=  &Grid{ gId, trollCells, funCells, itemsMap, updateMap }

    // add the food button to the grid
    foodButton := NewGridItem("FOODBUTTON", 9, 9)
    g.itemsMap[FOODBUTTON_ID] = foodButton
    g.funCells[9][9] = FOODBUTTON_ID

    // add the door(s)
    if (gId != 0) {
        door_previous := NewGridItem("DOOR-PREVIOUS", 0, DOOR_Y)
        g.itemsMap[DOOR_PREVIOUS_ID] = door_previous
        g.funCells[0][DOOR_Y] = DOOR_PREVIOUS_ID
    }
    door_next := NewGridItem("DOOR-NEXT", GRID_WIDTH-1, DOOR_Y)
    g.itemsMap[DOOR_NEXT_ID] = door_next
    g.funCells[GRID_WIDTH-1][DOOR_Y] = DOOR_NEXT_ID

    return g
}
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
        } else if(y >= GRID_HEIGHT) {
            y = 0
            x = int(math.Mod(float64(x + 1), GRID_WIDTH))
        } else {
            if ((g.trollCells[x][y] == 0 && g.funCells[x][y] == 0)) {
                return x, y
            }
            x += 1
            count += 1
        }
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

/* 
- Checks whether move valid (no collision with boundary)
    if collision -> returns (gridId, false)
- Checks whether Troll is walking through previous/next door
    if so -> returns (decremented/incremented gridID, true)
- Otherwise moves Troll -> returns (gridId, true)

Returns -> int:  gridID that Troll now lives in
           bool: true if requested move was valid (and troll therefore moved), false otherwise
*/
func (g *Grid) MoveTroll(trollID int, moveX int, moveY int) (int, bool) {
    gi := g.itemsMap[trollID]
    
    // retrieve troll client's current position
    currentX := gi.Coordinates["x"]
    currentY := gi.Coordinates["y"]
    
    // calculate requested new position coordinates
    requestedX := (currentX + moveX)
    requestedY := (currentY + moveY)
    
    currFunItem := g.funCells[currentX][currentY]
    
    // collision detection with left/right grid boundaries
    if (requestedX < 0 || requestedX >= GRID_WIDTH) {
        // if was on a door, then go through the door, otherwise invalid
        if (currFunItem == DOOR_PREVIOUS_ID) {
            return (g.id - 1), true
        } else if (currFunItem == DOOR_NEXT_ID) {
            return (g.id + 1), true
        }
        return g.id, false
    }

    // collision detection with top/bottom grid boundaries
    if (requestedY < 0 || requestedY >= GRID_HEIGHT) {
        return g.id, false
    }
    // collision detection with other trolls
    if (g.trollCells[requestedX][requestedY] != 0) { 
        return g.id, false
    }
    // check if ran into non-troll item
    nonTrollItem := g.funCells[requestedX][requestedY]
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
    return g.id, true
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
func (g *Grid) DeleteTroll(trollID int) error{
    gi := g.itemsMap[trollID]
    if (gi == nil) {
        return fmt.Errorf("Troll with id %i does not exist in Grid with id %i", trollID, g.id)
    }

    // set troll to be deleted in updateMap
    gi.Name = "DELETE"
    g.updateMap[trollID] = gi

    // delete troll GridItem
    g.removeFromCell(g.trollCells, gi)
    delete(g.itemsMap, trollID)
    return nil
}

func (g *Grid) removeFromCell(cells [][]int, gi *GridItem) {
    if (cells[gi.Coordinates["x"]][gi.Coordinates["y"]] == 0) {
        panic ("Tried to removeGridItem where no GridItem is located")
    }
    cells[gi.Coordinates["x"]][gi.Coordinates["y"]] = 0
}


