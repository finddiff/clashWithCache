package rules

import (
	"fmt"
	"github.com/finddiff/clashWithCache/common/cache"
	"strconv"
	"strings"

	"github.com/finddiff/clashWithCache/component/process"
	C "github.com/finddiff/clashWithCache/constant"
	"github.com/finddiff/clashWithCache/log"
)

var processCache = cache.NewLRUCache(cache.WithAge(2), cache.WithSize(64))

type Process struct {
	adapter string
	process string
}

func (ps *Process) RuleType() C.RuleType {
	return C.Process
}

func (ps *Process) Match(metadata *C.Metadata) bool {
	key := fmt.Sprintf("%s:%s:%s", metadata.NetWork.String(), metadata.SrcIP.String(), metadata.SrcPort)
	cached, hit := processCache.Get(key)
	if !hit {
		srcPort, err := strconv.Atoi(metadata.SrcPort)
		if err != nil {
			processCache.Set(key, "")
			return false
		}

		name, err := process.FindProcessName(metadata.NetWork.String(), metadata.SrcIP, srcPort)
		if err != nil {
			log.Debugln("[Rule] find process name %s error: %s", C.Process.String(), err.Error())
		}

		processCache.Set(key, name)

		cached = name
	}

	return strings.EqualFold(cached.(string), ps.process)
}

func (p *Process) Adapter() string {
	return p.adapter
}

func (p *Process) Payload() string {
	return p.process
}

func (p *Process) ShouldResolveIP() bool {
	return false
}

func NewProcess(process string, adapter string) (*Process, error) {
	return &Process{
		adapter: adapter,
		process: process,
	}, nil
}
