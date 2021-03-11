package outbound

import (
	"context"
	"github.com/finddiff/clashWithCache/component/resolver"
	"net"
	"strconv"

	"github.com/finddiff/clashWithCache/component/dialer"
	C "github.com/finddiff/clashWithCache/constant"
	"github.com/finddiff/clashWithCache/log"
)

type ReDirect struct {
	*Base
	targetHost    string
	targetPort    int
	targetPortStr string
	//targetTCPAddr net.Addr
	targetUDPAddr net.UDPAddr
	//pc C.PacketConn
}

type ReDirectOption struct {
	Name   string `proxy:"name"`
	Server string `proxy:"server"`
	Port   int    `proxy:"port"`
}

func (d *ReDirect) DialContext(ctx context.Context, metadata *C.Metadata) (C.Conn, error) {
	address := net.JoinHostPort(d.targetHost, d.targetPortStr)

	c, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, err
	}
	tcpKeepAlive(c)
	return NewConn(c, d), nil
}

func (d *ReDirect) DialUDP(metadata *C.Metadata) (C.PacketConn, error) {
	pc, err := dialer.ListenPacket("udp", "")
	if err != nil {
		return nil, err
	}
	//d.pc = newPacketConn(&redirectPacketConn{pc}, d)
	port, _ := strconv.Atoi(metadata.DstPort)
	return newPacketConn(
		&redirectPacketConn{
			pc,
			&d.targetUDPAddr,
			&net.UDPAddr{
				IP:   metadata.DstIP,
				Port: port,
			}}, d), nil
}

func (r *redirectPacketConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	log.Debugln("WriteTo:%v", r.targetUDPAddr)
	return r.PacketConn.WriteTo(b, r.targetUDPAddr)
}

func (r *redirectPacketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	log.Debugln("ReadFrom:%v, ", r.sourceUDPAddr)
	n, _, e := r.PacketConn.ReadFrom(b)
	if e != nil {
		return 0, nil, e
	}
	return n, r.sourceUDPAddr, nil
}

type redirectPacketConn struct {
	net.PacketConn
	targetUDPAddr *net.UDPAddr
	sourceUDPAddr *net.UDPAddr
}

func NewReDirect(op ReDirectOption) (*ReDirect, error) {
	ip, err := resolver.ResolveIP(op.Server)
	if err != nil {
		return nil, err
	}

	portstr := strconv.Itoa(op.Port)

	return &ReDirect{
		Base: &Base{
			name: op.Name,
			tp:   C.ReDirect,
			udp:  true,
		},
		targetHost: op.Server,
		targetPort: op.Port,
		targetUDPAddr: net.UDPAddr{
			IP:   ip,
			Port: op.Port,
		},
		targetPortStr: portstr,
		//targetTCPAddr:ip,
	}, nil
}
