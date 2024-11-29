package internalgrpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/storage"
)

func PbEventToStorageEvent(e *pb.Event) storage.Event {
	if e == nil {
		return storage.Event{}
	}

	return storage.Event{
		ID:          uint64(e.GetId()), //nolint:gosec
		Title:       e.GetTitle(),
		StartAt:     e.GetStartAt().AsTime(),
		EndAt:       e.GetEndAt().AsTime(),
		Description: e.GetDescription(),
		OwnerID:     uint64(e.GetOwnerId()), //nolint:gosec
		SendBefore:  e.SendBefore,
	}
}

func StorageEventToPbEvent(e storage.Event) *pb.Event {
	return &pb.Event{
		Id:          int64(e.ID), //nolint:gosec
		Title:       e.Title,
		StartAt:     timestamppb.New(e.StartAt),
		EndAt:       timestamppb.New(e.EndAt),
		Description: e.Description,
		OwnerId:     int64(e.OwnerID), //nolint:gosec
		SendBefore:  e.SendBefore,
	}
}

func StorageEventListToPbEventList(list []storage.Event) *pb.EventList {
	res := &pb.EventList{
		Events: make([]*pb.Event, 0, len(list)),
	}
	for _, e := range list {
		res.Events = append(res.Events, StorageEventToPbEvent(e))
	}

	return res
}
