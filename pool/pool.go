package pool

import (
	"os/exec"
	"time"
	"sync/atomic"
	"sync"
	"log"
	"strings"
	"syscall"

	"github.com/popu125/sShare/config"
)

const (
	cleanupDelay = 5 * time.Second

	ERR_FULL  = iota
	ERR_SPAWN
	DONE
)

type proc struct {
	cmd   *exec.Cmd
	start time.Time
	alive bool

	port string
	pass string
}

func (self *proc) Watch() {
	self.alive = true
	self.cmd.Run()
	self.alive = false
}

type Pool struct {
	procs map[uint]*proc
	count uint32
	limit uint32
	lock  sync.RWMutex
	ttl   time.Duration
	l     *log.Logger
	nca   bool // No Check Alive

	lastid uint
	errMap map[uint]uint

	cmd       string
	args      []string
	e_cmd     string
	e_args    []string
	e_enabled bool
}

func (self Pool) Count() uint32 {
	return atomic.LoadUint32(&self.count)
}

func (self *Pool) NewProc(port string, pass string) uint {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.count >= self.limit {
		return ERR_FULL
	}

	self.count += 1
	args := self.args
	for i, a := range args {
		switch a {
		case "{{pass}}":
			args[i] = pass
		case "{{port}}":
			args[i] = port
		}
	}
	np := exec.Command(self.cmd, args...)
	p := &proc{
		cmd:   np,
		start: time.Now(),
		port:  port, pass: pass,
	}

	self.lastid += 1
	procid := self.lastid
	self.procs[procid] = p
	go self.run(procid)

	self.l.Println("PROC_SPAWN", port, pass)
	return procid
}

func (self *Pool) cleanup() {
	clean_time := time.Now().Add(-self.ttl)
	for n, p := range self.procs {
		if (!self.nca && !p.alive) || p.start.Before(clean_time) {
			self.remove(n, p)
		}
	}
}

func (self *Pool) remove(n uint, p *proc) {
	self.lock.Lock()
	if p.alive {
		//p.cmd.Process.Kill()
		p.cmd.Process.Signal(syscall.SIGTERM)
	}
	delete(self.procs, n)
	self.count -= 1
	self.lock.Unlock()
	self.l.Println("CLEANUP", p.port, p.pass)

	args := self.e_args
	for i, a := range args {
		switch a {
		case "{{port}}":
			args[i] = p.port
		}
	}
	exec.Command(self.e_cmd, args...).Start()
}

func (self *Pool) Check(procid uint) bool {
	p := self.procs[procid]
	if p == nil {
		return false
	} else {
		return p.alive
	}
}

func (self *Pool) run(procid uint) {
	p := self.procs[procid]
	go p.Watch()
	if !p.alive {
		self.errMap[procid] = ERR_SPAWN
	} else {
		self.errMap[procid] = DONE
	}
}

func NewPool(conf config.Config, l log.Logger) *Pool {
	l.SetPrefix("[POOL] ")
	args, eargs := strings.Split(conf.RunCmd.Arg, " "), strings.Split(conf.ExitCmd.Arg, " ")
	p := &Pool{
		procs: map[uint]*proc{}, count: 0, lastid: 10, errMap: map[uint]uint{},
		limit: conf.Limit, ttl: conf.TTL, cmd: conf.RunCmd.Cmd, args: args, nca: conf.NoCheckAlive,
		e_cmd: conf.ExitCmd.Cmd, e_args: eargs, e_enabled: conf.ExitCmd.Enabled, l: &l,
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
