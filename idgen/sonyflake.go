package idgen

import (
	"fmt"
	"strconv"
	"sync"

	kitlog "libs/kit-logger/pkg/logger"

	"github.com/sony/sonyflake"
)

// ID is a custom type representing a unique identifier.
type ID uint64

var (
	sonyflakeInstance     *sonyflake.Sonyflake
	sonyflakeInstanceOnce sync.Once
)

func initSonyflake() error {
	var initErr error
	sonyflakeInstanceOnce.Do(func() {
		sonyflakeInstance = sonyflake.NewSonyflake(sonyflake.Settings{})
		if sonyflakeInstance == nil {
			initErr = fmt.Errorf("failed to create Sonyflake instance")
			kitlog.L().Error("Failed to create Sonyflake", "error", initErr)
		}
	})
	return initErr
}

// GenerateID generates a new unique ID.
func GenerateID() ID {
	if err := initSonyflake(); err != nil {
		kitlog.L().Error("Failed to initialize Sonyflake", "error", err)
		return 0
	}
	if sonyflakeInstance == nil {
		kitlog.L().Error("Sonyflake instance is nil")
		return 0
	}
	id, err := sonyflakeInstance.NextID()
	if err != nil {
		kitlog.L().Error("Failed to generate ID", "error", err)
		return 0
	}
	return ID(id)
}

// ToString converts an ID to its string representation.
func ToString(id ID) string {
	return fmt.Sprintf("%d", id)
}

// StringToSonyFlake converts a string representation of a Sonyflake ID to a uint64 ID.
func StringToSonyFlake(idStr string) (uint64, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse string to uint64: %v", err)
	}
	return id, nil
}
