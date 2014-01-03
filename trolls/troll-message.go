package trolls

import (
	//"fmt"
	"strconv"
)



type IncomingMessage struct {
	Type 		string `json:"message-type"`
	Data   		map[string]string `json:"data"`
	LocalTroll	int
}

type OutgoingMessage struct {
	Type 		string
	ItemsMap 	map[string]GridItem
	LocalTroll  int
}

func NewOutgoingMessage(msgType string, localTroll int, gridItemsMap map[int]*GridItem) *OutgoingMessage {
	var items map[string]GridItem
	items = nil
	if (gridItemsMap != nil) {
		items = JSONifyGridItemsMap(gridItemsMap)
	}

	return &OutgoingMessage{
		msgType,
		items,
		localTroll,
	}
}

func OutgoingItemsMessage(localTroll int, itemsMap map[int]*GridItem) *OutgoingMessage {
	return NewOutgoingMessage("items", localTroll, itemsMap)
}
func OutgoingUpdateMessage(localTroll int, itemsMap map[int]*GridItem) *OutgoingMessage {
	return NewOutgoingMessage("update", localTroll, itemsMap)
}

func OutgoingPingMessage(localTroll int) *OutgoingMessage {
	nilMap := make(map[int]*GridItem)
	nilMap = nil
	return NewOutgoingMessage("ping", localTroll, nilMap)
}

func JSONifyGridItemsMap(gridItemsMap map[int]*GridItem) map[string]GridItem {
	m := make(map[string]GridItem)
	for trollID, trollData := range gridItemsMap {
		trollIDString := strconv.Itoa(trollID) // json object can't have ints as keys
		
		if (trollData != nil) { // nil if we're signaling this troll removed
			m[trollIDString] = *trollData
		}
	}
	return m
}
