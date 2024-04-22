package gitlab

import (
	"fmt"
	"strconv"
	"time"

	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func ptr[T any](v T) *T {
	return &v
}

func parseID(id interface{}) (string, error) {
	switch v := id.(type) {
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("invalid ID type %#v, the ID must be an int or a string", id)
	}
}

func convertTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func convertUnixSeconds(ts int64) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{
		Seconds: ts,
		Nanos:   0,
	}
}

func convertUnixMilli(ts int64) *timestamppb.Timestamp {
	const msPerSecond int64 = 1_000
	const nsPerMilli int64 = 1_000
	return &timestamppb.Timestamp{
		Seconds: ts / msPerSecond,
		Nanos:   int32((ts % msPerSecond) * nsPerMilli),
	}
}

func convertUnixNano(ts int64) *timestamppb.Timestamp {
	const nsPerSecond int64 = 1_000_000_000
	return &timestamppb.Timestamp{
		Seconds: ts / nsPerSecond,
		Nanos:   int32(ts % nsPerSecond),
	}
}

func convertDuration(d float64) *durationpb.Duration {
	return durationpb.New(time.Duration(d * float64(time.Second)))
}
