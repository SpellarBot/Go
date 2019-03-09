package utils

import (
	"fmt"
	"math/rand"
)

type Linearmodel struct {
	Coef []float64
	Para float64
	Dimsion int
}

func(L *Linearmodel) Init(N int){
	L.Dimsion = N
	L.Coef = make([]float64,N,N*2)
	L.Para = 0.0
}


func (L *Linearmodel) Train( Data [][]float64, Target []float64 ){
	E := 1e-4
	N := 1000
	batch_size := 10
	eta := 0.01
	n := 0
	e := 1.0
	N_Data := len(Data)
	for{
		if n>N || e < E{
			break
		}
		begin := (n*batch_size)%N
		end := begin+batch_size
		if end > N_Data {
			end = N_Data
		}
		data := Data[begin:end]
		target := Target[begin:end]
		ParaGrad(L.Coef, &(L.Para), data, target, eta)
		Y_ := L.Predict(data)
		e = CalSquareError(target, Y_)
		n++
		fmt.Println("Result:",n,e)
	}
}



func (L *Linearmodel) Predict( Data [][]float64)([]float64){
	N1 := L.Dimsion
	N2 := len(Data)
	Result := make([]float64,N2,2*N2)
	for i:=0;i<N2;i++{
		for j:=0;j<N1;j++{
			//fmt.Println(i,j,Result[i])
			Result[i] += L.Coef[j]*Data[i][j]
		}
		Result[i] += L.Para
	}
	return Result
}

func ParaGrad(A []float64, b *float64, data [][]float64, target []float64, eta float64){
	N1 := len(A)
	N2 := len(data)
	A_ := make([]float64,N1,N1)
	b_ := 0.0
	for i:=0;i<N2;i++{
		for j:=0;j<N1;j++{
			A_[j] += 2*data[i][j]*data[i][j]*A[j] + 2*(*b)*data[i][j] -2*data[i][j]*target[i]
			b_ += 2*A[j]*data[i][j]
		}
		b_ += ( 2*(*b) -2*target[i])
	}
	*b = *b - eta*b_/float64(N2)
	for i:=0;i<N1;i++{
		A[i] = A[i] - eta*A_[i]/float64(N2)
	}
	//fmt.Println(A[0],*b,A_[0],b_)
	return
}

func CalSquareError(Y []float64, Y_ []float64)(float64){
	E := 0.0
	for k,value := range Y{
		x := Y_[k]-value
		E += x*x
	}
	E = E/float64(len(Y))
	return E
}


func X_Array(A []float64, B []float64)(float64){
	m := 0.0
	for k,value:= range A{
		m += value * B[k]
	}
	return m
}

func Add_Array(A []float64, B []float64)([]float64){
	m := make([]float64,len(A),2*len(A))
	for k,value:= range A{
		m[k] = value + B[k]
	}
	return m
}

func X_Array_number(A []float64, B float64)([]float64){
	m := make([]float64,len(A),2*len(A))
	for k,value:= range A{
		m[k] = value*B
	}
	return m
}


func Test(){
	fmt.Println("Begin to Train Data")
	var Data [][]float64
	var Target []float64
	for i:=0;i<1000;i++{
		Datai := []float64{rand.Float64()}
		Data = append(Data,Datai)
		Target = append(Target,2.0*Data[i][0] + 3.0)
	}
	L := Linearmodel{}
	L.Init(1)
	L.Train(Data, Target)
}