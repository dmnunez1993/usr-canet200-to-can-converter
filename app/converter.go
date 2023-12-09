package usrcanettocan

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type deviceConverterKey struct {
	host   string
	port   int64
	target string
}

type UsrCanetConverter struct {
	deviceConverters map[deviceConverterKey]*UsrCanetDeviceConverter
	mux              sync.Mutex
}

func (c *UsrCanetConverter) updateDeviceConverters() {
	config, err := GetOrCreateConverterConfig()

	if err != nil {
		log.Error(err.Error())
	}

	usedKeys := make(map[deviceConverterKey]deviceConverterKey)

	c.mux.Lock()
	defer c.mux.Unlock()

	// Add device if not already added
	for _, deviceConverterConfig := range config.DeviceConverters {
		deviceKey := deviceConverterKey{
			host:   deviceConverterConfig.Host,
			port:   deviceConverterConfig.Port,
			target: deviceConverterConfig.Target,
		}

		_, ok := c.deviceConverters[deviceKey]

		if !ok {
			newConverter := UsrCanetDeviceConverter{
				config: deviceConverterConfig,
			}

			go newConverter.run()

			c.deviceConverters[deviceKey] = &newConverter
		}

		usedKeys[deviceKey] = deviceKey
	}

	for deviceKey, deviceConverter := range c.deviceConverters {
		_, ok := usedKeys[deviceKey]

		if !ok {
			deviceConverter.stop()
			delete(c.deviceConverters, deviceKey)
		}
	}
}

func (c *UsrCanetConverter) run() {
	for {
		c.updateDeviceConverters()
		time.Sleep(time.Second * 5)
	}
}

func NewUsrCanetConverter() *UsrCanetConverter {
	c := UsrCanetConverter{}
	c.deviceConverters = make(map[deviceConverterKey]*UsrCanetDeviceConverter)

	go c.run()

	return &c
}
