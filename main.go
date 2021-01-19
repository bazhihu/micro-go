package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	main2 "micro-go/protoc"
)

func main() {
	data := &main2.Person_PhoneNumber{
		Number: "123456",
		Type:   main2.Person_HOME,
	}

	sendData, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("proto-Marshal-err", err, "sendData", sendData)
	}

	newData := &main2.Person_PhoneNumber{}
	err = proto.Unmarshal(sendData, newData)
	if err != nil {
		fmt.Println("proto-Unmarshal-err", err, "newData", newData)
	}

	fmt.Println("oldData", data, "newData", newData)
}
