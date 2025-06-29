package mhz19

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/tarm/serial"
)

const (
	cmdReadCo2    = 0x86
	cmdCalibrate  = 0x87
	cmdAutoCalOn  = 0x79
	cmdAutoCalOff = 0x79
)

type MHZ19 struct {
	port *serial.Port
}

func New(portName string) (*MHZ19, error) {
	config := &serial.Config{
		Name:        portName,
		Baud:        9600,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
		ReadTimeout: time.Second * 2,
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, fmt.Errorf("unable to open port %s: %v", portName, err)
	}

	return &MHZ19{port: port}, nil
}

func (m *MHZ19) Close() error {
	return m.port.Close()
}

func calculateChecksum(data []byte) byte {
	var sum byte = 0
	for i := 1; i < len(data)-1; i++ {
		sum += data[i]
	}
	return 0xFF - sum + 1
}

func (m *MHZ19) ReadCO2() (int, error) {
	cmd := []byte{0xFF, 0x01, cmdReadCo2, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	cmd[8] = calculateChecksum(cmd)

	_, err := m.port.Write(cmd)
	if err != nil {
		return 0, fmt.Errorf("error while sending a command: %v", err)
	}

	response := make([]byte, 9)
	n, err := m.port.Read(response)
	if err != nil {
		return 0, fmt.Errorf("response reading error: %v", err)
	}

	if n != 9 {
		return 0, fmt.Errorf("%d bits, instead of 9", n)
	}

	if response[0] != 0xFF || response[1] != 0x86 {
		return 0, fmt.Errorf("wrong header: %s", hex.EncodeToString(response))
	}

	expectedChecksum := calculateChecksum(response)
	if response[8] != expectedChecksum {
		return 0, fmt.Errorf("wrong checksum")
	}

	co2 := int(response[2])*256 + int(response[3])
	return co2, nil
}

func (m *MHZ19) CalibrateZero() error {
	cmd := []byte{0xFF, 0x01, cmdCalibrate, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	cmd[8] = calculateChecksum(cmd)

	_, err := m.port.Write(cmd)
	if err != nil {
		return fmt.Errorf("calibrate error: %v", err)
	}

	fmt.Println("calibrate successfully done")
	return nil
}

func (m *MHZ19) SetAutoCalibration(enable bool) error {
	var cmd []byte
	if enable {
		cmd = []byte{0xFF, 0x01, cmdAutoCalOn, 0xA0, 0x00, 0x00, 0x00, 0x00, 0x00}
	} else {
		cmd = []byte{0xFF, 0x01, cmdAutoCalOff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	}
	cmd[8] = calculateChecksum(cmd)

	_, err := m.port.Write(cmd)
	if err != nil {
		return fmt.Errorf("auto calibrate set error: %v", err)
	}

	if enable {
		fmt.Println("auto calibrate on")
	} else {
		fmt.Println("auto calibrate off")
	}
	return nil
}
