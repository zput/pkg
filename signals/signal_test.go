package signals

import (
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

func TestSetupSignalHandler1(t *testing.T) {
	resetOnlyOneSignalHandler()
	defer func() {
		resetOnlyOneSignalHandler()
	}()
	var i int

	defer func() {
		if err := recover(); err != nil {
			assert.Equal(t, 1, i)
			if i == 1 {
				return
			}
			t.Errorf("excepeted i be:%d, but actual:%d", 1, i)
		}
	}()
	SetupSignalHandler()
	i++
	SetupSignalHandler()
	i++
	t.Logf("finish setup signal handler first")
}

func TestSetupSignalHandler2(t *testing.T) {
	resetOnlyOneSignalHandler()

	defer func() {
		resetOnlyOneSignalHandler()
	}()

	var pid int
	stop := make(<-chan struct{})
	ch := make(chan bool)

	go func() {
		stop = SetupSignalHandler()
		pid = os.Getpid()
		t.Logf("current pid:%d", pid)
	}()

	time.Sleep(1 * time.Second)

	go func(ch chan bool) {
		_, ok := <-stop
		ch <- ok
	}(ch)

	cmd := exec.Command("kill", "-2", strconv.Itoa(pid))
	_, err := cmd.Output()
	if err != nil {
		t.Fatalf("exec kill once error:%v", err)
	}

	ok := <-ch
	if ok {
		t.Errorf("excepted closed, but received:%v", ok)
		return
	}
	t.Logf("kill once stop chan received closed msg, process still alive and do some shutdown jobs")
	// golang中并不原生支持多进程, 对于kill第二次的情况暂不做测试
}
