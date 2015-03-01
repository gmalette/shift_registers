package shift_registers

import (
	_ "fmt"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

type shiftRegister struct {
	dataPin  *gpio.LedDriver // ser
	clockPin *gpio.LedDriver // srclk
	latchPin *gpio.LedDriver // rclk
	clearPin *gpio.LedDriver // srclr
	lineSize int
}

func NewShiftRegister(lineSize int, ser *gpio.LedDriver, srclk *gpio.LedDriver, rclk *gpio.LedDriver, srclr *gpio.LedDriver) *shiftRegister {
	register := &shiftRegister{
		lineSize: lineSize,
		dataPin:  ser,
		clockPin: srclk,
		latchPin: rclk,
		clearPin: srclr,
	}

	register.dataPin.Off()
	register.clockPin.Off()
	register.latchPin.Off()
	register.Clear()

	return register
}

func (s *shiftRegister) Write(data []bool) {
	missingLen := s.lineSize - len(data)
	if missingLen < 0 {
		missingLen = 0
	}
	actualData := append(make([]bool, missingLen, missingLen), data...)
	actualData = actualData[(len(actualData) - s.lineSize):]

	for index := range actualData {
		pin := actualData[len(actualData)-index-1]
		if pin {
			s.dataPin.On()
		}
		s.clockPin.On()
		s.clockPin.Off()
		s.dataPin.Off()
	}

	s.latchPin.On()
	s.latchPin.Off()
}

func (s *shiftRegister) Clear() {
	if s.clearPin != nil {
	} else {
		empty := make([]bool, 0, 0)
		s.Write(empty)
	}
}
