package interfaces

import "PriceWatcher/internal/entities/config"

type Configer interface {
	GetConfig() (config.Config, error)
}
