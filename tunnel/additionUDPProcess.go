package tunnel

import "github.com/finddiff/clashWithCache/log"

func addWorker() {
	log.Debugln("addWorker")
	for i := 0; i < 4; i++ {
		go additionUDPProcess()
	}
}

// when udp pack block udp channel addition Proc start
func additionUDPProcess() {
	for {
		select {
		case conn := <-udpQueue:
			handleUDPConn(conn)
		default:
			log.Debugln("end additionUDPProcess")
			return
		}
	}
}
