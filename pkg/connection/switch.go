package connection

import (
	"context"
	"io"
	"io/ioutil"
	"log"

	p4info "github.com/p4lang/p4runtime/go/p4/config/v1"
	p4api "github.com/p4lang/p4runtime/go/p4/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SwitchConnection interface {
	io.Closer
	MasterArbitrationUpdate(ctx context.Context, dryRun bool) (*p4api.StreamMessageResponse_Arbitration, error)
	buildDeviceConfig(file string) ([]byte, error)
	SetForwardingPipelineConfig(ctx context.Context, dryRun bool, p4Info *p4info.P4Info, configFile string) (*p4api.SetForwardingPipelineConfigResponse, error)
}

func NewSwitchConnection(name string, addr string, deviceID uint64, protoDumpFile string) (SwitchConnection, error) {
	conn, err := connect(context.Background(), addr)
	if err != nil {
		return nil, err
	}

	p4cl := p4api.NewP4RuntimeClient(conn)

	return &switchConnection{
		grpcClient:      conn,
		p4runtimeClient: p4cl,
		name:            name,
		addr:            addr,
		deviceID:        deviceID,
		protoDumpFile:   protoDumpFile,
		p4info:          &p4info.P4Info{},
	}, nil
}

func connect(ctx context.Context, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

type switchConnection struct {
	name            string
	addr            string
	deviceID        uint64
	grpcClient      *grpc.ClientConn
	p4runtimeClient p4api.P4RuntimeClient
	p4info          *p4info.P4Info
	protoDumpFile   string
}

func (s *switchConnection) MasterArbitrationUpdate(ctx context.Context, dryRun bool) (*p4api.StreamMessageResponse_Arbitration, error) {
	electionID := &p4api.Uint128{High: 0, Low: 1}
	// set up bidirectional connection between switch and controller
	channel, err := s.p4runtimeClient.StreamChannel(ctx)
	if err != nil {
		return nil, err
	}

	request := &p4api.StreamMessageRequest{
		Update: &p4api.StreamMessageRequest_Arbitration{Arbitration: &p4api.MasterArbitrationUpdate{
			DeviceId:   s.deviceID,
			ElectionId: electionID,
		}},
	}
	if dryRun {
		log.Printf("Sending master arbitation request %v\n", request)
	}

	err = channel.Send(request)
	if err != nil {
		return nil, err
	}

	// Recieving update
	for {
		in, err := channel.Recv()
		if err != nil {
			return nil, err
		}

		switch v := in.Update.(type) {
		case *p4api.StreamMessageResponse_Arbitration:
			if dryRun {
				log.Printf("Received arbitration response %v\n", v)
			}
			if err != nil {
				return nil, err
			}
			return v, nil
		}
	}
}

func (s *switchConnection) SetForwardingPipelineConfig(ctx context.Context, dryRun bool, p4Info *p4info.P4Info, configFile string) (*p4api.SetForwardingPipelineConfigResponse, error) {
	deviceConfig, err := s.buildDeviceConfig(configFile)
	if err != nil {
		return nil, err
	}

	request := &p4api.SetForwardingPipelineConfigRequest{
		ElectionId: &p4api.Uint128{High: 0, Low: 1},
		DeviceId:   s.deviceID,
		Config:     &p4api.ForwardingPipelineConfig{},
		Action:     p4api.SetForwardingPipelineConfigRequest_VERIFY_AND_COMMIT,
	}
	request.Config.P4Info = p4Info
	request.Config.P4DeviceConfig = deviceConfig

	if dryRun {
		log.Printf("P4 runtime SetForwardingPipelineConfig request: %v\n", request)
	}
	return s.p4runtimeClient.SetForwardingPipelineConfig(ctx, request)
}

func (s *switchConnection) buildDeviceConfig(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func (s *switchConnection) Close() error {
	return s.grpcClient.Close()
}

var _ SwitchConnection = &switchConnection{}
