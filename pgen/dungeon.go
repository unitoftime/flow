package pgen

import (
	"fmt"
	"math"
	"math/rand"
	"github.com/unitoftime/flow/phy2"
	"github.com/unitoftime/flow/tile"
	"github.com/unitoftime/flow/ds"
)

// type RoomDagNode struct {
// 	Data string
// }

type RoomDag struct {
	Nodes []string // Holds the nodes in the order that they were added
	NodeMap map[string]int // holds the rank of the node (to prevent cycles)
	Edges map[string][]string // This is the dag
	rank int
}
func NewRoomDag() *RoomDag {
	return &RoomDag{
		Nodes: make([]string, 0),
		NodeMap: make(map[string]int),
		Edges: make(map[string][]string),
	}
}
func (d *RoomDag) AddNode(label string) bool {
	_, ok := d.NodeMap[label]
	if ok { return false } // already exists

	d.NodeMap[label] = d.rank
	d.Nodes = append(d.Nodes, label)
	d.Edges[label] = make([]string, 0)
	d.rank++
	return true
}

func (d *RoomDag) AddEdge(from, to string) {
	fromNode, ok := d.NodeMap[from]
	if !ok { return } // TODO: false?
	toNode, ok := d.NodeMap[to]
	if !ok { return } // TODO: false?

	// Only add edges from lower to higher ranks
	if fromNode < toNode {
		d.Edges[from] = append(d.Edges[from], to)
	} else {
		d.Edges[to] = append(d.Edges[to], from)
	}

}

func (d *RoomDag) HasEdgeEitherDirection(from, to string) bool {
	return d.HasEdge(from, to) || d.HasEdge(to, from)
}

func (d *RoomDag) HasEdge(from, to string) bool {
	edges := d.Edges[from]
	for i := range edges {
		if to == edges[i] {
			return true
		}
	}
	return false
}

// returns the node labels in topological order
func (d *RoomDag) TopologicalSort(start string) []string {
	visit := make(map[string]bool)
	queue := ds.NewQueue[string](0) // TODO: dynamically resized
	queue.Add(start)

	ret := make([]string, 0, len(d.Nodes))

	for {
		currentNode, ok := queue.Remove()
		if !ok { break }

		visited := visit[currentNode]
		if visited { continue }
		visit[currentNode] = true

		ret = append(ret, currentNode)

		children := d.Edges[currentNode]
		for _, child := range children {
			queue.Add(child)
		}
	}
	return ret
}

func GenerateRandomGridWalkDag(num int, numWalks int) *RoomDag {
	dag := NewRoomDag()
	pos := make([]tile.TilePosition, numWalks)
	lastLabel := make([]string, numWalks)

	for len(dag.Nodes) < num {
		for i := range pos {
			label := fmt.Sprintf("%d_%d", pos[i].X, pos[i].Y)
			_, exists := dag.NodeMap[label]
			if !exists {
				dag.AddNode(label)
			}

			dir := rand.Intn(4)
			if dir == 0 {
				pos[i].X++
			} else if dir == 1 {
				pos[i].X--
			} else if dir == 2 {
				pos[i].Y++
			} else {
				pos[i].Y--
			}

			if lastLabel[i] != "" {
				dag.AddEdge(lastLabel[i], label)
			}
			lastLabel[i] = label
		}
	}
	return dag
}

func GenerateRandomGridWalkDag2(rng *rand.Rand, num int, numWalks int) (*RoomDag, map[string]tile.TilePosition) {
	dag := NewRoomDag()
	roomPos := make(map[string]tile.TilePosition)

	pos := make([]tile.TilePosition, numWalks)
	lastLabel := make([]string, numWalks)

	for len(dag.Nodes) < num {
		for i := range pos {
			label := fmt.Sprintf("%d_%d", pos[i].X, pos[i].Y)
			_, exists := dag.NodeMap[label]
			if !exists {
				dag.AddNode(label)
				roomPos[label] = pos[i]
			}

			dir := rng.Intn(4)
			if dir == 0 {
				pos[i].X++
			} else if dir == 1 {
				pos[i].X--
			} else if dir == 2 {
				pos[i].Y++
			} else {
				pos[i].Y--
			}

			if lastLabel[i] != "" {
				dag.AddEdge(lastLabel[i], label)
			}
			lastLabel[i] = label
		}
	}
	return dag, roomPos
}

// func GenerateRoomDag(num int, edgeProbability float64, maxEdges int) *RoomDag {
// 	dag := NewRoomDag()
// 	for i := 0; i < num; i++ {
// 		dag.AddNode(fmt.Sprintf("%d", i))
// 	}

// 	for i := 0; i < num; i++ {
// 		for j := i+1; j < num; j++ {
// 			numIEdges := len(dag.Edges[fmt.Sprintf("%d", i)])
// 			if numIEdges >= maxEdges { break }

// 			if rand.Float64() < edgeProbability {
// 				dag.AddEdge(fmt.Sprintf("%d", i), fmt.Sprintf("%d", j))
// 			}
// 		}

// 	}
// 	return dag
// }

type RoomPlacement struct {
	// MapDef *MapDefinition // TODO: generic
	Rect tile.Rect
	GoalGap float64
	Static bool
	Repel float64
	Attract float64
	Depth int
	Mass float64
	Placed bool
}

type ForceBasedRelaxer struct {
	// targetGap int
	rng *rand.Rand
	dag *RoomDag
	placements map[string]RoomPlacement
	repel float64
	grav float64
}

func NewForceBasedRelaxer(rng *rand.Rand, dag *RoomDag, placements map[string]RoomPlacement, startingRepel float64, startingGrav float64) *ForceBasedRelaxer {
	return &ForceBasedRelaxer{
		rng: rng,
		dag: dag,
		placements: placements,
		repel: startingRepel,
		grav: startingGrav,
	}
}

func (r *ForceBasedRelaxer) Iterate() {
	for node, edges := range r.dag.Edges {
		for _, e := range edges {
			if HasEdgeIntersections(r.dag, r.placements, node, e) {
				room := r.placements[node]
				room.Attract = 1.2
				room.Repel += 1
				r.placements[node] = room
			} else {
				// r := r.placements[node]
				// r.Attract = 1.0
				// r.Repel = float64(r.Rect.W() * r.Rect.H())
				// rooms[node] = r
			}
		}
	}

	for label, p := range r.placements {
		if HasRectIntersections(r.placements, label) {
			// p.GoalGap++
			p.Repel += 1
			r.placements[label] = p
		} else {
			// p.GoalGap--
			// rooms[label] = p
		}
	}

	r.repel--
	r.grav = 0.1
	Wiggle(r.dag, r.placements, r.repel, r.grav, 1)
}

// Returns true if there were no intersections
func PSLDStep(rng *rand.Rand, dag *RoomDag, place map[string]RoomPlacement) bool {
	for node, edges := range dag.Edges {
		for _, e := range edges {
			if HasEdgeIntersections(dag, place, node, e) {
				// Move
				p1, ok := place[node]
				if !ok { panic("AAA") }
				if p1.Static { continue }

				p2, ok := place[e]
				if !ok { panic("AAA") }

				// We will find the vector from p1 to p2, then move both of them by that amount
				delta := p2.Rect.Center().Sub(p1.Rect.Center())
				p1.Rect = p1.Rect.Moved(delta)
				// p2.Rect = p2.Rect.Moved(delta)

				place[node] = p1
				place[e] = p2

				return false
			}
		}
	}
	return true
}

func PSLD(rng *rand.Rand, dag *RoomDag, place map[string]RoomPlacement) {
	for {
		noIntersections := PSLDStep(rng, dag, place)

		if noIntersections { break }
	}
}

type GridLayout struct {
	rng *rand.Rand
	queue *ds.Queue[string]
	dag *RoomDag
	place map[string]RoomPlacement
	nodeList []string
	topoIdx int
	topoMoved map[string]bool
	tolerance int
}
func NewGridLayout(rng *rand.Rand, dag *RoomDag, place map[string]RoomPlacement, start string, tolerance int) *GridLayout {
	queue := ds.NewQueue[string](len(place)) // TODO: dynamically resized
	queue.Add(start)

	return &GridLayout{
		rng: rng,
		queue: queue,
		dag: dag,
		place: place,
		nodeList: dag.TopologicalSort(start),

		topoMoved: make(map[string]bool),
		tolerance: tolerance,
	}
}

func (l *GridLayout) LayoutGrid(roomPos map[string]tile.TilePosition) {
	gridSize := 1
	for {
		for k, room := range l.place {
			pos, ok := roomPos[k]
			if !ok { panic("AAA") }

			rect := room.Rect.WithCenter(tile.TilePosition{pos.X * gridSize, pos.Y * gridSize})
			room.Rect = rect
			l.place[k] = room
		}

		if !HasAnyRectIntersections(l.place) {
			break
		}

		gridSize++
	}
}

func (l *GridLayout) Expand() bool {
	// Idea: try and move expand everything outwards until you have no more collisions
	// Goal: to minimize edge lengths and keep things generally grid-like (ie prioritize non diagonals)

	for _, parentNode := range l.nodeList {
		if !HasRectIntersections(l.place, parentNode) {
			continue
		}

		parent, ok := l.place[parentNode]
		if !ok { panic("AAA") }

		pos := parent.Rect.Center()
		pos.X = -pos.X
		pos.Y = -pos.Y
		parent.Rect = l.MoveTowards(parent.Rect, pos, 1)
		l.place[parentNode] = parent
	}

	return false
}

func (l *GridLayout) IterateGravity() bool {
	// Idea: try and move everything in by a small amount, if you start overlapping, then move back
	// Goal: to minimize edge lengths and keep things generally grid-like (ie prioritize non diagonals)
	moved := make(map[string]bool)

	done := true
	for _, parentNode := range l.nodeList {
		// parent, ok := l.place[parentNode]
		// if !ok { panic("AAA") }

		edges, ok := l.dag.Edges[parentNode]
		if !ok { panic("AAA") }
		for _, child := range edges {
			alreadyMoved := moved[child]
			if alreadyMoved { continue }
			moved[child] = true

			p, ok := l.place[child]
			if !ok { panic("AAA") }

			origCenter := p.Rect.Center()
			// pos := parent.Rect.Center()
			p.Rect = l.MoveTowards(p.Rect, origCenter, 1)
			l.place[child] = p
			// if HasRectIntersections(l.place, child) || AnyEdgesIntersect(l.dag, l.place) {
			if HasRectIntersections(l.place, child) || !l.AllEdgesAxisAligned(l.tolerance) {
				p.Rect = p.Rect.WithCenter(origCenter)
				l.place[child] = p
			} else {
				done = false // If we ever moved something then we aren't done
			}
		}
	}

	return done
}

//TODO: tolerance?
func (l *GridLayout) AllEdgesAxisAligned(tolerance int) bool {
	for _, n := range l.nodeList {
		edges, ok := l.dag.Edges[n]
		if !ok { panic("AAA") }

		a := l.place[n]
		r1 := a.Rect
		for _, e := range edges {
			b := l.place[e]
			r2 := b.Rect

			if !AxisAligned(tolerance, r1, r2) {
				return false
			}
		}
	}
	return true
}

func AxisAligned(tol int, a, b tile.Rect) bool {
	xOverlap := (tol <= b.Max.X-a.Min.X) && (a.Max.X-b.Min.X >= tol)
	if xOverlap { return true }
	yOverlap := (tol <= b.Max.Y-a.Min.Y) && (a.Max.Y-b.Min.Y >= tol)
	if yOverlap { return true }
	return false
}

func (l *GridLayout) IterateTowardsParent() bool {
	// Idea: try and move everything in by a small amount, if you start overlapping, then move back
	// Goal: to minimize edge lengths and keep things generally grid-like (ie prioritize non diagonals)
	// TODO: Try moving halfway each time?
	moved := make(map[string]bool)
	movedAmount := make(map[string]tile.TilePosition) // track the movement amount of each node so they can compound

	done := true
	for _, parentNode := range l.nodeList {
		parent, ok := l.place[parentNode]
		if !ok { panic("AAA") }

		// parentMove, ok := movedAmount[parentNode]
		// if !ok { panic("AAA") }

		edges, ok := l.dag.Edges[parentNode]
		if !ok { panic("AAA") }
		for _, child := range edges {
			alreadyMoved := moved[child]
			if alreadyMoved { continue }
			moved[child] = true

			p, ok := l.place[child]
			if !ok { panic("AAA") }

			origCenter := p.Rect.Center()
			// p.Rect = l.MoveTowards(p.Rect, tile.TilePosition{}, 1)
			// l.place[child] = p
			// if HasRectIntersections(l.place, child) || AnyEdgesIntersect(l.dag, l.place) {
			// 	p.Rect = p.Rect.WithCenter(origCenter)
			// 	l.place[child] = p
			// }


			// Move x direction
			{
				origCenter = p.Rect.Center()
				relCenter := origCenter.Sub(parent.Rect.Center()) // Relative to parent
				relCenter.Y = 0
				p.Rect = l.MoveTowards(p.Rect, relCenter, 1)
				// p.Rect = p.Rect.Moved(parentMove)
				l.place[child] = p
				// if HasRectIntersections(l.place, child) || AnyEdgesIntersect(l.dag, l.place) {
				if HasRectIntersections(l.place, child) || !l.AllEdgesAxisAligned(l.tolerance) {
					p.Rect = p.Rect.WithCenter(origCenter)
					l.place[child] = p
					movedAmount[child] = tile.TilePosition{}
				} else {
					done = false
				}
				movedAmount[child] = p.Rect.Center().Sub(origCenter)
			}

			// Move y direction
			{
				origCenter = p.Rect.Center()
				relCenter := origCenter.Sub(parent.Rect.Center()) // Relative to parent
				relCenter.X = 0
				p.Rect = l.MoveTowards(p.Rect, relCenter, 1)
				// p.Rect = p.Rect.Moved(parentMove)
				l.place[child] = p
				// if HasRectIntersections(l.place, child) || AnyEdgesIntersect(l.dag, l.place) {
				if HasRectIntersections(l.place, child) || !l.AllEdgesAxisAligned(l.tolerance) {
					p.Rect = p.Rect.WithCenter(origCenter)
					l.place[child] = p
					movedAmount[child] = tile.TilePosition{}
				} else {
					done = false
				}
				movedAmount[child] = p.Rect.Center().Sub(origCenter)
			}
		}
	}

	return done
}

func (l *GridLayout) Iterate() bool {
	// Idea: Pick a random node and optimize all of its edges,
	// Goal: to minimize edge lengths and keep things generally grid-like (ie prioritize non diagonals)
	// TODO: maybe pick longest edge?

	// rngIndex := l.rng.Intn(len(l.nodeList))
	parentNode := l.nodeList[l.topoIdx]
	parent, ok := l.place[parentNode]
	if !ok { panic("AAA") }

	edges, ok := l.dag.Edges[parentNode]
	if !ok { panic("AAA") }
	for _, child := range edges {
		alreadyMoved := l.topoMoved[child]
		if alreadyMoved { continue }

		p, ok := l.place[child]
		if !ok { panic("AAA") }

		for {
			l.topoMoved[child] = true

			origCenter := p.Rect.Center()
			relCenter := origCenter.Sub(parent.Rect.Center()) // Relative to parent

			// fmt.Println(p.Rect, origCenter, relCenter)
			p.Rect = l.MoveTowards(p.Rect, relCenter, 1)
			// fmt.Println(p.Rect)
			l.place[child] = p
			if HasRectIntersections(l.place, child) || AnyEdgesIntersect(l.dag, l.place) {
				// fmt.Println("INTERSECT")
				p.Rect = p.Rect.WithCenter(origCenter)
				l.place[child] = p
				break
			}
		}
	}

	l.topoIdx = (l.topoIdx + 1) % len(l.nodeList)
	return false
}

func (l *GridLayout) MoveTowards(rect tile.Rect, pos tile.TilePosition, dist int) tile.Rect {
	move := tile.TilePosition{}
	magX := pos.X
	if magX < 0 { magX = -magX }
	magY := pos.Y
	if magY < 0 { magY = -magY }

	if magX > magY {
		if pos.X > 0 {
			move.X = -1
		} else if pos.X < 0 {
			move.X = +1
		}
	} else {
		if pos.Y > 0 {
			move.Y = -1
		} else if pos.Y < 0 {
			move.Y = +1
		}
	}

	return rect.Moved(move)
}

// type HeirarchicalLayout struct {
// 	rng *rand.Rand
// 	queue *ds.Queue[string]
// 	dag *RoomDag
// 	place map[string]RoomPlacement

// 	depthGroup [][]string
// }
// func NewHeirarchicalLayout(rng *rand.Rand, dag *RoomDag, place map[string]RoomPlacement, start string) *HeirarchicalLayout {
// 	queue := ds.NewQueue[string](len(place)) // TODO: dynamically resized
// 	queue.Add(start)

// 	return &HeirarchicalLayout{
// 		rng: rng,
// 		queue: queue,
// 		dag: dag,
// 		place: place,
// 	}
// }
// func (l *HeirarchicalLayout) GroupLayers(start string) {
// 	visit := make(map[string]bool)
// 	queue := ds.NewQueue[string](len(place)) // TODO: dynamically resized
// 	queue.Add(start)

// 	for {
// 		currentNode, ok := queue.Remove()
// 		if !ok { break }

// 		currentPlacement, ok := d.place[currentNode]
// 		if !ok { panic("MUST BE SET") }
// 		currentDepth := currentPlacement.Depth

// 		children := d.dag.Edges[currentNode]
// 		for _, child := range children {
// 			visited := visit[child]
// 			if visited { continue }
// 			visit[child] = true

// 			// Set the depth of the child
// 			room := d.place[child]
// 			room.Depth = currentDepth + 1
// 			d.place[child] = room

// 			// Add the child to the relevant depthGroup
// 			if room.Depth >= len(d.depthGroup) {
// 				d.depthGroup = append(d.depthGroup, make([]string, 0))
// 			}
// 			group := d.depthGroup[room.Depth]
// 			group = append(group, child)
// 			d.depthGroup[room.Depth] = group

// 			queue.Add(child)
// 		}
// 	}
// }

// func (l *HeirarchicalLayout) Iterate() bool {
// 	for _, label := range l.depthGroup {
		
// 	}
// }

type DepthFirstLayout struct {
	rng *rand.Rand
	stack *ds.Stack[string]
	dag *RoomDag
	place map[string]RoomPlacement
	dist float64
}
func NewDepthFirstLayout(rng *rand.Rand, dag *RoomDag, place map[string]RoomPlacement, start string, distance float64) *DepthFirstLayout {
	stack := ds.NewStack[string]()
	stack.Add(start)

	return &DepthFirstLayout{
		rng: rng,
		stack: stack,
		dag: dag,
		place: place,
		dist: distance,
	}
}
func (d *DepthFirstLayout) Reset() {
	for k, p := range d.place {
		if !p.Static {
			p.Placed = false
			d.place[k] = p
		}
	}
}
func (d *DepthFirstLayout) CutCrossed() bool {
	cutList := make([]string, 0)
	for k, _ := range d.place {
		if NodeHasEdgeIntersections(d.dag, d.place, k) {
		// if NodeHasEdgeIntersectionsWithMoreShallowEdge(d.dag, d.place, k) {
			cutList = append(cutList, k)
		}
	}

	if len(cutList) <= 0 {
		return false
	}

	for _, k := range cutList {
		d.CutBelow(k)
		// d.Cut(k)

		d.stack.Add(k)
	}

	// // TODO: Add them all. a bit hacky
	// for k := range d.place {
	// 	d.stack.Add(k)
	// }
	return true
}

func (d *DepthFirstLayout) Cut(label string) {
	p, ok := d.place[label]
	if !ok { return }

	// fmt.Println("Cut: ", label)
	p.Placed = false
	d.place[label] = p
}

func (d *DepthFirstLayout) CutBelow(label string) {
	// fmt.Println("CutBelow: ", label)

	visit := make(map[string]bool)

	stack := ds.NewStack[string]()
	stack.Add(label)
	for {
		curr, ok := stack.Remove()
		if !ok { break }
		visited := visit[curr]
		if visited { continue }
		visit[curr] = true

		children := d.dag.Edges[curr]
		for _, child := range children {
			// fmt.Println("Cut: ", curr)
			p, ok := d.place[child]
			if !ok { continue }
			p.Placed = false
			d.place[child] = p

			stack.Add(child)
		}
	}
}

// Returns true to indicate it hasn't finished
func (d *DepthFirstLayout) Iterate() bool {
	currentNode, ok := d.stack.Remove()
	if !ok { return false }

	children := d.dag.Edges[currentNode]
	for _, child := range children {
		room := d.place[child]
		if room.Placed { continue }

		currentPlacement, ok := d.place[currentNode]
		if !ok { panic("MUST BE SET") }
		currentRect := currentPlacement.Rect
		currentDepth := currentPlacement.Depth

		for i := 0; i < 50; i++ {
			// // Average Position
			// squareRadius := int(d.dist)
			// pos := FindNodeNeighborAveragePosition(d.dag, d.place, child)
			// moveX := pos.X + d.rng.Intn(2 * squareRadius) - squareRadius
			// moveY := pos.Y + d.rng.Intn(2 * squareRadius) - squareRadius
			// room.Rect = room.Rect.WithCenter(tile.TilePosition{moveX, moveY})
			// room.Depth = currentDepth + 1

			// Random square placement
			squareRadius := int(d.dist)
			moveX := currentRect.Center().X + d.rng.Intn(2 * squareRadius) - squareRadius
			moveY := currentRect.Center().Y + d.rng.Intn(2 * squareRadius) - squareRadius
			room.Rect = room.Rect.WithCenter(tile.TilePosition{moveX, moveY})
			room.Depth = currentDepth + 1

			room.Placed = true
			d.place[child] = room

			repel := 90000.0
			grav := 0.0
			Wiggle(d.dag, d.place, repel, grav, 10)

			if !NodeHasEdgeIntersections(d.dag, d.place, child) {
				break
			}
		}

		d.stack.Add(child)
	}
	return true
}

type BreadthFirstLayout struct {
	rng *rand.Rand
	queue *ds.Queue[string]
	dag *RoomDag
	place map[string]RoomPlacement
	dist float64
}
func NewBreadthFirstLayout(rng *rand.Rand, dag *RoomDag, place map[string]RoomPlacement, start string, distance float64) *BreadthFirstLayout {
	queue := ds.NewQueue[string](len(place)) // TODO: dynamically resized
	queue.Add(start)

	return &BreadthFirstLayout{
		rng: rng,
		queue: queue,
		dag: dag,
		place: place,
		dist: distance,
	}
}

func (d *BreadthFirstLayout) Reset() {
	for k, p := range d.place {
		if !p.Static {
			p.Placed = false
			d.place[k] = p
		}
	}
}

// Returns true to indicate it hasn't finished
func (d *BreadthFirstLayout) Iterate() bool {
	currentNode, ok := d.queue.Remove()
	if !ok { return false }

	children := d.dag.Edges[currentNode]
	for _, child := range children {
		room := d.place[child]
		if room.Placed { continue }

		currentPlacement, ok := d.place[currentNode]
		if !ok { panic("MUST BE SET") }
		currentRect := currentPlacement.Rect
		currentDepth := currentPlacement.Depth

		for i := 0; i < 50; i++ {
			// Random square placement
			squareRadius := int(d.dist)
			moveX := currentRect.Center().X + d.rng.Intn(2 * squareRadius) - squareRadius
			moveY := currentRect.Center().Y + d.rng.Intn(2 * squareRadius) - squareRadius
			room.Rect = room.Rect.WithCenter(tile.TilePosition{moveX, moveY})
			room.Depth = currentDepth + 1

			room.Placed = true

			repel := 50000.0
			grav := 0.0
			Wiggle(d.dag, d.place, repel, grav, 20)

			d.place[child] = room

			if !NodeHasEdgeIntersections(d.dag, d.place, child) {
				break
			}
		}

		d.queue.Add(child)
	}
	return true
}


func PlaceDepthFirst(rng *rand.Rand, dag *RoomDag, rooms map[string]RoomPlacement, start string, distance float64) {
	stack := ds.NewStack[string]()
	stack.Add(start)

	fan := 0
	fanX := []int{0, 1, 1, 1}
	fanY := []int{1, 1, 1, -1}
	placed := make(map[string]bool)
	for {
		currentNode, ok := stack.Remove()
		if !ok { break }

		children := dag.Edges[currentNode]
		for _, child := range children {
			_, alreadyPlaced := placed[child]
			if alreadyPlaced { continue }

			currentPlacement, ok := rooms[currentNode]
			if !ok { panic("MUST BE SET") }
			currentRect := currentPlacement.Rect
			currentDepth := currentPlacement.Depth

			squareRadius := int(distance)
			// squareRadius := int(distance / float64(2 * (currentDepth + 1)))
			// moveX := currentRect.Center().X + rng.Intn(2 * squareRadius) - squareRadius
			// moveY := currentRect.Center().Y + rng.Intn(2 * squareRadius) - squareRadius
			fan = (fan + 1) % len(fanX)
			mX := fanX[fan]
			mY := fanY[fan]
			moveX := (mX * squareRadius) + currentRect.Center().X
			moveY := (mY * squareRadius) + currentRect.Center().Y

			room := rooms[child]
			room.Rect = room.Rect.WithCenter(tile.TilePosition{moveX, moveY})
			room.Depth = currentDepth + 1
			// TODO: Recalc gap goal?
			rooms[child] = room

			// // Random square placement
			// // TODO: hardcoded. Could I favor non-diagonals here?
			// squareRadius := int(distance / float64(2 * (currentDepth + 1)))
			// // moveX := currentRect.Center().X + rng.Intn(2 * squareRadius) - squareRadius
			// // moveY := currentRect.Center().Y + rng.Intn(2 * squareRadius) - squareRadius
			// mX := rng.Intn(3) - 1
			// mY := rng.Intn(3) - 1
			// moveX := (mX * squareRadius) + currentRect.Center().X
			// moveY := (mY * squareRadius) + currentRect.Center().Y

			// room := rooms[child]
			// room.Rect = room.Rect.WithCenter(tile.TilePosition{moveX, moveY})
			// room.Depth = currentDepth + 1
			// // TODO: Recalc gap goal?
			// rooms[child] = room

			// // Random circle placement
			// rngAngle := rng.Float64() * math.Pi / 4
			// posAngle := currentRect.Center().Angle()
			// moveX := int(distance * math.Cos(rngAngle)) + currentRect.Center().X
			// moveY := int(distance * math.Sin(rngAngle)) + currentRect.Center().Y
			// room := rooms[child]
			// room.Rect = room.Rect.WithCenter(tile.TilePosition{moveX, moveY})
			// room.Depth = currentDepth + 1
			// fmt.Println(child, tile.TilePosition{moveX, moveY}, room.Rect, room.Depth)
			// // TODO: Recalc gap goal?
			// rooms[child] = room

			repel := 1000.0
			grav := 0.0
			Wiggle(dag, rooms, repel, grav, 100)

			placed[child] = true
			stack.Add(child)
		}
	}
}

func PlaceBreadthFirst(rng *rand.Rand, dag *RoomDag, rooms map[string]RoomPlacement, start string, distance float64) {
	childQueue := ds.NewQueue[string](len(rooms)) // TODO: probably doesnt need to be this big, would be nice if we could have a dynamically resizable queue
	childQueue.Add(start)

	placed := make(map[string]bool)
	for {
		currentNode, ok := childQueue.Remove()
		if !ok { break }
		// currentPlacement, ok := rooms[currentNode]
		// if !ok { panic("MUST BE SET") }
		// currentRect := currentPlacement.Rect
		// currentDepth := currentPlacement.Depth

		children := dag.Edges[currentNode]
		for _, child := range children {
			_, alreadyPlaced := placed[child]
			if alreadyPlaced { continue }

			currentPlacement, ok := rooms[currentNode]
			if !ok { panic("MUST BE SET") }
			currentRect := currentPlacement.Rect
			currentDepth := currentPlacement.Depth

			// Random square placement
			squareRadius := int(distance) // TODO: hardcoded. Could I favor non-diagonals here?
			// moveX := currentRect.Center().X + rng.Intn(2 * squareRadius) - squareRadius
			// moveY := currentRect.Center().Y + rng.Intn(2 * squareRadius) - squareRadius
			mX := rng.Intn(3) - 1
			mY := rng.Intn(3) - 1
			moveX := (mX * squareRadius) + currentRect.Center().X
			moveY := (mY * squareRadius) + currentRect.Center().Y

			room := rooms[child]
			room.Rect = room.Rect.WithCenter(tile.TilePosition{moveX, moveY})
			room.Depth = currentDepth + 1
			// TODO: Recalc gap goal?
			rooms[child] = room

			// // Random circle placement
			// rngAngle := rng.Float64() * math.Pi / 4
			// posAngle := currentRect.Center().Angle()
			// moveX := int(distance * math.Cos(rngAngle)) + currentRect.Center().X
			// moveY := int(distance * math.Sin(rngAngle)) + currentRect.Center().Y
			// room := rooms[child]
			// room.Rect = room.Rect.WithCenter(tile.TilePosition{moveX, moveY})
			// room.Depth = currentDepth + 1
			// fmt.Println(child, tile.TilePosition{moveX, moveY}, room.Rect, room.Depth)
			// // TODO: Recalc gap goal?
			// rooms[child] = room

			placed[child] = true
			childQueue.Add(child)

			// for i := 0; i < 1; i++ {
			// 	PSLDStep(rng, dag, rooms)
			// }

			// repel := 1000.0
			// grav := 0.0
			// Wiggle(dag, rooms, repel, grav, 100)
		}

		// // All children are added, wiggle everything
		// repel := 100.0
		// grav := 1.0
		// Wiggle(dag, rooms, repel, grav, 10)

		// bailoutMax--
		// if bailoutMax <= 0 { break }
	}
}

func Untangle(dag *RoomDag, rooms map[string]RoomPlacement, repelMultiplier, gravityConstant float64, iterations int) {
	for i := 0; i < iterations; i++ {
		for node, edges := range dag.Edges {
			for _, e := range edges {
				if HasEdgeIntersections(dag, rooms, node, e) {
					r := rooms[node]
					r.Attract = 1.2
					r.Repel += 1
					rooms[node] = r

					// noIntersections = false
				} else {
					r := rooms[node]
					// r.Attract = 1.0
					// r.Repel = float64(r.Rect.W() * r.Rect.H())
					rooms[node] = r
				}
			}
		}
		stable := Wiggle(dag, rooms, repelMultiplier, gravityConstant, 1)
		if stable {
			break
		}
	}
}

func Wiggle(dag *RoomDag, rooms map[string]RoomPlacement, repelMultiplier, gravityConstant float64, iterations int) bool {
	// TODO: everything inside of here probably needs to be fixed point, or integer based
	forces := make(map[string]phy2.Vec)

	// repelMultiplier = 100.0
	// repelConstant := 50000.0
	// gravityConstant := 0.0001
	// gravityConstant := 0.0
	attractConstant := 0.2
	mass := 4.0

	stable := true
	// TODO: You gotta make this loop deterministic
	// https://editor.p5js.org/JeromePaddick/sketches/bjA_UOPip
	for i := 0; i < iterations; i++ {
		// Calculate Forces
		for key, placement := range rooms {
			if placement.Static { continue } // Dont add forces to static placements
			if !placement.Placed { continue } // Skip if not yet placed

			force := phy2.Vec{}
			for key2, placement2 := range rooms {
				if key == key2 { continue } // Skip if same
				if !placement2.Placed { continue } // Skip if not yet placed

				deltaTile := placement2.Rect.Center().Sub(placement.Rect.Center())

				delta := phy2.Vec{float64(deltaTile.X), float64(deltaTile.Y)}

				// dist := float64(tile.ManhattanDistance(placement.Rect.Center(), placement2.Rect.Center()))
				dist := delta.Len()

				// Connections
				if dag.HasEdgeEitherDirection(key2, key) {
					attractVal := attractConstant * placement.Attract * placement2.Attract
					mag := attractVal * (dist - (placement.GoalGap + placement2.GoalGap))
					if mag > 100 { mag = 100 }
					vec := delta.Norm().Scaled(mag)
					vec.X *= 4
					force = force.Add(vec)
					// fmt.Println("Conn: ", vec)
				} else {
					if dist < 100 { dist = 100 }
					// Repulsions
					// mag := -1 * repelConstant / (dist * dist)
					repelVal := repelMultiplier * placement.Repel * placement2.Repel
					mag := -1 * repelVal / (dist * dist * dist * dist)
					vec := delta.Norm().Scaled(mag)
					force = force.Add(vec)
					// fmt.Println("Repel: ", vec)
				}
			}

			// Add a gravity
			pos := phy2.Vec{
				float64(placement.Rect.Center().X),
				float64(placement.Rect.Center().Y),
			}
			// force = force.Add(pos.Scaled(-gravityConstant / (10 * float64(placement.Depth + 1))))
			force = force.Add(pos.Scaled(-gravityConstant))
			// fmt.Println("Grav: ", pos.Scaled(-gravityConstant))

			// // Add very weak inverse gravity
			// {
			// 	ceil := 1000.0
			// 	mag := pos.Len()
			// 	if mag > ceil { mag = ceil }
			// 	force = force.Add(pos.Norm().Scaled(mag))
			// }

			// // Add a upward Y force
			// {
			// 	vec := phy2.Vec{
			// 		0.0,
			// 		-100.0,
			// 	}
			// 	force = force.Add(vec)
			// }

			forces[key] = force
		}

		// Apply Forces
		for k, f := range forces {
			placement, ok := rooms[k]
			if !ok { panic("AAAA") }

			tileMove := tile.TilePosition{
				int(math.Round(f.X / mass)),
				int(math.Round(f.Y / mass)),
			}

			if (tileMove == tile.TilePosition{}) { continue } // Skip if there is no movement

			stable = false
			placement.Rect = placement.Rect.Moved(tileMove)
			rooms[k] = placement
		}
	}
	return stable
}

// func (g roomAndHallwayDungeon) pickRoom(rooms []string) *MapDefinition {
// 	// Random block
// 	idx := g.rng.Intn(len(rooms))
// 	name := rooms[idx]

// 	mapDef, err := LoadMap(name)
// 	if err != nil {
// 		panic(err)
// 	}
// 	mapDef.RecalculateBounds() // TODO: Move this to editor so we dont have to do it during dungeon generation time
// 	return mapDef
// }

// func (g roomAndHallwayDungeon) findNonCollidingPosition(rects map[string]RoomPlacement, last tile.Rect, newRect tile.Rect) tile.Rect {
// 	expansion := 12
// 	combinedWidth := (last.W() + newRect.W()) / 2
// 	combinedHeight := (last.H() + newRect.H()) / 2

// 	for {
// 		newPos := last.Center()

// 		dir := g.rng.Intn(4)
// 		posneg := g.rng.Intn(2)
// 		if dir == 0 {
// 			newPos.X += combinedWidth + expansion
// 			if posneg == 0 {
// 				newPos.Y += combinedHeight / 2
// 			} else {
// 				newPos.Y -= combinedHeight / 2
// 			}
// 		} else if dir == 1 {
// 			newPos.X -= (combinedWidth + expansion)
// 			if posneg == 0 {
// 				newPos.Y += combinedHeight / 2
// 			} else {
// 				newPos.Y -= combinedHeight / 2
// 			}
// 		} else if dir == 2 {
// 			newPos.Y += combinedHeight + expansion
// 			if posneg == 0 {
// 				newPos.X += combinedWidth / 2
// 			} else {
// 				newPos.X -= combinedWidth / 2
// 			}
// 		} else if dir == 3 {
// 			newPos.Y -= (combinedHeight + expansion)
// 			if posneg == 0 {
// 				newPos.X += combinedWidth / 2
// 			} else {
// 				newPos.X -= combinedWidth / 2
// 			}
// 		}

// 		rect := newRect.WithCenter(newPos)

// 		if !hasRectIntersections(rects, rect) {
// 			return rect
// 		}
// 	}
// }

func NodeHasEdgeIntersectionsWithMoreShallowEdge(dag *RoomDag, placements map[string]RoomPlacement, label string) bool {
	edges, ok := dag.Edges[label]
	if !ok { panic("AAA") }

	p1, ok := placements[label]
	if !ok { return false }

	for _, e := range edges {
		p2, ok := placements[e]
		if !ok { continue }

		if p1.Depth < p2.Depth { continue } // skip if we are more shallow

		if HasEdgeIntersections(dag, placements, label, e) {
			return true
		}
	}

	return false
}

func AnyEdgesIntersect(dag *RoomDag, rects map[string]RoomPlacement) bool {
	for node, edges := range dag.Edges {
		for _, e := range edges {
			if HasEdgeIntersections(dag, rects, node, e) {
				return true
			}
		}
	}
	return false
}

func NodeHasEdgeIntersections(dag *RoomDag, rects map[string]RoomPlacement, label string) bool {
	edges, ok := dag.Edges[label]
	if !ok { panic("AAA") }

	for _, e := range edges {
		if HasEdgeIntersections(dag, rects, label, e) {
			return true
		}
	}

	return false
}

// def ccw(A,B,C):
//     return (C.y-A.y) * (B.x-A.x) > (B.y-A.y) * (C.x-A.x)

// # Return true if line segments AB and CD intersect
// def intersect(A,B,C,D):
//     return ccw(A,C,D) != ccw(B,C,D) and ccw(A,B,C) != ccw(A,B,D)
// https://stackoverflow.com/questions/3838329/how-can-i-check-if-two-segments-intersect
func ccw(a, b, c tile.TilePosition) bool {
	return (c.Y - a.Y) * (b.X - a.X) >= (b.Y - a.Y) * (c.X - a.X)
}

// Returns true if line ab intersects line cd
func intersect(a, b, c, d tile.TilePosition) bool {
	return (ccw(a,c,d) != ccw(b,c,d)) && (ccw(a,b,c) != ccw(a,b,d))
}

func HasEdgeIntersections(dag *RoomDag, rects map[string]RoomPlacement, label1, label2 string) bool {
	// Line ab is going to go from label1 to label2
	rA, ok := rects[label1]
	if !ok { return false }
	rB, ok := rects[label2]
	if !ok { return false }

	if !rA.Placed { return false } // Skip if not placed
	if !rB.Placed { return false } // Skip if not placed


	a := rA.Rect.Center()
	b := rB.Rect.Center()

	// We will skip all edges that are connect to labels that are the same as ours
	for k, edges := range dag.Edges {
		if k == label1 { continue }
		if k == label2 { continue }

		for _, e := range edges {
			if e == label1 { continue }
			if e == label2 { continue }

			rC, ok := rects[k]
			if !ok { continue }
			rD, ok := rects[e]
			if !ok { continue }
			if !rC.Placed { continue } // Skip if not placed
			if !rD.Placed { continue } // Skip if not placed

			c := rC.Rect.Center()
			d := rD.Rect.Center()
			if intersect(a, b, c, d) {
				return true
			}
		}
	}
	return false
}

func HasAnyRectIntersections(rects map[string]RoomPlacement) bool {
	for label := range rects {
		if HasRectIntersections(rects, label) {
			return true
		}
	}
	return false
}

func HasRectIntersections(rects map[string]RoomPlacement, label string) bool {
	src, ok := rects[label]
	if !ok { return false }

	for key, r := range rects {
		if key == label { continue } // skip self
		if !r.Placed { continue } // skip if not placed

		if r.Rect.Intersects(src.Rect) {
			return true
		}
	}
	return false
}

func FindNodeNeighborAveragePosition(dag *RoomDag, place map[string]RoomPlacement, label string) tile.TilePosition {
	total := tile.TilePosition{}
	count := 0

	for key, p := range place {
		if key == label { continue } // Skip self
		if !p.Placed { continue } // skip unplaced rooms

		if dag.HasEdgeEitherDirection(key, label) {
			total = total.Add(p.Rect.Center())
			count++
		}
	}

	if count == 0 {
		return total
	}
	return total.Div(count)
}
