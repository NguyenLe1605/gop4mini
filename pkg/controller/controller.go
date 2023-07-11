package controller

import (
	"context"
	"log"

	"github.com/NguyenLe1605/gop4mini/pkg/connection"
)

func Run(p4InfoFilePath string, bmv2JsonFilePath string) error {
	helper, err := NewP4InfoHelper(p4InfoFilePath)
	if err != nil {
		return err
	}

	s1, err := connection.NewSwitchConnection(
		"s1",
		"127.0.0.1:50051",
		0,
		"logs/s1",
	)
	if err != nil {
		return err
	}
	defer s1.Close()

	s2, err := connection.NewSwitchConnection(
		"s2",
		"127.0.0.1:50052",
		1,
		"logs/s2",
	)
	if err != nil {
		return err
	}
	defer s2.Close()

	_, err = s1.MasterArbitrationUpdate(context.Background(), true)
	if err != nil {
		return err
	}

	s2.MasterArbitrationUpdate(context.Background(), true)
	if err != nil {
		return err
	}

	_, err = s1.SetForwardingPipelineConfig(context.Background(), false, helper.P4Info, bmv2JsonFilePath)
	if err != nil {
		return err
	}
	log.Println("Install p4 program using SetForwardingPipelineConfig on s1")

	_, err = s2.SetForwardingPipelineConfig(context.Background(), false, helper.P4Info, bmv2JsonFilePath)
	if err != nil {
		return err
	}
	log.Println("Install p4 program using SetForwardingPipelineConfig on s2")

	return nil
}
