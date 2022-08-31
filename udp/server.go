package udp

import (
	"fmt"
	"net"
	"time"
)

func ServerStart() {
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8080,
	})
	if err != nil {
		fmt.Println("listening failed!", err)
		return
	}
	fmt.Println("Start Server")
	defer socket.Close()

	data := make([]byte, 1)
	count := 0
	time.AfterFunc(10*time.Second, func() { fmt.Printf("RECEIVE %v \n", count) })

	for {
		read, _, _ := socket.ReadFromUDP(data)
		if read != 0 {
			count += 1
		}
		// fmt.Printf("%s \n", data[:read])
		// read, remoteAddr, err := socket.ReadFromUDP(data)
		// if err != nil {
		// 	fmt.Println("read data failed!", err)
		// 	continue
		// }

		// if read == 0 && remoteAddr == nil {
		// 	continue
		// }
		// senddata := []byte("hello client!")
		// _, err = socket.WriteToUDP(senddata, remoteAddr)
		// if err != nil {
		// 	fmt.Println("send data failed!", err)
		// 	return
		// }
	}

}
