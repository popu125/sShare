package pool

import (
	"bytes"
	"errors"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type proc struct {
	cmd   *exec.Cmd
	start time.Time
	alive bool

	port string
	info string
}

func (p *proc) Watch() {
	p.alive = true
	p.cmd.Run()
	p.alive = false
}

func (pool *Pool) NewProc() (int, string, error) {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	if pool.count >= pool.limit {
		return 0, "", errors.New("Error: Full")
	}

	var port int
	for {
		port = pool.portStart + rand.Intn(pool.portLimit)
		if na, ok := pool.ports.Load(port); !ok || !na.(bool) {
			break
		}
	}
	ports := strconv.Itoa(port)

	pool.count += 1
	g := newArgData(ports)
	buf := new(bytes.Buffer) //TODO: There should be checked carefully
	check(pool.arg.Execute(buf, g))
	args := strings.Split(buf.String(), " ")
	info := g.Data()

	np := exec.Command(pool.cmd, args...)
	p := &proc{
		cmd:   np,
		start: time.Now(),
		port:  ports, info: info,
	}

	pool.procs[port] = p
	go p.Watch()

	pool.l.Println("PROC_SPAWN", ports, info)
	return port, info, nil
}

func (pool *Pool) remove(n int, p *proc) {
	pool.lock.Lock()
	if p.alive {
		p.cmd.Process.Kill()
	}
	delete(pool.procs, n)
	pool.count -= 1
	pool.lock.Unlock()
	pool.l.Println("CLEANUP", p.port, p.info)

	buf := new(bytes.Buffer)
	check(pool.e_arg.Execute(buf, newExitArgData(p.port)))
	args := strings.Split(buf.String(), " ")
	exec.Command(pool.e_cmd, args...).Start()
}
