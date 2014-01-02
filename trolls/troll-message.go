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
	TrollsMap 	map[string]TrollData
	LocalTroll  int
}

func NewOutgoingMessage(msgType string, localTroll int, trollsDataMap map[int]*TrollData) *OutgoingMessage {
	var trolls map[string]TrollData
	trolls = nil
	if (trollsDataMap != nil) {
		trolls = JSONifyTrollsDataMap(trollsDataMap)
	}

	return &OutgoingMessage{
		msgType,
		trolls,
		localTroll,
	}
}

func OutgoingTrollsMessage(localTroll int, trollsMap map[int]*TrollData) *OutgoingMessage {
	return NewOutgoingMessage("trolls", localTroll, trollsMap)
}
func OutgoingUpdateMessage(localTroll int, trollsMap map[int]*TrollData) *OutgoingMessage {
	return NewOutgoingMessage("update", localTroll, trollsMap)
}

func OutgoingTestMessage(localTroll int) *OutgoingMessage {
	nilMap := make(map[int]*TrollData)
	nilMap = nil
	return NewOutgoingMessage("test", localTroll, nilMap)
}

func JSONifyTrollsDataMap(trollsDataMap map[int]*TrollData) map[string]TrollData {
	m := make(map[string]TrollData)
	for trollID, trollData := range trollsDataMap {
		trollIDString := strconv.Itoa(trollID) // json object can't have ints as keys
		
		if (trollData != nil) { // nil if we're signaling this troll removed
			m[trollIDString] = *trollData
		}
	}
	return m
}
