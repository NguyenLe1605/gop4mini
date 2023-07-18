package controller

import (
	"context"
	"fmt"

	"github.com/NguyenLe1605/gop4mini/pkg/connection"
)

const (
	SWITCH_TO_HOST_PORT   uint16 = 1
	SWITCH_TO_SWITCH_PORT uint16 = 2
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
	_, err = s2.MasterArbitrationUpdate(context.Background(), true)
	if err != nil {
		return err
	}

	_, err = s1.SetForwardingPipelineConfig(context.Background(), false, helper.P4Info, bmv2JsonFilePath)
	if err != nil {
		return err
	}
	fmt.Println("Install p4 program using SetForwardingPipelineConfig on s1")

	_, err = s2.SetForwardingPipelineConfig(context.Background(), false, helper.P4Info, bmv2JsonFilePath)
	if err != nil {
		return err
	}
	fmt.Println("Install p4 program using SetForwardingPipelineConfig on s2")

	if err = writeTunnelRules(helper, s1, s2, 100, "08:00:00:00:02:22", "10.0.2.2"); err != nil {
		return err
	}

	if err = writeTunnelRules(helper, s2, s1, 200, "08:00:00:00:01:11", "10.0.1.1"); err != nil {
		return err
	}

	return nil
}

func writeTunnelRules(
	helper *P4InfoHelper,
	ingressSw connection.SwitchConnection,
	egressSw connection.SwitchConnection,
	tunnelId uint16,
	dstEthAddr string,
	dstIpAddr string,
) error {
	entry, err := helper.BuildTableEntry(
		"MyIngress.ipv4_lpm",
		map[string]interface{}{
			"hdr.ipv4.dstAddr": struct {
				string
				int32
			}{
				dstIpAddr, 32,
			},
		},
		"MyIngress.myTunnel_ingress",
		map[string]interface{}{
			"dst_id": tunnelId,
		},
	)

	if err != nil {
		return err
	}

	if _, err = ingressSw.WriteTableEntry(context.Background(), entry, false); err != nil {
		return err
	}
	fmt.Printf("Installed ingress tunnel rule on %s\n", ingressSw.GetName())

	entry, err = helper.BuildTableEntry(
		"MyIngress.myTunnel_exact",
		map[string]interface{}{
			"hdr.myTunnel.dst_id": tunnelId,
		},
		"MyIngress.myTunnel_forward",
		map[string]interface{}{
			"port": SWITCH_TO_SWITCH_PORT,
		},
	)
	if err != nil {
		return err
	}

	if _, err = ingressSw.WriteTableEntry(context.Background(), entry, false); err != nil {
		return err
	}
	fmt.Println("Install transit tunnel rule")

	entry, err = helper.BuildTableEntry(
		"MyIngress.myTunnel_exact",
		map[string]interface{}{
			"hdr.myTunnel.dst_id": tunnelId,
		},
		"MyIngress.myTunnel_egress",
		map[string]interface{}{
			"dstAddr": dstEthAddr,
			"port":    SWITCH_TO_HOST_PORT,
		},
	)
	if err != nil {
		return err
	}

	if _, err = egressSw.WriteTableEntry(context.Background(), entry, false); err != nil {
		return err
	}
	fmt.Printf("Install egress tunnel rule on %s\n", egressSw.GetName())

	return nil
}
