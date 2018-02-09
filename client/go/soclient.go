package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
	"flag"
	"github.com/golang/protobuf/proto"
	"github.com/moris351/ipgeo/message"
)
var(
	port = flag.String("port", ":5000", "server port")
)
func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", *port)
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return 
	}

	ff,err:=os.Open("ips")
	defer ff.Close()
	if err!=nil{
		return
	}

	reader := bufio.NewReader(ff)
	buf:=make([]byte,1512)
	i:=0
	for {
		//if i>10{break}
		i++
		str,_,err:=reader.ReadLine()
		if err!=nil{
			fmt.Println("read failed,err:",err)
			return 
		}

		fmt.Println("str:",string(str))
		query:=&message.Query{string(str)}
		bquery,err:=proto.Marshal(query)
		if err!=nil{
			fmt.Println("Marshal failed,err:",err)
			return
		}
		
		n, err := conn.Write(bquery)
		if err != nil{
			fmt.Println("write failed,err:",err,",n:",n)
		}
		n, err = conn.Read(buf)
		if err != nil{
			fmt.Println("write failed,err:",err)
			return
		}

		ans:=&message.Answer{}
		err=proto.Unmarshal(buf[0:n],ans)

		if err != nil{
			fmt.Println("Unmarshal failed,err:",err)
			return
		}
		fmt.Println("ans:",ans)
		time.Sleep(1 * time.Millisecond)

	}
}
