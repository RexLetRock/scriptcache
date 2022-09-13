package zgobws

import "time"

func MainGobws() {
	go MainServerGobws()
	time.Sleep(1 * time.Second)
	MainClientGobws()
}
