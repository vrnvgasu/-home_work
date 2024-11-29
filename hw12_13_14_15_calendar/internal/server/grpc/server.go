package internalgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/config"
	validate "github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/server/grpc/interceptors"
	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Logger interface {
	Warn(msg string)
	Info(msg string)
	Error(msg string)
	File(msg string)
}

type Server struct {
	pb.UnimplementedCalendarServer
	app    *app.App
	l      Logger
	server *grpc.Server
}

func NewServer(logger Logger, app *app.App) *Server {
	return &Server{
		app: app,
		l:   logger,
		server: grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				validate.UnaryServerRequestLoggerInterceptor(logger),
			),
		),
	}
}

func (s *Server) Add(ctx context.Context, pbE *pb.Event) (*pb.Event, error) {
	event, err := s.app.CreateEvent(ctx, PbEventToStorageEvent(pbE))
	if err != nil {
		s.l.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "server grpc Add CreateEvent failed")
	}

	return StorageEventToPbEvent(event), nil
}

func (s *Server) Update(ctx context.Context, pbE *pb.Event) (*pb.Event, error) {
	err := s.app.UpdateEvent(ctx, PbEventToStorageEvent(pbE))
	if err != nil {
		s.l.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "server grpc Update UpdateEvent failed")
	}

	return pbE, nil
}

func (s *Server) List(ctx context.Context, period *pb.Period) (*pb.EventList, error) {
	d := period.GetDuration().Number()

	params := storage.Params{
		StartAtGEq: period.GetStart().AsTime(),
	}
	params.SetStartAtLEqFromDuration(storage.Duration(d))

	events, err := s.app.EventList(ctx, params)
	if err != nil {
		s.l.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "server grpc List EventList failed")
	}

	return StorageEventListToPbEventList(events), nil
}

func (s *Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp",
		fmt.Sprintf("%s:%d", config.Cfg.GRPSServer.Host, config.Cfg.GRPSServer.Port))
	if err != nil {
		return fmt.Errorf("start grps server Listen: %w", err)
	}

	pb.RegisterCalendarServer(s.server, s)

	s.l.Info(fmt.Sprintf("starting server on %s", lsn.Addr().String()))
	if err = s.server.Serve(lsn); err != nil {
		return fmt.Errorf("start grpc server Serve: %w", err)
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop() error {
	s.server.GracefulStop()

	return nil
}
