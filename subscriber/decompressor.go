package main

// #define LZ4_DEBUG 1
//#define SUPRESS_WARNING 1
// #include "lib/lz4.h"
import "C"

import (
	"log"
	"unsafe"

	"subscriber/ipc"

	"google.golang.org/protobuf/proto"
)

func DecompressSubscribedData() {
	for {
		msg := <-SubscribedDatChan
		data := msg.Payload()
		decompressedData := make([]byte, 100*len(data))

		n := C.LZ4_decompress_safe((*C.char)(unsafe.Pointer(&data[0])), (*C.char)(unsafe.Pointer(&decompressedData[0])), C.int(len(data)), C.int(cap(decompressedData)))
		if int(n) > cap(decompressedData) {
			log.Println("uncompressed length greater than buffer: ")
		} else {
			m := ipc.MqttBatch{}
			proto.Unmarshal(decompressedData, &m)

			cnt := 0
			for range m.Msgbatch {
				cnt++
			}
		}
	}
}
