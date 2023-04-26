package emongo

import (
	"github.com/gotomicro/ego/core/elog"
)

const PackageName = "component.emongo"

// Component client (cmdable and config)
type Component struct {
	config *config
	dbName string // dbname 解析后的dbName放到这里
	client *Client
	logger *elog.Component
}

// Client returns emongo Client
func (c *Component) Client() *Client {
	return c.client
}

// DbName returns emongo Client
func (c *Component) DbName() string {
	return c.dbName
}
