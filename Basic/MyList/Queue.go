package MyList

type MyQueue struct {
	data []int
	Size int
}
func (d *MyQueue) Init(){
	d.data = make([]int,0,100)
	d.Size = 0
}
func (d*MyQueue) Init_Array(a []int){
	d.Init()
	copy(d.data, a)
	d.Size = len(a)

}
func (d *MyQueue) Push(a int){
	d.Size++
	d.data = append(d.data,0)
	for i:=d.Size-1;i>0;i--{
		d.data[i] = d.data[i-1]
	}
	d.data[0] = a
}
func (d *MyQueue) Pop(){
	N := d.Size
	d.data = d.data[1:N]
	d.Size--
}