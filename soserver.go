package ipgeo

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/moris351/ipgeo/message"
	"io"
	"net"
)

const DBNAME = "geoip.store"

type sockets struct {
}

func ServeAt(port string) {
	s := &sockets{}
	s.Listen(port)
}
func (s *sockets) Listen(p string) {
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

func (s *sockets) Serve(conn net.Conn) {
	buf := make([]byte, 1512)

	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			fmt.Println("end of client")
			return
		}
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return
		}

		fmt.Printf("Received data: %x\n", buf[:n])
		//if n!=0{
		l := Locator(DBNAME)

		msg := &message.Query{}
		err = proto.Unmarshal(buf[0:n], msg)
		if err != nil {
			fmt.Println("proto Unmarshal failed,err:", err)
			return
		}

		fmt.Println("msg:", msg)
		geo, err := l.FindGeo(msg.Ip)
		status := message.StatusType_OK
		if err != nil {
			fmt.Println("Receive FindGeo failed,err:", err)
			status = message.StatusType_ERR
		}
		ans := &message.Answer{Status: status, City: geo}

		fmt.Println("ans:", ans)
		bans, err := proto.Marshal(ans)
		if err != nil {
			fmt.Println("proto Marshal failed,err:", err)
		}

		if _, err := conn.Write(bans); err != nil {
			fmt.Println("Receive answer failed,err:", err)
			return
		}
		//}

	}
}
