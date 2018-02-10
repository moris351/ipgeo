package ipgeo

import (
	_ "bytes"
	"encoding/binary"
	_ "fmt"
	"net"
)

var full uint32 = 0xffffffff

func IPNet(ipmask string) []byte {

	_, ipnet, err := net.ParseCIDR(ipmask)
	if err != nil {
		lb.Debug("IPNet ParseCIDR err:", err, ",ipmask:", ipmask)
		return nil
	}
	return append(ipnet.IP, ipnet.Mask...)
}

func IP2IPNet(ip string, n uint32) []byte {
	b := net.ParseIP(ip)
	if b == nil {
		lb.Debug("ip:", ip, "invalid")
		return nil
	}
	b = b.To4()
	nmask := (full >> n) << n
	bmask := make([]byte, 4)
	binary.BigEndian.PutUint32(bmask, nmask)

	for i := 0; i < 4; i++ {
		b[i] = b[i] & bmask[i]
	}
	return append(b, bmask...)

}
