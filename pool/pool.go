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
	procs  *sync.Map
	count  uint32
	limit  uint32
	ttl    time.Duration
	logger *log.Logger
	nca    bool // No Check Alive

	ports     *sync.Map
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
	pool.procs.Range(func(port, p interface{}) bool {
		if !pool.nca && !p.(*proc).alive {
			pool.logger.Println("PROC_DEAD", port)
			pool.remove(port.(int))
		} else if p.(*proc).start.Before(clean_time) {
			pool.logger.Println("TIMED_OUT", port)
			pool.remove(port.(int))
		}
		return true
	})
}

func NewPool(conf *config.Config, l log.Logger) *Pool {
	l.SetPrefix("[POOL] ")
	arg, earg := template.Must(template.New("arg").Parse(conf.RunCmd.Arg)), template.Must(template.New("earg").Parse(conf.ExitCmd.Arg))
	pool := &Pool{
		procs: new(sync.Map), count: 0, errMap: map[int]int{},
		limit: conf.Limit, ttl: conf.TTL, cmd: conf.RunCmd.Cmd, arg: arg, nca: conf.NoCheckAlive,
		e_cmd: conf.ExitCmd.Cmd, e_arg: earg, e_enabled: conf.ExitCmd.Enabled, logger: &l,
		portStart: int(conf.PortStart), portLimit: int(conf.PortRange),
		ports: new(sync.Map),
	}

	go func() {
		defer func() {
			pool.procs.Range(func(port, p interface{}) bool {
				pool.remove(port.(int))
				return true
			})
		}()
		for {
			pool.cleanup()
			time.Sleep(cleanupDelay)
		}
	}()

	return pool
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
