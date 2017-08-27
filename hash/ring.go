package hash

import (
	"math"
	"crypto/sha1"
	"strconv"
	"sort"
)

const (
	factor_hash int = 4
)

type node struct {
	hash        uint32
	destination string
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
		slots:   slots,
		weights: make(map[string]int),
	}
	return
}

func (r *Ring) Add(destination string, weight int) {
	r.weights[destination] = weight
	r.build()
}

func (r *Ring) AddMulti(weights map[string]int) {
	for destination, weight := range weights {
		r.weights[destination] = weight
	}
	r.build()
}

func (r *Ring) Remove(destination string) {
	delete(r.weights, destination)
	r.build()
}

func (r *Ring) Update(destination string, weight int) {
	r.weights[destination] = weight
	r.build()
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
	sha := sha1.New()
	sha.Write([]byte(key))
	hBytes := sha.Sum(nil)
	hValue := getHash(hBytes[2:2+factor_hash])
	i := sort.Search(len(r.nodes), func(i int) bool {
		return r.nodes[i].hash >= hValue
	})
	if i == len(r.nodes) {
		return r.nodes[0].destination
	}
	return r.nodes[i].destination
}

func (r *Ring) build() {
	var totalWeight int
	for _, weight := range r.weights {
		totalWeight += weight
	}
	totalSlots := r.slots * len(r.weights)
	sha := sha1.New()
	for destination, weight := range r.weights {
		// 计算Key对应的虚拟节点数
		slots := int(math.Floor(float64(weight) / float64(totalWeight) * float64(totalSlots)))
		for i := 1; i <= slots; i++ {
			sha.Write([]byte(destination + ":" + strconv.Itoa(i)))
			n := node{
				hash:        getHash(sha.Sum(nil)[2:2+factor_hash]),
				destination: destination,
			}
			r.nodes = append(r.nodes, n)
			sha.Reset()
		}
	}
	sort(r.nodes)
}

func getHash(bs []byte) uint32 {
	if len(bs) < factor_hash {
		return 0
	}
	return (uint32(bs[3]) << 24) | (uint32(bs[2]) << 16) | (uint32(bs[1]) << 8) | (uint32(bs[0]))
}
