package metric

import (
	"math"
	"sync"
	"time"
)

type latencyMetrics struct {
	mu sync.Mutex

	total uint64

	min time.Duration
	max time.Duration
	avg float64

	percentilesTree *node
}

type node struct {
	left, right *node

	value      time.Duration
	itemsCount int
	count      int
	height     int
}

func NewLatencyMetrics() LatencyMetrics {
	return &latencyMetrics{
		mu:              sync.Mutex{},
		total:           0,
		min:             -1,
		max:             -1,
		avg:             0,
		percentilesTree: nil,
	}
}

func (lm *latencyMetrics) ConsumeResult(res *Result) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.total++
	if lm.min > res.Duration || lm.min == -1 {
		lm.min = res.Duration
	}
	if lm.max < res.Duration || lm.max == -1 {
		lm.max = res.Duration
	}

	prevSum := lm.avg * float64(lm.total-1)
	curSum := prevSum + float64(res.Duration)
	lm.avg = curSum / float64(lm.total)

	if lm.percentilesTree == nil {
		lm.percentilesTree = newNode(res.Duration)
		return
	}
	lm.percentilesTree.insert(res.Duration)
}

func (lm *latencyMetrics) GetPercentile(p int) time.Duration {
	if lm.percentilesTree == nil || lm.percentilesTree.count < p {
		return 0
	}

	return lm.percentilesTree.findPercentile(p)
}

func (lm *latencyMetrics) GetMin() time.Duration {
	return lm.min
}

func (lm *latencyMetrics) GetMax() time.Duration {
	return lm.max
}

func (lm *latencyMetrics) GetAvg() time.Duration {
	return time.Duration(uint64(lm.avg))
}

func newNode(val time.Duration) *node {
	return &node{
		left:       nil,
		right:      nil,
		value:      val,
		count:      1,
		height:     1,
		itemsCount: 1,
	}
}

func (n *node) getHeight() int {
	if n == nil {
		return 0
	}

	return n.height
}

func (n *node) getBalanceFactor() int {
	return n.right.getHeight() - n.left.getHeight()
}

func (n *node) fixHeight() {
	n.height = max(n.left.getHeight(), n.right.getHeight()) + 1
}

func (n *node) getCount() int {
	if n == nil {
		return 0
	}

	return n.count
}

func (n *node) getItemsCount() int {
	if n == nil {
		return 0
	}

	return n.itemsCount
}

func (n *node) recalculate() {
	n.fixHeight()
	n.count = n.itemsCount + n.left.getCount() + n.right.getCount()
}

func (n *node) findPercentile(percentile int) time.Duration {
	if n.count == 1 {
		return n.value
	}

	x := float64(percentile)/float64(100)*float64(n.count-1) + 1
	abs, frac := math.Modf(x)

	return n.getValueByIndex(int(abs), 0) + time.Duration(frac)*(n.getValueByIndex(int(abs)+1, 0)-n.getValueByIndex(int(abs)+1, 0))
}

func (n *node) getValueByIndex(idx, passed int) time.Duration {
	lc := n.left.getItemsCount()

	if passed+lc >= idx {
		return n.left.getValueByIndex(idx, 0)
	}
	if passed+n.itemsCount+lc >= idx {
		return n.value
	}

	return n.right.getValueByIndex(idx, passed+n.itemsCount+lc)
}

func (n *node) insert(v time.Duration) {
	defer func() {
		n.fixHeight()
		n.count++
		n.rebalance()
	}()

	if v > n.value {
		if n.right == nil {
			n.right = newNode(v)
			return
		}

		n.right.itemsCount++
		return
	}
	if v < n.value {
		if n.left == nil {
			n.left = newNode(v)
			return
		}

		n.left.itemsCount++
		return
	}

	n.itemsCount++
}

func (n *node) rebalance() {
	bf := n.getBalanceFactor()

	if bf > 1 && n.left.getBalanceFactor() >= 0 {
		n.rightRotate()
	}
	if bf < -1 && n.right.getBalanceFactor() <= 0 {
		n.leftRotate()
	}
	if bf > 1 && n.left.getBalanceFactor() < 0 {
		n.left.leftRotate()
		n.rightRotate()
	}
	if bf < -1 && n.right.getBalanceFactor() > 0 {
		n.right.rightRotate()
		n.leftRotate()
	}
}

func (n *node) swap(n2 *node) {
	n.itemsCount, n.value, n2.itemsCount, n2.value = n2.itemsCount, n2.value, n.itemsCount, n.value
}

func (n *node) rightRotate() {
	r := n.left
	defer n.recalculate()
	defer r.recalculate()

	n.swap(r)
	n.left = r.left
	r.left = r.right
	r.right = n.right
	n.right = r
}

func (n *node) leftRotate() {
	l := n.right
	defer n.recalculate()
	defer l.recalculate()

	n.swap(l)
	n.right = l.right
	l.right = l.left
	l.left = n.left
	n.left = l
}
