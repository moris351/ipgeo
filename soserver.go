package ipgeo

import(
	"fmt"
	"github.com/moris351/ipgeo/message"
	"net"
	"github.com/golang/protobuf/proto"
)

const DBNAME = "geoip.store"
type sockets struct{
}
func ServeAt(port string){
	s:=&sockets{}
	s.Listen(port)
}
func(s *sockets)Listen(p string) {
	fmt.Println("Starting the server ...")
	listener, err := net.Listen("tcp", p)
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return 
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return 
		}
		go s.Serve(conn)
	}
}

func (s *sockets)Serve(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return 
		}

		fmt.Printf("Received data: %v", string(buf[:len]))
		if len!=0{
			l:=Locator(DBNAME)

			var msg *message.Query
			err:=proto.Unmarshal(buf,msg)
			if err!=nil{
				fmt.Println("proto Unmarshal failed,err:",err)
				return
			}

			geo,err:=l.FindGeo(msg.Ip)
			if err!=nil{
				fmt.Println("Receive FindGeo failed,err:",err)
				return
			}
			
			ans:=&message.Answer{geo}
			bans,err:=proto.Marshal(ans)
			if err!=nil{
				fmt.Println("proto Marshal failed,err:",err)
			}

			if _,err:=conn.Write(bans);err!=nil{
				fmt.Println("Receive answer failed,err:",err)
				return
			}
		}

	}
}
	

