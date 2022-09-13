package zmws

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/smallnest/epoller"
)

const HostAddress = "0.0.0.0:8000"

var poller epoller.Poller

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}
	if err := poller.Add(conn); err != nil {
		log.Printf("Failed to add connection %v", err)
		conn.Close()
	}
}

func MainZmws() {
	// Increase resources limitations
	// var rLimit syscall.Rlimit
	// if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
	// 	panic(err)
	// }
	// rLimit.Cur = rLimit.Max
	// if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
	// 	panic(err)
	// }

	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	// Start epoll
	var err error
	poller, err = epoller.NewPollerWithBuffer(128)
	if err != nil {
		panic(err)
	}

	go Start()

	http.HandleFunc("/", wsHandler)
	log.Printf("Server listen at %v", HostAddress)
	if err := http.ListenAndServe(HostAddress, nil); err != nil {
		log.Fatal(err)
	}
}

func Start() {
	count := 0
	for {
		connections, err := poller.WaitWithBuffer()
		if err != nil {
			if err.Error() != "bad file descriptor" {
				log.Printf("failed to poll: %v", err)
			}
			continue
		}

		for _, conn := range connections {
			if conn == nil {
				break
			}
			if _, _, err := wsutil.ReadClientData(conn); err != nil {
				if err := poller.Remove(conn); err != nil {
					log.Printf("Failed to remove %v", err)
				}
				conn.Close()
			} else {
				count++
				if count%100000 == 0 || count == 1 {
					log.Printf("Msg received: %v", count)
				}
			}
		}
	}
}
