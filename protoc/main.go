package protoc

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"micro-go/protoc/pd"
)

func main() {
	data := &pd.Person_PhoneNumber{
		Number: "123456",
		Type:   pd.Person_HOME,
	}

	sendData, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("proto-Marshal-err", err, "sendData", sendData)
	}

	newData := &pd.Person_PhoneNumber{}
	err = proto.Unmarshal(sendData, newData)
	if err != nil {
		fmt.Println("proto-Unmarshal-err", err, "newData", newData)
	}

	fmt.Println("oldData", data, "newData", newData)
}
