package pgen

import (
	"iter"
	"math/rand"

	"github.com/unitoftime/flow/glm"
	"github.com/unitoftime/flow/spatial"
)

type PoissonPlacement struct {
	HashSize int // The size of the spatial hash

	pointmap *spatial.Pointmap[struct{}]
}

func (p *PoissonPlacement) Add(pos glm.Vec2) {
	p.pointmap.Add(pos, struct{}{})
}

func (p *PoissonPlacement) Check(pos glm.Vec2, bounds glm.Rect) bool {
	return p.pointmap.Collides(bounds.WithCenter(pos))
}

func (p *PoissonPlacement) Reset() {
	p.pointmap = spatial.NewPointmap[struct{}]([2]int{p.HashSize, p.HashSize}, 0)
}

// Returns an iterator that generates all acceptable placement positions
func (p *PoissonPlacement) All(rng *rand.Rand, bounds glm.Rect, attempts int, minDist float64) iter.Seq[glm.Vec2] {
	// A pointmap to ensure we dont clump together certain objects
	if p.pointmap == nil {
		p.pointmap = spatial.NewPointmap[struct{}]([2]int{p.HashSize, p.HashSize}, 0)
	}

	clusterCheckRect := glm.CR(minDist)

	return func(yield func(glm.Vec2) bool) {
		for range attempts {
			checkPos := SeededRandomPositionInRect(rng, bounds)
			checkRect := clusterCheckRect.WithCenter(checkPos)
			if p.pointmap.Collides(checkRect) { continue } // Skip, too close to something else

			if !yield(checkPos) {
				break
			}
		}
	}
}
