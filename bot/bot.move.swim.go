package bot

import (
	"cmp"
	"container/heap"
	"math"
)

const (
	swimCell   = 20.0
	swimReplan = 1000
	swimStall  = 5000
)

var swimDirs = [8][2]int{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}}

// swimTo takes the bot through the water of a swim map to the goal node: along the shortest path of water cells,
// steering toward a point a little ahead on it and searching again when the path is old or the bot drifts off it.
// A goal on a floor is reached by sinking onto that floor and walking; a goal on a rope by catching it.
func (w *walker) swimTo(goal int, tx, ty float64) bool {
	end := w.g.nodes[goal]
	aimX, aimY := tx, ty
	if end.chain >= 0 {
		floor, _ := w.g.surface(end.chain, tx)
		aimY = floor - swimCell
	} else {
		aimX = float64(end.rope.X)
	}

	best, bestAt := math.Inf(1), w.spent
	for w.spent <= w.limit {
		switch {
		case end.chain >= 0 && w.fh != 0 && w.g.chainOf[w.fh] == end.chain:
			return w.walkTo(tx) && w.settle()
		case end.chain < 0 && w.rope == end.rope.ID:
			return w.climbTo(ty) && w.settle()
		case w.rope != 0:
			dir := 1
			if aimX < w.x {
				dir = -1
			}
			if !w.step(dir, 0, true) {
				return false
			}
			continue
		}
		if d := math.Hypot(aimX-w.x, aimY-w.y); d < best-5 {
			best, bestAt = d, w.spent
		} else if w.spent-bestAt > swimStall {
			return false
		}

		// Standing on a floor, the search starts just above it so the floor still holds the first step.
		fromY := w.y
		if w.fh != 0 {
			fromY--
		}
		path, ok := w.g.swimPath(w.x, fromY, aimX, aimY)
		if !ok {
			return false
		}
		for planned := w.spent; w.spent-planned < swimReplan; {
			for len(path) > 1 && math.Hypot(path[0][0]-w.x, path[0][1]-w.y) < swimCell*1.5 {
				path = path[1:]
			}
			if !w.stroke(path[0][0], path[0][1]) {
				return false
			}
			if w.rope != 0 || end.chain >= 0 && w.fh != 0 && w.g.chainOf[w.fh] == end.chain {
				break
			}
		}
	}
	return false
}

// stroke steers one tick toward (x, y). Every key but the stroke sinks the bot, up the least, so it holds up unless
// it needs to dive. A stroke lifts it some 160 px, so it strokes only for a point well above and once the last
// stroke has worn off. On a floor it walks or jumps.
func (w *walker) stroke(x, y float64) bool {
	side := 0
	if math.Abs(x-w.x) > 4 {
		side = cmp.Compare(x, w.x)
	}
	vertical, jump := -1, false
	switch {
	case w.fh != 0:
		vertical, jump = 0, y < w.y-swimCell/2
	case y < w.y-swimCell*2:
		jump = w.vy >= 0
	case y > w.y+swimCell/2:
		vertical = 1
	}
	return w.step(side, vertical, jump)
}

// swimPath is the shortest path by A* over water cells from (x0, y0) to (x1, y1), as the centers of its cells after
// the first, ending at (x1, y1). A step between cells may not cross a floor from above, nor a wall or a ceiling from
// its blocking side.
func (g *routeMap) swimPath(x0, y0, x1, y1 float64) ([][2]float64, bool) {
	if g.tree == nil {
		return nil, false
	}
	bounds := g.tree.Bounds
	left, top := float64(bounds.Left)-swimCell, float64(bounds.Top)-swimCell*10
	cols := int((float64(bounds.Right)+swimCell-left)/swimCell) + 1
	rows := int((float64(bounds.Bottom)-top)/swimCell) + 1
	cell := func(x, y float64) int {
		return min(max(int((y-top)/swimCell), 0), rows-1)*cols + min(max(int((x-left)/swimCell), 0), cols-1)
	}
	center := func(n int) (float64, float64) {
		return left + (float64(n%cols)+0.5)*swimCell, top + (float64(n/cols)+0.5)*swimCell
	}

	start, goal := cell(x0, y0), cell(x1, y1)
	cost := map[int]float64{start: 0}
	prev := make(map[int]int)
	open := &swimQueue{{cell: start, f: math.Hypot(x1-x0, y1-y0)}}
	for open.Len() > 0 {
		cur := heap.Pop(open).(swimEntry)
		if cur.cell == goal {
			path := [][2]float64{{x1, y1}}
			for at := goal; at != start; {
				if at = prev[at]; at == start {
					break
				}
				cx, cy := center(at)
				path = append([][2]float64{{cx, cy}}, path...)
			}
			return path, true
		}
		if cur.cost > cost[cur.cell] {
			continue
		}
		cx, cy := center(cur.cell)
		if cur.cell == start {
			cx, cy = x0, y0
		}
		col, row := cur.cell%cols, cur.cell/cols
		for d, dir := range swimDirs {
			c, r := col+dir[0], row+dir[1]
			if c < 0 || c >= cols || r < 0 || r >= rows {
				continue
			}
			next := r*cols + c
			nx, ny := center(next)
			if g.waterCut(cur.cell*8+d, cur.cell == start, cx, cy, nx, ny) {
				continue
			}
			step := cur.cost + math.Hypot(nx-cx, ny-cy)
			if old, ok := cost[next]; ok && old <= step {
				continue
			}
			cost[next], prev[next] = step, cur.cell
			heap.Push(open, swimEntry{cell: next, cost: step, f: step + math.Hypot(x1-nx, y1-ny)})
		}
	}
	return nil, false
}

// waterCut reports a step between water cells that a foothold blocks; steps between cell centers are kept per map,
// a step from the bot's own spot is not.
func (g *routeMap) waterCut(key int, fromBot bool, x0, y0, x1, y1 float64) bool {
	if cut, ok := g.waterCuts[key]; ok && !fromBot {
		return cut
	}
	cut := false
	for _, fh := range g.query(min(x0, x1), min(y0, y1), max(x0, x1), max(y0, y1)) {
		if g.crosses(fh, x0, y0, x1, y1) {
			cut = true
			break
		}
	}
	if !fromBot {
		g.waterCuts[key] = cut
	}
	return cut
}

type swimEntry struct {
	cell    int
	cost, f float64
}

type swimQueue []swimEntry

func (q swimQueue) Len() int           { return len(q) }
func (q swimQueue) Less(i, j int) bool { return q[i].f < q[j].f }
func (q swimQueue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *swimQueue) Push(x any)        { *q = append(*q, x.(swimEntry)) }
func (q *swimQueue) Pop() any {
	old := *q
	e := old[len(old)-1]
	*q = old[:len(old)-1]
	return e
}
