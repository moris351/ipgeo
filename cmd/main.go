package main 

import(
	"fmt"
	"github.com/moris351/ipgeo"
)

func main(){
	l:=ipgeo.Locator()
	if l==nil{
		return
	}
	
	fmt.Println(l)
	if err:=l.InitDB();err!=nil{
		return
	}

	l.Show()
	l.Close()
}

