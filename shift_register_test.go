package shift_registers

import (
	_ "fmt"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"reflect"
	"testing"
)

type gpioTestAdaptor struct {
	name string
	port string
}

var datum []byte
var latched bool
var state byte

var testAdaptorDigitalWrite = func(p string, d byte) (err error) {
	if p == "1" {
		state = d
	} else if p == "2" && d&1 == 1 {
		datum = append([]byte{state}, datum[0:3]...)
	} else if d&1 == 1 {
		latched = true
	}
	return nil
}

func (t *gpioTestAdaptor) DigitalWrite(s string, b byte) (err error) {
	return testAdaptorDigitalWrite(s, b)
}
func (t *gpioTestAdaptor) ServoWrite(string, byte) (err error)     { return nil }
func (t *gpioTestAdaptor) PwmWrite(string, byte) (err error)       { return nil }
func (t *gpioTestAdaptor) AnalogRead(string) (val int, err error)  { return 0, nil }
func (t *gpioTestAdaptor) DigitalRead(string) (val int, err error) { return 0, nil }
func (t *gpioTestAdaptor) Connect() (errs []error)                 { return }
func (t *gpioTestAdaptor) Finalize() (errs []error)                { return }
func (t *gpioTestAdaptor) Name() string                            { return t.name }
func (t *gpioTestAdaptor) Port() string                            { return t.port }

func newGpioTestAdaptor(name string) *gpioTestAdaptor {
	return &gpioTestAdaptor{
		name: name,
		port: "/dev/null",
	}
}
func initTestLedDriver(conn gpio.DigitalWriter, pin string) *gpio.LedDriver {
	return gpio.NewLedDriver(conn, "bot", pin)
}

func initShiftRegister() *shiftRegister {
	datum = []byte{1, 1, 1, 1}
	state = 1
	latched = false

	adapter := newGpioTestAdaptor("adaptor")
	data := initTestLedDriver(adapter, "1")
	clock := initTestLedDriver(adapter, "2")
	latch := initTestLedDriver(adapter, "3")

	register := NewShiftRegister(
		4,
		data,
		clock,
		latch,
		nil,
	)

	return register
}

func TestShiftRegisterInitializationClears(t *testing.T) {
	initShiftRegister()

	expected := []byte{0, 0, 0, 0}
	if !reflect.DeepEqual(datum, expected) {
		t.Errorf("Bytes were not reset, was %+v", datum)
	}

	if !latched {
		t.Error("Data was not latched")
	}
}

func TestShiftRegisterWrite(t *testing.T) {
	register := initShiftRegister()
	latched = false
	register.Write([]bool{true, false, true, true})

	if !reflect.DeepEqual(datum, []byte{1, 0, 1, 1}) {
		t.Errorf("Bytes were not written properly, was %+v", datum)
	}

	if !latched {
		t.Error("Data was not latched")
	}
}

func TestShiftRegisterWriteFillsData(t *testing.T) {
	register := initShiftRegister()
	register.Write([]bool{true, false})

	if !reflect.DeepEqual(datum, []byte{1, 0, 0, 0}) {
		t.Errorf("Bytes were not written properly, was %+v", datum)
	}
}

func TestShiftRegisterWriteTruncatesData(t *testing.T) {
	register := initShiftRegister()
	register.Write([]bool{true, false, false, false, true, false})

	if !reflect.DeepEqual(datum, []byte{1, 0, 0, 0}) {
		t.Errorf("Bytes were not written properly, was %+v", datum)
	}
}
