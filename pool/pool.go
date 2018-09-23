package pool

import (
	"log"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/popu125/sShare/config"
)

const (
	cleanupDelay = 5 * time.Second
)

type Pool struct {
	procs map[int]*proc
	count uint32
	limit uint32
	lock  sync.RWMutex
	ttl   time.Duration
	l     *log.Logger
	nca   bool // No Check Alive

	ports     *sync.Map //TODO: check if works
	portStart int
	portLimit int

	errMap map[int]int

	cmd       string
	arg       *template.Template
	e_cmd     string
	e_arg     *template.Template
	e_enabled bool
}

func (pool Pool) Count() uint32 {
	return atomic.LoadUint32(&pool.count)
}

func (pool *Pool) cleanup() {
	clean_time := time.Now().Add(-pool.ttl)
	for n, p := range pool.procs {
		if (!pool.nca && !p.alive) || p.start.Before(clean_time) {
			pool.remove(n, p)
		}
	}
}

func NewPool(conf config.Config, l log.Logger) *Pool {
	l.SetPrefix("[POOL] ")
	arg, earg := template.Must(template.New("arg").Parse(conf.RunCmd.Arg)), template.Must(template.New("earg").Parse(conf.ExitCmd.Arg))
	p := &Pool{
		procs: map[int]*proc{}, count: 0, errMap: map[int]int{},
		limit: conf.Limit, ttl: conf.TTL, cmd: conf.RunCmd.Cmd, arg: arg, nca: conf.NoCheckAlive,
		e_cmd: conf.ExitCmd.Cmd, e_arg: earg, e_enabled: conf.ExitCmd.Enabled, l: &l,
		portStart: int(conf.PortStart), portLimit: int(conf.PortRange),
		ports: new(sync.Map),
	}

	go func() {
		defer func() {
			for n, pr := range p.procs {
				p.remove(n, pr)
			}
		}()
		for {
			p.cleanup()
			time.Sleep(cleanupDelay)
		}
	}()

	return p
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
