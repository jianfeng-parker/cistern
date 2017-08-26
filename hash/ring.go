package hash

import (
	"math"
	"crypto/sha1"
	"strconv"
)

type node struct {
	hash     uint32
	identity string
}

type Ring struct {
	slots   int
	nodes   []node
	weights map[string]int
}

func NewRing(slots int) (r *Ring) {
	if slots <= 0 {
		slots = 160 // 默认160个虚拟节点
	}
	r = &Ring{
		slots: slots,
	}
	return
}

func (r *Ring) Add(identity string, weight int) {

}

func (r *Ring) AddMulti() {

}

func (r *Ring) Remove(identity string) {

}

func (r *Ring) Update(identity string, weight int) {

}

func (r *Ring) Len() int {
	return len(r.nodes)
}

func (r *Ring) Swap() {

}

func (r *Ring) GetNode(key string) string {
	if len(r.nodes) == 0 {
		return ""
	}
	
}

func (r *Ring) build() {
	var totalWeight int
	for _, weight := range r.weights {
		totalWeight += weight
	}
	totalSlots := r.slots * len(r.weights)
	sha := sha1.New()
	for identity, weight := range r.weights {
		// 计算Key对应的虚拟节点数
		slots := int(math.Floor(float64(weight) / float64(totalWeight) * float64(totalSlots)))
		for i := 1; i <= slots; i++ {
			sha.Write([]byte(identity + ":" + strconv.Itoa(i)))
			n := node{
				hash:     getHash(sha.Sum(nil)[2:6]),
				identity: identity,
			}
			r.nodes = append(r.nodes, n)
			sha.Reset()
		}
	}
	sort(r.nodes)
}

func getHash(bs []byte) uint32 {

}
