/* The server has a GridMap tand this GridMap does a lot of work managing the Grids */

package trolls

import (
        "fmt"
)



type GridMap struct {
    minGridId       int
    grids           map[int]*Grid
}
func NewGridMap() *GridMap {
    minGridId       := 0
    grids           := make(map[int]*Grid)
    return &GridMap{ minGridId, grids }
}

// getter function to get Grid by ID from GridMap
func (gm *GridMap) Grid(gId int) *Grid {
    return gm.grids[gId]
}
/* Returns -> int:  gridID that Troll now lives in
			  bool: true if requested move was valid (and troll therefore moved), false otherwise
*/
func (gm *GridMap) MoveTroll(trollID int, gridID int, moveX int, moveY int) (int, bool) {
	var grid *Grid = gm.grids[gridID]
	newGridID, validMove := grid.MoveTroll(trollID, moveX, moveY)
	
	/* check and handle Troll moving to new Grid */
	if (newGridID != gridID) {
		gm.grids[gridID].DeleteTroll(trollID)
		
		if (newGridID >= len(gm.grids)) {
			gm.AddGrid()
		}
		gm.grids[newGridID].AddTroll(trollID)
	}
	return newGridID, validMove
}
func (gm *GridMap) AddGrid() {
	gId := len(gm.grids)
	var grid *Grid = NewGrid(gId)
	gm.grids[gId] = grid
}
/* finds the next available Grid to add Troll to or creates new grid
    returns GridID that Troll was added to */
func (gm *GridMap) AddTroll(tId int) int {
    gId := gm.minGridId
    for ((gId < len(gm.grids)) && gm.grids[gId].IsFull()) {
        gId ++
    }
    if (gm.grids[gId] == nil) {
        gm.AddGrid()
        fmt.Println("******** Added new Grid - Now ", len(gm.grids), "grids.")
    }

    gm.grids[gId].AddTroll(tId)
    fmt.Println("GridMap grids", gm.grids, "gId", gId)
    return gId
}
/* Deletes Troll from its Grid 
    and checks if that Grid is ready for removal - if it is it removes it
*/
func (gm *GridMap) DeleteTroll(gId int, tId int) error{
    grid := gm.grids[gId]
    if (grid == nil) {
        return fmt.Errorf("GridMap: No Grid with id %i exists in GridMap", gId)
    }
    err  := grid.DeleteTroll(tId)
    if (err != nil) {
        return err
    }
    /* if grid is empty and the last grid -- then remove it */
    if (gm.SafelyRemove(gId)) {
        fmt.Println("Removed Grid", gId, "from GridMap - Now ", len(gm.grids), "grids.")
    }
    return nil
}
/* checks if Grid with id gId is safe to remove
    if true: removes it and returns true
    otherwise: returns false
        only ever called by GridMap.DeleteTroll 
*/
func (gm *GridMap) SafelyRemove (gId int) bool{
    /* if there is a grid after this one, we don't want to leave it as an island */
    if (gId < (len(gm.grids) - 1)) {
        return false
    }

    /* if there are any trolls still in the itemsMap, then can't remove */
    for itemId, _ := range gm.grids[gId].itemsMap {
        if (itemId > 0) {
            return false
        }
    }
    /* Grid is safe for removal! */
    delete(gm.grids, gId)
    return true
}