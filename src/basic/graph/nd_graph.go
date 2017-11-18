package MyGraph

type NDGrapht struct {
	Node []string
	Dist [][]float64
}
type Slide struct {
	NodeA string
	NodeB string
	Distance float64
}

func (N *NDGrapht) Init(n []string, D [][]float64){
	N.Node = n
	N.Dist = D
}

func (N *NDGrapht) InitWithSlide(n []string,M []Slide){
	var m1 map[string]int
	for i,j:=range n{
		m1[j] = i
	}
	N.Node = n
	nn := len(n)
	var D [nn][nn]float64
	for _,m:=range M{
		i:=m1[m.NodeA]
		j:=m1[m.NodeB]
		D[i][j] = m.Distance
	}
	N.Dist = D
}