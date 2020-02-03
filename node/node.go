package node

import (
	"fmt"
	pb "github.com/harlanc/gobalan/proto"
	"sync"
)

var (
	//NodeContainer Work node list
	NodeContainer *WorkerNodeContainer = &WorkerNodeContainer{nodeList: make([]*Node, MaxNodeNumber), workID2Index: make(map[uint32]int)}
	//MaxNodeNumber max node number
	MaxNodeNumber uint = 64
)

//Node store a server info
type Node struct {
	IP       string
	Port     string
	WorkerID uint32
	Stat     *pb.Stat
}

//WorkerNodeContainer collection worker node list
type WorkerNodeContainer struct {
	nodeList     []*Node
	workID2Index map[uint32]int

	mu sync.RWMutex
}

//InsertNode insert a worker node
func (ws *WorkerNodeContainer) InsertNode(workerid uint32, node *Node) {

	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.workID2Index[workerid] = len(ws.nodeList)
	ws.nodeList = append(ws.nodeList, node)

}

//DeleteNode delete a worker node
func (ws *WorkerNodeContainer) DeleteNode(workerid uint32) {

	ws.mu.Lock()
	defer ws.mu.Unlock()

	var idx int
	var ok bool

	if idx, ok = ws.workID2Index[workerid]; !ok {
		fmt.Printf("The worker id  %d does not exist!", workerid)
		return
	}

	l := len(ws.nodeList)

	copy(ws.nodeList[idx:], ws.nodeList[idx+1:])
	ws.nodeList[l-1] = nil
	ws.nodeList = ws.nodeList[:l-1]
	delete(ws.workID2Index, workerid)
}

//UpdateNode delete a worker node
func (ws *WorkerNodeContainer) UpdateNode(workerid uint32, stat *pb.Stat) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	var idx int
	var ok bool

	if idx, ok = ws.workID2Index[workerid]; !ok {
		fmt.Printf("The worker id  %d does not exist!", workerid)
		return
	}
	ws.nodeList[idx].Stat = stat

}

//GetNodeList get the node list
func (ws *WorkerNodeContainer) GetNodeList() []*Node {

	ws.mu.RLock()
	defer ws.mu.RLock()

	return ws.nodeList
}
