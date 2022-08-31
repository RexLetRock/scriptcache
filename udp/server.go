package udp

import (
	"fmt"
	"net"
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

	data := make([]byte, 1024*10000)
	for {
		read, _, _ := socket.ReadFromUDP(data)
		fmt.Printf("%v \n", read)
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
