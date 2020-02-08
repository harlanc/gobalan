package node

import (
	"sync"

	"github.com/harlanc/gobalan/logger"
	pb "github.com/harlanc/gobalan/proto"
)

var (
	//NodeContainer Work node list
	NodeContainer *WorkerNodeContainer = &WorkerNodeContainer{nodeList: make([]*Node, 0, MaxNodeNumber), workID2Index: make(map[uint32]int)}
	//MaxNodeNumber max node number
	MaxNodeNumber uint = 64
)

//Node store a server info
type Node struct {
	IP            string
	Port          int32
	WorkerID      uint32
	Stat          *pb.Stat
	ServiceStatus pb.ServiceStatus
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

	logger.LogWarnf("Delete worker id %d\n", workerid)

	var idx int
	var ok bool

	if idx, ok = ws.workID2Index[workerid]; !ok {
		logger.LogErrf("The worker id  %d does not exist!\n", workerid)
		return
	}

	l := len(ws.nodeList)

	for _, v := range ws.nodeList[idx+1:] {
		ws.workID2Index[v.WorkerID]--
	}

	copy(ws.nodeList[idx:], ws.nodeList[idx+1:])
	ws.nodeList[l-1] = nil
	ws.nodeList = ws.nodeList[:l-1]
	delete(ws.workID2Index, workerid)
}

//UpdateNodeStat delete a worker node
func (ws *WorkerNodeContainer) UpdateNodeStat(workerid uint32, stat *pb.Stat) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	var idx int
	var ok bool

	if idx, ok = ws.workID2Index[workerid]; !ok {
		logger.LogErrf("The worker id  %d does not exist!\n", workerid)
		return
	}
	ws.nodeList[idx].Stat = stat
}

//UpdateNodeServiceStatus delete a worker node
func (ws *WorkerNodeContainer) UpdateNodeServiceStatus(workerid uint32, status pb.ServiceStatus) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	var idx int
	var ok bool

	if idx, ok = ws.workID2Index[workerid]; !ok {
		logger.LogErrf("The worker id  %d does not exist!\n", workerid)
		return
	}
	ws.nodeList[idx].ServiceStatus = status
}

//GetNodeListLen get the length
func (ws *WorkerNodeContainer) GetNodeListLen() int {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	return len(ws.nodeList)
}

//GetNodeList get the node list
func (ws *WorkerNodeContainer) GetNodeList() []*Node {

	ws.mu.RLock()
	defer ws.mu.RUnlock()

	//deep copy
	newlist := make([]*Node, 0, ws.GetNodeListLen())

	for _, v := range ws.nodeList {

		newstat := new(pb.Stat)
		if v.Stat != nil {
			*newstat = *v.Stat
		}
		newnode := &Node{IP: v.IP, Port: v.Port, WorkerID: v.WorkerID, Stat: newstat}
		newlist = append(newlist, newnode)

	}

	return newlist
}
