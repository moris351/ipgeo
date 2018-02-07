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
func IpMaskGetUint32(ipmask string) uint32 {
	_, ipnet, err := net.ParseCIDR(ipmask)
	if err != nil {
		return 0
	}
	//lb.Debug(ipnet.IP)
	ner := ipnet.IP.To4()
	u32 := binary.BigEndian.Uint32(ner)
	mask := binary.BigEndian.Uint32(ipnet.Mask)
	//lb.Debug("mask:%d\n",mask)
	///l,_:=ipnet.Mask.Size()

	return (mask & u32)
}

func IpGetUint32(ip string) uint32 {
	ipt := net.ParseIP(ip)

	ner := ipt.To4()
	u32 := binary.BigEndian.Uint32(ner)
	return u32
}
func Match(u32 uint32, u32mask uint32, mask uint32) bool {
	return u32&u32mask == u32mask
}
func Contains(ipmask string, ip string) bool {
	_, ipnet, err := net.ParseCIDR(ipmask)
	if err != nil {
		lb.Debug(ipmask, err)
		return false
	}
	bip := net.ParseIP(ip)
	return ipnet.Contains(bip)
}

func Search(ip string, ipmask []byte) []byte {
	u32 := IpGetUint32(ip)
	u32mask := IpMaskGetUint32(string(ipmask))
	lb.Debug("ipmask:\n%b\n", u32mask)
	for i := uint(1); i < 32; i++ {
		u := (u32 >> i) << i
		lb.Debug("%b\n", u)
		if u == u32mask {
			lb.Debug("match")
			return ipmask
		}
	}
	return nil
}
