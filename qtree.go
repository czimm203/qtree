package qtree

import (
	"log"
)

const MAX_NODES int = 6

type Node struct {
    Loc Point
    Val any
}

var EmptyNode Node = Node{}

type Point struct {
    X float64
    Y float64
}

func NewQTree (topLeft Point, bottomRight Point) *QTree { 
    return &QTree{
        node: []Node{},
        topLeft: topLeft,
        bottomRight: bottomRight,
        isLeaf: true,
        nw: nil,
        ne: nil,
        sw: nil,
        se: nil,
    }
}

type QTree struct {
    node []Node
    topLeft Point
    bottomRight Point
    isLeaf bool
    nw *QTree
    ne *QTree
    sw *QTree
    se *QTree
}

func (qt *QTree) Find(p Point) []Node {
    // Walk tree to leaf and return it's contents
    pts := make([]Node,0)
    if !qt.SurroundsPoint(p) {
        log.Println("Point not contained")
        return pts
    } else if qt.SurroundsPoint(p) && !qt.isLeaf {
        switch {
        case qt.nw.SurroundsPoint(p):
            return qt.nw.Find(p)
        case qt.ne.SurroundsPoint(p):
            return qt.ne.Find(p)
        case qt.sw.SurroundsPoint(p):
            return qt.sw.Find(p)
        case qt.se.SurroundsPoint(p):
            return qt.se.Find(p)
        }
    } else {
        pts = append(pts, qt.node...)
    }
    return pts
}

func (qt *QTree) BoundedBy(ul Point, lr Point) bool {
    return (qt.topLeft.X <= lr.X) && (qt.topLeft.Y >= lr.Y) && (qt.bottomRight.X >= ul.X) && (qt.bottomRight.Y <= ul.Y)
}

func (qt *QTree) FindInBoundingBox(ul Point, lr Point) []Node {
    pts := make([]Node, 0)

    if !qt.BoundedBy(ul, lr) {
        log.Println("Point not contained")
        return pts
    } else if qt.isLeaf {
        for i := range qt.node {
            if SurroundsPoint(qt.node[i].Loc, ul, lr) {
                pts = append(pts, qt.node[i])
                if(len(pts)) > MAX_NODES {
                    log.Panic("too many nodes")
                }
            }
        }
    } else {
        if qt.nw.BoundedBy(ul, lr) {
            for _, n := range qt.nw.FindInBoundingBox(ul, lr){
                pts = append(pts, n)
            }
        }
        if qt.ne.BoundedBy(ul, lr) {
            for _, n := range qt.ne.FindInBoundingBox(ul, lr){
                pts = append(pts, n)
            }
        }
        if qt.sw.BoundedBy(ul, lr) {
            for _, n := range qt.sw.FindInBoundingBox(ul, lr){
                pts = append(pts, n)
            }
        }
        if qt.se.BoundedBy(ul, lr) {
            for _, n := range qt.se.FindInBoundingBox(ul, lr){
                pts = append(pts, n)
            }
        }
    } 
    return pts
}

func (qt *QTree) Insert(n Node) { 
    if !qt.SurroundsPoint(n.Loc) {
        log.Fatalf("Node out of bounds: %v\n", n)
    }

    if qt.isLeaf && len(qt.node) < MAX_NODES {
        qt.node = append(qt.node, n)
    } else if qt.isLeaf {
        qt.divideQuad()
        qt.Insert(n)
    } else {
        if qt.nw.SurroundsPoint(n.Loc) {
            qt.nw.Insert(n)
        } else if qt.ne.SurroundsPoint(n.Loc) {
            qt.ne.Insert(n)
        } else if qt.sw.SurroundsPoint(n.Loc) {
            qt.sw.Insert(n)
        } else if qt.se.SurroundsPoint(n.Loc) {
            qt.se.Insert(n)
        }
        if(len(qt.nw.node) > MAX_NODES || len(qt.ne.node) > MAX_NODES ||len(qt.sw.node) > MAX_NODES ||len(qt.se.node) > MAX_NODES) {
            log.Panic("too many nodes: Insert")
        }
    }
    if len(qt.node) > MAX_NODES {
        log.Fatalf("sad\n")
    }   
}

func (qt *QTree) IsLeaf() bool {
    return qt.isLeaf 
}

func (qt *QTree) divideQuad() {
    halfX := (qt.bottomRight.X + qt.topLeft.X)/2
    halfY := (qt.topLeft.Y + qt.bottomRight.Y)/2

    qt.nw = NewQTree(qt.topLeft, Point{X: halfX, Y: halfY})
    qt.ne = NewQTree(Point{X: halfX, Y: qt.topLeft.Y }, Point{X: qt.bottomRight.X, Y: halfY})
    qt.sw = NewQTree(Point{X: qt.topLeft.X, Y: halfY}, Point{X: halfX, Y: qt.bottomRight.Y })
    qt.se = NewQTree(Point{X: halfX, Y: halfY},  qt.bottomRight)

    for _, pt := range qt.node {
        if qt.nw.SurroundsPoint(pt.Loc) {
            qt.nw.Insert(pt)
        } else if qt.ne.SurroundsPoint(pt.Loc) {
            qt.ne.Insert(pt)
        } else if qt.sw.SurroundsPoint(pt.Loc) {
            qt.sw.Insert(pt)
        } else if qt.se.SurroundsPoint(pt.Loc) {
            qt.se.Insert(pt)
        }
    }
    qt.node = []Node{}
    qt.isLeaf = false
}

func SurroundsPoint(p Point, ul Point, lr Point) bool {
    return p.X <= lr.X && p.X >= ul.X && p.Y <= ul.Y && p.Y >= lr.Y
}

func (qt *QTree) SurroundsPoint(p Point) bool {
    return SurroundsPoint(p, qt.topLeft, qt.bottomRight)
}

