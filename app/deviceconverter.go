package usrcanettocan

import (
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type UsrCanetDeviceConverter struct {
	config         UsrCanetDeviceConverterConfig
	conn           net.Conn
	connectionOpen bool
	mux            sync.Mutex
	stopped        bool
}

func (dc *UsrCanetDeviceConverter) isStopped() bool {
	dc.mux.Lock()
	defer dc.mux.Unlock()

	return dc.stopped
}

func (dc *UsrCanetDeviceConverter) isConnectionOpen() bool {
	dc.mux.Lock()
	defer dc.mux.Unlock()

	return dc.connectionOpen
}

func (dc *UsrCanetDeviceConverter) stop() bool {
	dc.closeConn()

	dc.mux.Lock()
	defer dc.mux.Unlock()

	dc.stopped = true
	return dc.stopped
}

func (dc *UsrCanetDeviceConverter) initializeConn() bool {
	dc.mux.Lock()
	defer dc.mux.Unlock()
	d := net.Dialer{Timeout: 5 * time.Second}
	conn, err := d.Dial(
		"tcp",
		fmt.Sprintf(
			"%s:%d",
			dc.config.Host,
			dc.config.Port,
		),
	)

	if err != nil {
		log.Error(err.Error())
		return false
	} else {
		log.Info(
			fmt.Sprintf(
				"Connected to Device with Host: %s and Port: %d",
				dc.config.Host,
				dc.config.Port,
			),
		)
	}

	dc.conn = conn
	dc.connectionOpen = true

	return dc.connectionOpen
}

func (dc *UsrCanetDeviceConverter) closeConn() {
	dc.mux.Lock()
	defer dc.mux.Unlock()

	if dc.connectionOpen {
		dc.conn.Close()
		dc.connectionOpen = false
	}

	log.Info(
		fmt.Sprintf(
			"Disconnected from Device with Host: %s and Port: %d",
			dc.config.Host,
			dc.config.Port,
		),
	)
}

func (dc *UsrCanetDeviceConverter) run() {
	for {
		stopped := dc.isStopped()

		if stopped {
			break
		}

		if !dc.isConnectionOpen() {
			connectionOpen := dc.initializeConn()

			if !connectionOpen {
				time.Sleep(time.Second)
				continue
			}
		}

		dc.mux.Lock()
		conn := dc.conn
		dc.mux.Unlock()

		buff := make([]byte, 13)

		_, err := conn.Read(buff)

		if err != nil {
			log.Error(err.Error())
			dc.closeConn()
		}

		// TODO: implement conversion

	}
}
