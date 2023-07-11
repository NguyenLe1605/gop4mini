package controller

import (
	"fmt"
	"io/ioutil"

	p4info "github.com/p4lang/p4runtime/go/p4/config/v1"
	"google.golang.org/protobuf/encoding/prototext"
)

type P4InfoHelper struct {
	P4Info *p4info.P4Info
}

func NewP4InfoHelper(p4InfoFilePath string) (*P4InfoHelper, error) {
	data, err := ioutil.ReadFile(p4InfoFilePath)
	if err != nil {
		return nil, err
	}

	info := &p4info.P4Info{}
	err = prototext.Unmarshal(data, info)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v", info)

	return &P4InfoHelper{
		P4Info: info,
	}, nil
}
