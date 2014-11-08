package cloudstack

import (
	"fmt"
	"net/url"
	"os"

	"github.com/atsaki/golang-cloudstack-library"
)

// Config is the configuration structure used to instantiate the CloudStack
// provider.
type Config struct {
	EndPoint  string
	ApiKey    string
	SecretKey string

	client *cloudstack.Client
}

func (c *Config) loadAndValidate() error {

	if c.EndPoint == "" {
		c.EndPoint = os.Getenv("CLOUDSTACK_ENDPOINT")
	}
	if c.ApiKey == "" {
		c.ApiKey = os.Getenv("CLOUDSTACK_APIKEY")
	}
	if c.SecretKey == "" {
		c.SecretKey = os.Getenv("CLOUDSTACK_SECRETKEY")
	}

	endpoint, err := url.Parse(c.EndPoint)
	if err != nil {
		fmt.Errorf("Error parse endpoint (%s): %s",
			c.EndPoint, err)
	}

	c.client, err = cloudstack.NewClient(*endpoint, c.ApiKey, c.SecretKey, "", "")
	if err != nil {
		fmt.Errorf("Error failed to create new client. %s", err)
	}

	return nil
}
