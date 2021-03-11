package inbound

import (
	"net"

	"github.com/finddiff/clashWithCache/component/socks5"
	C "github.com/finddiff/clashWithCache/constant"
	"github.com/finddiff/clashWithCache/context"
)

// NewSocket recieve TCP inbound and return ConnContext
func NewSocket(target socks5.Addr, conn net.Conn, source C.Type) *context.ConnContext {
	metadata := parseSocksAddr(target)
	metadata.NetWork = C.TCP
	metadata.Type = source
	if ip, port, err := parseAddr(conn.RemoteAddr().String()); err == nil {
		metadata.SrcIP = ip
		metadata.SrcPort = port
	}

	return context.NewConnContext(conn, metadata)
}
