// KV最大堆

package utils

type Kv struct {
	Key string `json:"key"`
	Val int    `json:"value"`
}

type KvMaxHeap []Kv

func (h KvMaxHeap) Len() int            { return len(h) }
func (h KvMaxHeap) Less(i, j int) bool  { return h[i].Val > h[j].Val }
func (h KvMaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *KvMaxHeap) Push(x interface{}) { *h = append(*h, x.(Kv)) }
func (h *KvMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type KvMinHeap []Kv

func (h KvMinHeap) Len() int            { return len(h) }
func (h KvMinHeap) Less(i, j int) bool  { return h[i].Val < h[j].Val }
func (h KvMinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *KvMinHeap) Push(x interface{}) { *h = append(*h, x.(Kv)) }
func (h *KvMinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
