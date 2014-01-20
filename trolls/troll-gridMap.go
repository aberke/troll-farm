/* The server has a GridMap tand this GridMap does a lot of work managing the Grids */

package trolls

import (
        "fmt"
)



type GridMap struct {
    totalGrids      int
    minGridId       int
    grids           map[int]*Grid
}
func NewGridMap() *GridMap {
    totalGrids      := 0
    minGridId       := 0
    grids           := make(map[int]*Grid)

    return &GridMap{ totalGrids, minGridId, grids }
}

// getter function to get Grid by ID from GridMap
func (gm *GridMap) Grid(gId int) *Grid {
    return gm.grids[gId]
}
/* finds the next available Grid to add Troll to or creates new grid
    returns GridID that Troll was added to */
func (gm *GridMap) AddTroll(tId int) int {
    gId := gm.minGridId
    for ((gId < gm.totalGrids) && gm.grids[gId].IsFull()) {
        gId ++
    }
    if (gm.grids[gId] == nil) {
        g := NewGrid(gId)
        gm.grids[gId] = g
        gm.totalGrids ++
        fmt.Println("******** Added new Grid - Now ", len(gm.grids), "grids.")
    }

    gm.grids[gId].AddTroll(tId)
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
    if (gId < (gm.totalGrids -1)) {
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
    gm.totalGrids --
    return true
}