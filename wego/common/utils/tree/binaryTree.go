package tree

import (
	"fmt"
)

type BinaryNode struct{
	Value float64
	Left *BinaryNode
	Right *BinaryNode
}
func (B *BinaryNode) Set_value(a float64){
	(*B).Value = a
}
func (B *BinaryNode) Set_left(b1 BinaryNode){
	(*B).Left = &b1
}
func (B *BinaryNode) Set_right(b1 BinaryNode){
	(*B).Right = &b1
}


type BinaryTree struct {
	Parent BinaryNode
}
func (B *BinaryTree) InitWithArray(a []float64){
	p0 := &(*B).Parent
	N := len(a)
	(*p0).Value = a[0]
	InitWithArray(p0, a, 0, N)
	return
}

func InitWithArray(B *BinaryNode,a []float64, k int, N int){
	fmt.Println(B.Value,N,k)
	if 2*k+1 < N{
		B1 := &BinaryNode{}
		(*B1).Value = a[2*k+1]
		B.Left = B1
		InitWithArray(B1, a, 2*k+1, N)
	}
	if 2*k+2 < N {
		B2 := &BinaryNode{}
		(*B2).Value = a[2*k+2]
		B.Right = B2
		InitWithArray(B2, a, 2*k+2, N)
	}
	return
}


func (B BinaryTree) PreOrder(value *[]float64){
	a := B.Parent
	*value = PreOrder(a, *value)
}

func PreOrder(N BinaryNode, a []float64)([]float64){
	a = append(a, N.Value)
	fmt.Println(N.Value,len(a))
	if N.Left == nil && N.Right== nil {
		return a
	} else {
		//fmt.Println("---------")
		if N.Left != nil{
			a = PreOrder(*(N.Left), a)
		}
		if N.Right != nil{
			a = PreOrder(*(N.Right), a)
		}
	}
	return a
}

