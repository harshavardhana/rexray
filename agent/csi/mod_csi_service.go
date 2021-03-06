package csi

import (
	"context"
	"net"

	xctx "golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/codedellemc/gocsi"
	"github.com/codedellemc/gocsi/csi"
	"github.com/codedellemc/goioc"
)

// ServiceProvider is a gRPC endpoint that provides the CSI
// services: Controller, Identity, Node.
type ServiceProvider interface {

	// Serve accepts incoming connections on the listener lis, creating
	// a new ServerTransport and service goroutine for each. The service
	// goroutine read gRPC requests and then call the registered handlers
	// to reply to them. Serve returns when lis.Accept fails with fatal
	// errors.  lis will be closed when this method returns.
	// Serve always returns non-nil error.
	Serve(ctx context.Context, lis net.Listener) error

	// Stop stops the gRPC server. It immediately closes all open
	// connections and listeners.
	// It cancels all active RPCs on the server side and the corresponding
	// pending RPCs on the client side will get notified by connection
	// errors.
	Stop(ctx context.Context)

	// GracefulStop stops the gRPC server gracefully. It stops the server
	// from accepting new connections and RPCs and blocks until all the
	// pending RPCs are finished.
	GracefulStop(ctx context.Context)
}

type csiService struct {
	ServiceProvider

	serviceName string
	serviceType string
	sp          ServiceProvider
	conn        PipeConn
}

func newService(
	ctx context.Context,
	serviceName, serviceType string) *csiService {

	if sp, ok := goioc.New(serviceType).(ServiceProvider); ok {
		return &csiService{
			serviceName: serviceName,
			serviceType: serviceType,
			sp:          sp,
			conn:        NewPipeConn(serviceName),
		}
	}

	return nil
}

func (s *csiService) Serve(
	ctx context.Context, lis net.Listener) error {

	if lis == nil {
		lis = s.conn
	}
	return s.sp.Serve(ctx, lis)
}

func (s *csiService) Stop(ctx context.Context) {
	s.sp.Stop(ctx)
	s.conn.Close()
}

func (s *csiService) GracefulStop(ctx context.Context) {
	s.sp.GracefulStop(ctx)
	s.conn.Close()
}

func (s *csiService) dial(
	ctx xctx.Context) (client *grpc.ClientConn, err error) {

	return grpc.DialContext(
		ctx,
		s.serviceName,
		grpc.WithInsecure(),
		grpc.WithDialer(s.conn.DialGrpc),
		grpc.WithUnaryInterceptor(gocsi.ChainUnaryClient(
			gocsi.ClientCheckReponseError,
			gocsi.ClientResponseValidator)))
}

func (s *csiService) dialController(
	ctx xctx.Context) (csi.ControllerClient, error) {

	c, err := s.dial(ctx)
	if err != nil {
		return nil, err
	}
	return csi.NewControllerClient(c), nil
}

func (s *csiService) dialIdentity(
	ctx xctx.Context) (csi.IdentityClient, error) {

	c, err := s.dial(ctx)
	if err != nil {
		return nil, err
	}
	return csi.NewIdentityClient(c), nil
}

func (s *csiService) dialNode(
	ctx xctx.Context) (csi.NodeClient, error) {

	c, err := s.dial(ctx)
	if err != nil {
		return nil, err
	}
	return csi.NewNodeClient(c), nil
}

////////////////////////////////////////////////////////////////////////////////
//                            Controller Service                              //
////////////////////////////////////////////////////////////////////////////////

func (s *csiService) CreateVolume(
	ctx xctx.Context,
	req *csi.CreateVolumeRequest) (
	*csi.CreateVolumeResponse, error) {

	c, err := s.dialController(ctx)
	if err != nil {
		return nil, err
	}
	return c.CreateVolume(ctx, req)
}

func (s *csiService) DeleteVolume(
	ctx xctx.Context,
	req *csi.DeleteVolumeRequest) (
	*csi.DeleteVolumeResponse, error) {

	c, err := s.dialController(ctx)
	if err != nil {
		return nil, err
	}
	return c.DeleteVolume(ctx, req)
}

func (s *csiService) ControllerPublishVolume(
	ctx xctx.Context,
	req *csi.ControllerPublishVolumeRequest) (
	*csi.ControllerPublishVolumeResponse, error) {

	c, err := s.dialController(ctx)
	if err != nil {
		return nil, err
	}
	return c.ControllerPublishVolume(ctx, req)
}

func (s *csiService) ControllerUnpublishVolume(
	ctx xctx.Context,
	req *csi.ControllerUnpublishVolumeRequest) (
	*csi.ControllerUnpublishVolumeResponse, error) {

	c, err := s.dialController(ctx)
	if err != nil {
		return nil, err
	}
	return c.ControllerUnpublishVolume(ctx, req)
}

func (s *csiService) ValidateVolumeCapabilities(
	ctx xctx.Context,
	req *csi.ValidateVolumeCapabilitiesRequest) (
	*csi.ValidateVolumeCapabilitiesResponse, error) {

	c, err := s.dialController(ctx)
	if err != nil {
		return nil, err
	}
	return c.ValidateVolumeCapabilities(ctx, req)
}

func (s *csiService) ListVolumes(
	ctx xctx.Context,
	req *csi.ListVolumesRequest) (
	*csi.ListVolumesResponse, error) {

	c, err := s.dialController(ctx)
	if err != nil {
		return nil, err
	}
	return c.ListVolumes(ctx, req)
}

func (s *csiService) GetCapacity(
	ctx xctx.Context,
	req *csi.GetCapacityRequest) (
	*csi.GetCapacityResponse, error) {

	c, err := s.dialController(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetCapacity(ctx, req)
}

func (s *csiService) ControllerGetCapabilities(
	ctx xctx.Context,
	req *csi.ControllerGetCapabilitiesRequest) (
	*csi.ControllerGetCapabilitiesResponse, error) {

	c, err := s.dialController(ctx)
	if err != nil {
		return nil, err
	}
	return c.ControllerGetCapabilities(ctx, req)
}

////////////////////////////////////////////////////////////////////////////////
//                             Identity Service                               //
////////////////////////////////////////////////////////////////////////////////

func (s *csiService) GetSupportedVersions(
	ctx xctx.Context,
	req *csi.GetSupportedVersionsRequest) (
	*csi.GetSupportedVersionsResponse, error) {

	c, err := s.dialIdentity(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetSupportedVersions(ctx, req)
}

func (s *csiService) GetPluginInfo(
	ctx xctx.Context,
	req *csi.GetPluginInfoRequest) (
	*csi.GetPluginInfoResponse, error) {

	c, err := s.dialIdentity(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetPluginInfo(ctx, req)
}

////////////////////////////////////////////////////////////////////////////////
//                               Node Service                                 //
////////////////////////////////////////////////////////////////////////////////

func (s *csiService) NodePublishVolume(
	ctx xctx.Context,
	req *csi.NodePublishVolumeRequest) (
	*csi.NodePublishVolumeResponse, error) {

	c, err := s.dialNode(ctx)
	if err != nil {
		return nil, err
	}
	return c.NodePublishVolume(ctx, req)
}

func (s *csiService) NodeUnpublishVolume(
	ctx xctx.Context,
	req *csi.NodeUnpublishVolumeRequest) (
	*csi.NodeUnpublishVolumeResponse, error) {

	c, err := s.dialNode(ctx)
	if err != nil {
		return nil, err
	}
	return c.NodeUnpublishVolume(ctx, req)
}

func (s *csiService) GetNodeID(
	ctx xctx.Context,
	req *csi.GetNodeIDRequest) (
	*csi.GetNodeIDResponse, error) {

	c, err := s.dialNode(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetNodeID(ctx, req)
}

func (s *csiService) ProbeNode(
	ctx xctx.Context,
	req *csi.ProbeNodeRequest) (
	*csi.ProbeNodeResponse, error) {

	c, err := s.dialNode(ctx)
	if err != nil {
		return nil, err
	}
	return c.ProbeNode(ctx, req)
}

func (s *csiService) NodeGetCapabilities(
	ctx xctx.Context,
	req *csi.NodeGetCapabilitiesRequest) (
	*csi.NodeGetCapabilitiesResponse, error) {

	c, err := s.dialNode(ctx)
	if err != nil {
		return nil, err
	}
	return c.NodeGetCapabilities(ctx, req)
}
