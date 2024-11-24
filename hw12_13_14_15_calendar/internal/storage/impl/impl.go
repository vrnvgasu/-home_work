package impl

import (
	"fmt"

	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

func NewIStorage() (storage.IStorage, error) {
	switch config.Cfg.DBType {
	case config.DBTypeSQL:
		return sqlstorage.New(), nil
	case config.DBTypeMemory:
		return memorystorage.New(), nil
	default:
		return nil, fmt.Errorf("unsupported db type: %s", config.Cfg.DBType)
	}
}
