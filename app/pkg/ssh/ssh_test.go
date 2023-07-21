package ssh

import (
	"context"
	"fmt"
	"testing"
	"time"
)

//func TestSshClient_Key(t *testing.T) {
//	_ = Init(&Config{Timeout: time.Second * 10})
//	sf := ServerConfig{
//		User:     "ipfs",
//		Host:     "192.168.1.213",
//		Password: "ipfs",
//		Port:     22,
//	}
//	r, err := NewTerminal(sf, 400, 500)
//	fmt.Println(err)
//	time.Sleep(time.Second * 7)
//	r.Close()
//	out, err := RunCmd(sf, "pwd")
//	time.Sleep(time.Second * 5)
//	fmt.Println(string(out), err)
//	time.Sleep(time.Second * 20)
//	fmt.Println(r)
//
//}
//
//type singleWriter struct {
//	b  bytes.Buffer
//	mu sync.Mutex
//}
//
//func (w *singleWriter) Write(p []byte) (int, error) {
//	w.mu.Lock()
//	defer w.mu.Unlock()
//	fmt.Println("singleWriter: ", string(p))
//	return w.b.Write(p)
//}
//func TestNewSession(t *testing.T) {
//
//	a, b := user.Current()
//	fmt.Println(a, b)
//	_ = Init(&Config{Timeout: time.Second * 10})
//	sf := ServerConfig{
//		User:     "ipfs",
//		Host:     "192.168.1.213",
//		Password: "ipfs",
//		Port:     22,
//	}
//	sess, err := NewSession(sf)
//	fmt.Println("NewSession: ", err)
//	//w := &singleWriter{}
//	//sess.session.Stdout = w
//	go func() {
//		t := time.NewTimer(time.Second * 6)
//		for {
//			select {
//			case <-t.C:
//				err := sess.session.Signal(ssh.SIGINT)
//				err = sess.session.Signal(ssh.SIGKILL)
//				fmt.Println("stop Signal:", err)
//				return
//			}
//		}
//	}()
//	//err = sess.session.Run("ping baidu.com")
//	bb, err := sess.session.CombinedOutput("ping baidu.com")
//	fmt.Println("bb: ", string(bb))
//	//fmt.Println(string(w.b.Bytes()))
//	if e, ok := err.(*ssh.ExitError); ok {
//		fmt.Println("Signal:", e.Signal())
//		fmt.Println("ExitStatus:", e.ExitStatus())
//		fmt.Println("Msg:", e.Msg())
//		fmt.Println("Error:", e.Error())
//	}
//	fmt.Printf("err : %+v\r\n", err)
//}

type s struct {
	typ string
}

func (st *s) Write(p []byte) (n int, err error) {
	fmt.Println(st.typ, string(p))
	return len(p), nil
}

func TestNewSession(t *testing.T) {
	_ = Init(&Config{Timeout: time.Second * 10})
	//cmd := "echo ${DD}"
	sess, err := NewRemoteExec(ServerConfig{
		Host:     "192.168.1.180",
		User:     "ipfs",
		Password: "ipfs",
		Port:     22,
	})
	if err == nil {
		fmt.Println(err)
		if err == nil {
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				t := time.NewTimer(time.Second * 5)
				select {
				case <-t.C:
					cancel()
					return
				}
			}()
			sess.WithCtx(ctx)
			sess.WithEnvs(NewEnvsBySliceKV([]string{"BD=baidu.com"}))
			b, err := sess.Run("ping $BD")
			fmt.Println(string(b), err)
		}
	}
	fmt.Println(111, err)
}
