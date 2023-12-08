package buildinfo

import (
	"sync"
	"time"

	"github.com/celestiaorg/validator-da-tracker/pkg/models/metricprom"
)

const (
	defaultCacheDuration = 20 * time.Second
)

type Cache struct {
	mutex      sync.RWMutex
	data       map[string]metricprom.BuildInfo
	lastUpdate time.Time
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]metricprom.BuildInfo),
	}
}

func (c *Cache) Update(metrics []metricprom.BuildInfo) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, metric := range metrics {
		c.data[metric.Metric.Instance] = metric
	}
	c.lastUpdate = time.Now()
}

func (c *Cache) GetMetrics() []metricprom.BuildInfo {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	metrics := make([]metricprom.BuildInfo, 0, len(c.data))
	for _, metric := range c.data {
		metrics = append(metrics, metric)
	}
	return metrics
}

func (c *Cache) IsStale() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	elapsed := time.Since(c.lastUpdate)
	return elapsed > defaultCacheDuration
}
