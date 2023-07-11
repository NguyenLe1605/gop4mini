package connection

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConnection(t *testing.T) {
	conn, err := NewSwitchConnection("s1", "127.0.0.1:50051", 0, "logs/s1")
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Close()
	fmt.Printf("%v\n", conn)
}

func TestMasterArbitrationUpdate(t *testing.T) {
	s1, err := NewSwitchConnection("s1", "127.0.0.1:50051", 0, "logs/s1")
	require.NoError(t, err)
	require.NotNil(t, s1)
	defer s1.Close()

	s2, err := NewSwitchConnection("s2", "127.0.0.1:50052", 1, "logs/s2")
	require.NoError(t, err)
	require.NotNil(t, s2)
	defer s2.Close()

	s1Resp, err := s1.MasterArbitrationUpdate(context.Background(), true)
	require.NoError(t, err)
	require.NotNil(t, s1Resp)


	s2Resp, err := s2.MasterArbitrationUpdate(context.Background(), true)
	require.NoError(t, err)
	require.NotNil(t, s2Resp)
}