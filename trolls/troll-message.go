package trolls

import (
	"fmt"
	"strconv"
)



type IncomingMessage struct {
	Type 		string `json:"message-type"`
	Data   		string `json:"data"`
}

type OutgoingMessage struct {
	Type 		string
	LocalTroll  int
	TrollsMap 	map[string]TrollData
}

func NewOutgoingMessage(msgType string, localTroll int, trollsDataMap map[int]*TrollData) *OutgoingMessage {
	fmt.Println("888888 NewOutgoingMessage")
	var trolls map[string]TrollData
	trolls = nil
	if (trollsDataMap != nil) {
		fmt.Println("trollsDataMap not nil")
		trolls = JSONifyTrollsDataMap(trollsDataMap)
	} 
	fmt.Println("trolls", trolls)

	return &OutgoingMessage{
		msgType,
		localTroll,
		trolls,
	}
}

func OutgoingTrollsMessage(localTroll int, trollsMap map[int]*TrollData) *OutgoingMessage {
	return NewOutgoingMessage("trolls", localTroll, trollsMap)
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
		m[trollIDString] = *trollData
	}
	return m
}
