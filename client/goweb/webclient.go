package main

import (
	"io"
	"net/http"
	"net"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/moris351/ipgeo/message"
)

const form = `
	<html><body>
		<h2>Input ip</h2>
		<form action="#" method="post" name="bar">
			<input type="text" name="in" />
			<input type="submit" value="submit"/>
		</form>
	</body></html>
`
func FormServer(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch request.Method {
	case "GET":
		/* display the form to the user */
		io.WriteString(w, form)
	case "POST":
		/* handle the form data, note that ParseForm must
		   be called before we can extract form data */
		//request.ParseForm();
		//io.WriteString(w, request.Form["in"][0])

		ip:=request.FormValue("in")

		out:=queryGeo(ip)
		io.WriteString(w,out )
	}
}
func connectIpGeo(port string)(net.Conn){

	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return nil
	}
	return conn
}
func queryGeo(ip string)string{

	fmt.Println("ip:",string(ip))
	query:=&message.Query{string(ip)}
	bquery,err:=proto.Marshal(query)
	if err!=nil{
		fmt.Println("Marshal failed,err:",err)
		return ""
	}
	
	n, err := connIpGeo.Write(bquery)
	if err != nil{
		fmt.Println("write failed,err:",err,",n:",n)
		return ""
	}
	buf:=make([]byte,512)
	n, err = connIpGeo.Read(buf)
	if err != nil{
		fmt.Println("write failed,err:",err)
		return ""
	}

	ans:=&message.Answer{}
	err=proto.Unmarshal(buf[0:n],ans)

	if err != nil{
		fmt.Println("Unmarshal failed,err:",err)
		return ""
	}
	fmt.Println("ans:",ans)
	return ans.City

}
var connIpGeo net.Conn
func main() {
	connIpGeo=connectIpGeo(":5000")

	http.HandleFunc("/", FormServer)
	if err := http.ListenAndServe(":8088", nil); err != nil {
		panic(err)
	}
}
