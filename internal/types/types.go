package types

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func ConvertUnixSeconds(ts int64) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{
		Seconds: ts,
		Nanos:   0,
	}
}

func ConvertUnixMilli(ts int64) *timestamppb.Timestamp {
	const msPerSecond int64 = 1_000
	const nsPerMilli int64 = 1_000
	return &timestamppb.Timestamp{
		Seconds: ts / msPerSecond,
		Nanos:   int32((ts % msPerSecond) * nsPerMilli),
	}
}

func ConvertUnixNano(ts int64) *timestamppb.Timestamp {
	const nsPerSecond int64 = 1_000_000_000
	return &timestamppb.Timestamp{
		Seconds: ts / nsPerSecond,
		Nanos:   int32(ts % nsPerSecond),
	}
}

func ConvertDuration(d float64) *durationpb.Duration {
	return durationpb.New(time.Duration(d * float64(time.Second)))
}
