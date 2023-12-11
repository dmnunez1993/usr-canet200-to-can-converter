package usrcanettocan

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/go-daq/canbus"
	log "github.com/sirupsen/logrus"
)

// USR Canet 200 sends packages with size of 13 per the documentation
const UsrCanet200PacketSize = 13

// According to datasheet, bit 7 from frame info word defines
// if the package is RTR
const UsrCanetRtrMask = (1 << 6)

// According to datasheet, bit 8 from frame info word defines
// the type of package (EFF or SFF)
const UsrCanetEffMask = (1 << 7)

type UsrCanetDeviceConverter struct {
	config      UsrCanetDeviceConverterConfig
	netConn     net.Conn
	netConnOpen bool
	canConn     *canbus.Socket
	canConnOpen bool
	mux         sync.Mutex
	stopped     bool
}

func (dc *UsrCanetDeviceConverter) isStopped() bool {
	dc.mux.Lock()
	defer dc.mux.Unlock()

	return dc.stopped
}

func (dc *UsrCanetDeviceConverter) isConnOpen() bool {
	dc.mux.Lock()
	defer dc.mux.Unlock()

	return dc.netConnOpen && dc.canConnOpen
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

	dc.netConn = conn
	dc.netConnOpen = true

	dc.canConn, err = canbus.New()

	if err != nil {
		log.Error(err.Error())
		return false
	}

	err = dc.canConn.Bind(dc.config.Target)

	if err != nil {
		log.Error(err.Error())
		return false
	} else {
		log.Info(
			fmt.Sprintf("Connected to can port: %s", dc.config.Target),
		)
	}

	dc.canConnOpen = true

	return dc.netConnOpen && dc.canConnOpen
}

func (dc *UsrCanetDeviceConverter) closeConn() {
	dc.mux.Lock()
	defer dc.mux.Unlock()

	if dc.netConnOpen {
		dc.netConn.Close()
		dc.netConnOpen = false
	}

	log.Info(
		fmt.Sprintf(
			"Disconnected from Device with Host: %s and Port: %d",
			dc.config.Host,
			dc.config.Port,
		),
	)

	if dc.canConnOpen {
		dc.canConn.Close()
		dc.canConnOpen = false
	}
}

func (dc *UsrCanetDeviceConverter) run() {
	for !dc.isStopped() {

		if !dc.isConnOpen() {
			netConnOpen := dc.initializeConn()

			if !netConnOpen {
				time.Sleep(time.Second)
				continue
			}
		}

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go dc.netToCan(wg)
		go dc.canToNet(wg)

		wg.Wait()

		time.Sleep(time.Second)
	}
}

func (dc *UsrCanetDeviceConverter) netToCan(wg *sync.WaitGroup) {
	defer wg.Done()
	for !dc.isStopped() {
		dc.mux.Lock()
		netConn := dc.netConn
		dc.mux.Unlock()

		buff := make([]byte, UsrCanet200PacketSize)

		_, err := netConn.Read(buff)

		if err != nil {
			log.Error(err.Error())
			dc.closeConn()
			return
		}

		frameInfo := buff[0]
		frameId := buff[1:5]
		payload := buff[5:]

		msg := canbus.Frame{
			ID:   uint32(frameId[3]) + uint32(frameId[2])<<8 + uint32(frameId[1])<<16 + uint32(frameId[0])<<24,
			Data: payload,
		}

		if int(frameInfo)&UsrCanetEffMask > 0 {
			msg.Kind = canbus.EFF
		} else {
			msg.Kind = canbus.SFF
		}

		if int(frameInfo)&UsrCanetRtrMask > 0 {
			msg.Kind = canbus.RTR
		}

		dc.mux.Lock()
		canConn := dc.canConn
		dc.mux.Unlock()

		_, err = canConn.Send(msg)

		if err != nil {
			log.Error(err.Error())
			dc.closeConn()
			return
		}
	}
}

func (dc *UsrCanetDeviceConverter) canToNet(wg *sync.WaitGroup) {
	defer wg.Done()

	for !dc.isStopped() {
		dc.mux.Lock()
		canConn := dc.canConn
		dc.mux.Unlock()

		msg, err := canConn.Recv()

		if err != nil {
			log.Error(err.Error())
			dc.closeConn()
			return
		}

		// Set the length of the package
		frameInfo := byte(len(msg.Data))

		if msg.Kind == canbus.RTR {
			// Setting bit 6 of the frame info word indicates a RTR message
			frameInfo = frameInfo | 0x6
		}

		if msg.Kind == canbus.EFF {
			// Setting bit 7 of the frame info word indicates an EFF message
			frameInfo = frameInfo | 0x7
		}

		frameId := make([]byte, 4)

		// Split ID into bytes
		frameId[0] = byte(msg.ID >> 24 & 0xFF)
		frameId[1] = byte(msg.ID >> 16 & 0xFF)
		frameId[2] = byte(msg.ID >> 8 & 0xFF)
		frameId[3] = byte(msg.ID & 0xFF)

		payload := msg.Data

		netPackage := make([]byte, 13)
		// Add frame info
		netPackage[0] = frameInfo
		// Add frame id
		copy(netPackage[1:5], frameId)
		// Add payload
		copy(netPackage[6:], payload)

		dc.mux.Lock()
		netConn := dc.netConn
		dc.mux.Unlock()

		_, err = netConn.Write(netPackage)

		if err != nil {
			log.Error(err.Error())
			dc.closeConn()
			return
		}
	}
}
