package main

//#define LZ4_DEBUG 1
//#define SUPRESS_WARNING 1
//#include "lib/lz4.h"
import "C"
import (
	"errors"
	"fmt"
	"log"
	ipc "publisher/ipc"
	"time"
	"unsafe"

	"google.golang.org/protobuf/proto"
)

func CompressData(data []byte) ([]byte, error) {
	CompressedData := make([]byte, len(data))
	if len(data) <= 0 {
		return nil, errors.New("data is null")
	}
	// log.Println("Data Cap:", cap(data), "CompressedData Cap: ", cap(CompressedData))
	compressedSize := C.LZ4_compress_default((*C.char)(unsafe.Pointer(&data[0])), (*C.char)(unsafe.Pointer(&CompressedData[0])), C.int(len(data)), C.int(len(CompressedData)))
	log.Printf("Data Len: %dkb | CompressedData Len: %dkb", len(data)/1024, compressedSize/1024)
	if int(compressedSize) == 0 {
		return nil, errors.New("failed to compress")
	}
	// log.Println("Compressed Data Length: ", compressedSize, len(data))
	return CompressedData, nil
}

func createBatch(batchSize int) ([]byte, error) {
	mqttBatch := ipc.MqttBatch{}
	mqttMsg := ipc.MqttMsg{}
	pkt_cnt := 0
	for pkt_cnt < batchSize {
		mqttMsg.Timestamp = uint64(time.Now().Unix())
		mqttMsg.Key = "BCM_ODO"
		mqttMsg.Value = fmt.Sprint(pkt_cnt)
		pkt_cnt++
		mqttBatch.Msgbatch = append(mqttBatch.Msgbatch, &mqttMsg)
	}
	tempbuf, err := proto.Marshal(&mqttBatch)
	if err != nil {
		log.Fatalln("Cannot Marshal")
		goto ERR_RET
	}
	return tempbuf, nil

ERR_RET:
	return nil, err
}
