package buildinfo

import (
	"github.com/celestiaorg/validator-da-tracker/pkg/models/metricprom"
	"sync"
	"time"
)

const (
	//defaultCacheDuration = 20 * time.Minute
	defaultCacheDuration = 20 * time.Second
)

var (
	cacheDuration = defaultCacheDuration
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

	var metrics []metricprom.BuildInfo
	for _, metric := range c.data {
		metrics = append(metrics, metric)
	}
	return metrics
}

func (c *Cache) IsStale() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	elapsed := time.Since(c.lastUpdate)
	return elapsed > cacheDuration
}
