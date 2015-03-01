# shift_registers

Eases the use of Shift Registers. Can simulate any length of register

## Usage

```ruby
package main

import (
  "github.com/gmalette/shift_registers"
  "github.com/hybridgroup/gobot"
  "github.com/hybridgroup/gobot/platforms/gpio"
  "github.com/hybridgroup/gobot/platforms/raspi"
  "time"
)

func main() {
  gbot := gobot.NewGobot()

  r := raspi.NewRaspiAdaptor("raspi")
  data := gpio.NewLedDriver(r, "led", "16")
  clock := gpio.NewLedDriver(r, "led", "8")
  latch := gpio.NewLedDriver(r, "led", "12")

  reg := shift_registers.NewShiftRegister(
    8,
    data,
    clock,
    latch,
    nil,
  )

  c := true
  work := func() {
    gobot.Every(1*time.Second, func() {
      reg.Write([]bool{c, c, c, c, c, c, c, c})
      c = !c
    })
  }

  robot := gobot.NewRobot("pi",
    []gobot.Connection{r},
    []gobot.Device{data, clock, latch},
    work,
  )

  gbot.AddRobot(robot)

  gbot.Start()
}
```
