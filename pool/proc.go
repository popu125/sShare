package pool

import (
	"bytes"
	"errors"
	"math/rand"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"
)

type proc struct {
	cmd   *exec.Cmd
	start time.Time
	alive bool

	port    int
	info    string
	argData *runArgData
}

func (p *proc) Watch() {
	p.alive = true
	p.cmd.Run()
	p.alive = false
}

func (pool *Pool) NewProc() (int, *runArgData, error) {
	if pool.count >= pool.limit {
		return 0, nil, errors.New("Error: Full")
	}

	var port int
	for {
		port = pool.portStart + rand.Intn(pool.portLimit)
		if na, ok := pool.ports.Load(port); !ok || !na.(bool) {
			pool.ports.Store(port, true)
			break
		}
	}

	atomic.AddUint32(&pool.count, 1)
	argData := newArgData(port)
	buf := new(bytes.Buffer)
	check(pool.arg.Execute(buf, argData))
	args := strings.Split(buf.String(), " ")
	info := argData.Data()

	np := exec.Command(pool.cmd, args...)
	p := &proc{
		cmd:   np,
		start: time.Now(),
		port:  port, info: info,
		argData: argData,
	}

	pool.procs.Store(port, p)
	go p.Watch()

	pool.logger.Println("PROC_SPAWN", port, info)
	return port, argData, nil
}

func (pool *Pool) remove(port int) {
	var p *proc
	if tmpProc, ok := pool.procs.Load(port); ok {
		p = tmpProc.(*proc)
	} else {
		return
	}

	if p.alive {
		p.cmd.Process.Kill()
	}

	pool.ports.Delete(port)
	pool.procs.Delete(port)
	atomic.AddUint32(&pool.count, ^uint32(0))
	pool.logger.Println("REMOVE", p.port, p.info)

	buf := new(bytes.Buffer)
	check(pool.e_arg.Execute(buf, p.argData))
	args := strings.Split(buf.String(), " ")
	exec.Command(pool.e_cmd, args...).Start()
}
