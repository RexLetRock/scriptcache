package main

import (
	"scriptcache/tcp"
)

var (
	kMaxUInt64 uint64 = 0xFFFFFFFFFFFFFFFF
	// kEpoch              uint64 = 1288834974657
	// kWorkerIdBits       uint64 = 14
	// kWorkerMaxId        uint64 = kMaxUInt64 ^ (kMaxUInt64 << kWorkerIdBits)
	// kWorkerIdMask       uint64 = kMaxUInt64 ^ (kMaxUInt64 << kWorkerIdBits)
	// kSequenceBits       uint64 = 12
	// kSequenceMask       uint64 = kMaxUInt64 ^ (kMaxUInt64 << kSequenceBits)
	// kWorkerIdShift      uint64 = kSequenceBits
	// kTimestampLeftShift uint64 = kSequenceBits + kWorkerIdBits
)

func main() {
	// cmd.MultiChannel()

	// go udp.ServerStart()
	// go udp.ClientStart()

	// go tcpevio.ServerStart()
	// time.Sleep(1 * time.Second)
	// go tcpevio.ClientStart()

	// fmt.Printf("kMaxUInt64 %b %02x \n", kMaxUInt64, kMaxUInt64)
	// fmt.Printf("kWorkerMaxId %b %02x %v \n", kWorkerMaxId, kWorkerMaxId, kWorkerMaxId)
	// fmt.Printf("kWorkerIdMask %b %02x %v \n", kWorkerIdMask, kWorkerIdMask, kWorkerIdMask)
	// fmt.Printf("kSequenceMask %b %02x \n", kSequenceMask, kSequenceMask)
	// fmt.Printf("kTimestampLeftShift %b %02x \n", kTimestampLeftShift, kTimestampLeftShift)

	// workerId := uint64(1)
	// if workerId > kWorkerMaxId {
	// 	workerId = workerId & kWorkerIdMask
	// }
	// lastSeqId := (123 + 1) & kSequenceMask

	// d := (uint64(time.Now().UnixMilli())-kEpoch)<<kTimestampLeftShift | (workerId << kWorkerIdShift) | lastSeqId
	// fmt.Printf("%v \n", d)

	// a := uint64(6592462082522747004) & 0xFFF
	// b := uint64(6592462082589855869) & 0xFFF
	// c := uint64(6592462082656964734) & 0xFFF
	// fmt.Println(a, b, c)

	tcp.ServerStart()
	// time.Sleep(1 * time.Second)
	// go tcp.ClientStart()

	select {}
}
