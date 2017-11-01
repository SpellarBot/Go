package MyList

type MyStack struct{
	data []int
	Size int
}

func (d *MyStack) Init(){
	d.data = make([]int,0,100)
	d.Size = 0
}
func (d*MyStack) Init_Array(a []int){
	d.Init()
	copy(d.data, a)
	d.Size = len(a)

}

func (d *MyStack) Push(a int){
	d.Size++
	d.data = append(d.data,a)
}
func (d *MyStack) Pop(){
	N := d.Size
	d.data = d.data[0:N-1]
	d.Size--
}
