package constellation

import (
	"encoding/json"
	"fmt"
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

type NodeType int

// define NodeTypes
const (
	Satellite NodeType = iota
	GroundStation
)

// what is a Node
type Node struct {
	ID        string
	Type      NodeType
	Position  Vec3
	Velocity  Vec3
	CommRange float64
}

type ExportData struct {
	Nodes []ExportNode `json:"nodes"`
	Edges []ExportEdge `json:"edges"`
	Route []string     `json:"route,omitempty"`
}

type ExportNode struct {
	ID   string  `json:"id"`
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
}

type ExportEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Constellation struct {
	Storage      map[string]*Node
	AdjacentList map[string][]string
}

func New() *Constellation {
	return &Constellation{
		Storage:      make(map[string]*Node),
		AdjacentList: make(map[string][]string),
	}
}

func (c *Constellation) AddNode(node *Node) {
	// check if it's in map
	if _, exists := c.Storage[node.ID]; exists {
		fmt.Println("node exists in storage")
		return
	}

	// add to storage
	c.Storage[node.ID] = node

	// init empty adjacent list
	c.AdjacentList[node.ID] = make([]string, 0)
}

// basically, rebuild the adjacency list
func (c *Constellation) UpdateLinks() {
	// clear
	for id := range c.AdjacentList {
		c.AdjacentList[id] = []string{}
	}

	// for all the pairs in Storage, check if they're in range
	for _, nodeA := range c.Storage {
		for _, nodeB := range c.Storage {
			if nodeA.ID >= nodeB.ID {
				continue // skip self and duplicates
			}

			// check if nodeA and nodeB are in range
			dis := c.distance(nodeB.Position, nodeA.Position)
			if dis <= min(nodeA.CommRange, nodeB.CommRange) {
				// add for both directions
				c.AdjacentList[nodeA.ID] = append(c.AdjacentList[nodeA.ID], nodeB.ID)
				c.AdjacentList[nodeB.ID] = append(c.AdjacentList[nodeB.ID], nodeA.ID)
			}
		}
	}

}

// return back what the fromID reaches
func (c Constellation) GetReachable(fromID string) []string {
	path, exists := c.AdjacentList[fromID]
	if !exists {
		return nil
	}
	return path
}

// how to route from source to dest
func (c *Constellation) Route(fromID, toID string) []string {
	// best known distance for each node
	dist := map[string]float64{}

	// track where you came from
	prev := map[string]string{}

	// track visited
	visited := map[string]bool{}

	// init all distances to infinity, except source
	for id := range c.Storage {
		dist[id] = math.Inf(1)
	}
	dist[fromID] = 0

	// pick the unvisisted node with smallest distance
	for {
		// find unvisited with minimum distance
		minDist := math.Inf(1)
		current := ""
		for id, d := range dist {
			if !visited[id] && d < minDist {
				minDist = d
				current = id
			}
		}

		// if no node found or reached, then stop
		if current == "" || current == toID {
			break
		}

		// add to visited
		visited[current] = true

		// update neighbors
		for _, neighborID := range c.AdjacentList[current] {
			// new distance through current node
			n := c.Storage[neighborID]
			cur := c.Storage[current]
			newDist := dist[current] + c.distance(cur.Position, n.Position)
			if newDist < dist[neighborID] {
				dist[neighborID] = newDist
				prev[neighborID] = current
			}
		}
	}

	path := []string{}
	current := toID

	for current != "" {
		path = append(path, current)
		current = prev[current]
	}

	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

func (c *Constellation) distance(pos1, pos2 Vec3) float64 {
	dx := pos2.X - pos1.X
	dy := pos2.Y - pos1.Y
	dz := pos2.Z - pos1.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (c *Constellation) ExportJSON() ([]byte, error) {
	data := ExportData{}

	for id, node := range c.Storage {
		nodeType := "satellite"
		if node.Type == GroundStation {
			nodeType = "ground"
		}
		data.Nodes = append(data.Nodes, ExportNode{
			ID:   id,
			Type: nodeType,
			X:    node.Position.X,
			Y:    node.Position.Y,
			Z:    node.Position.Z,
		})
	}

	// Only add each edge once
	seen := map[string]bool{}
	for from, neighbors := range c.AdjacentList {
		for _, to := range neighbors {
			key := from + "-" + to
			keyRev := to + "-" + from
			if !seen[key] && !seen[keyRev] {
				data.Edges = append(data.Edges, ExportEdge{From: from, To: to})
				seen[key] = true
			}
		}
	}

	return json.MarshalIndent(data, "", "  ")
}

func (c *Constellation) ExportJSONWithRoute(fromID, toID string) ([]byte, error) {
	data := ExportData{}

	for id, node := range c.Storage {
		nodeType := "satellite"
		if node.Type == GroundStation {
			nodeType = "ground"
		}
		data.Nodes = append(data.Nodes, ExportNode{
			ID:   id,
			Type: nodeType,
			X:    node.Position.X,
			Y:    node.Position.Y,
			Z:    node.Position.Z,
		})
	}

	// Only add each edge once
	seen := map[string]bool{}
	for from, neighbors := range c.AdjacentList {
		for _, to := range neighbors {
			key := from + "-" + to
			keyRev := to + "-" + from
			if !seen[key] && !seen[keyRev] {
				data.Edges = append(data.Edges, ExportEdge{From: from, To: to})
				seen[key] = true
			}
		}
	}

	// Add the route
	data.Route = c.Route(fromID, toID)

	return json.MarshalIndent(data, "", "  ")
}
