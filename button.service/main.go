package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

const t = 1*time.Second + 500*time.Millisecond

var start, end time.Time

func main() {
	pi := raspi.NewAdaptor()
	button := gpio.NewButtonDriver(pi, "10")
	led := gpio.NewLedDriver(pi, "7")

	work := func() {
		var timer *time.Timer
		turnOffLED(led)

		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("button pressed")
			start = time.Now()
			turnOnLED(led)

			timer = time.AfterFunc(t, func() { turnOffLED(led) })
		})

		button.On(gpio.ButtonRelease, func(data interface{}) {
			timer.Stop()
			fmt.Println("button released")
			end = time.Now()
			elapsed := end.Sub(start)

			if elapsed > t {
				turnOffLED(led)
				shutdown()
				os.Exit(0)
			}

			defer func() {
				start = time.Time{}
				end = time.Time{}
				turnOffLED(led)
			}()
		})
	}

	robot := gobot.NewRobot("push button",
		[]gobot.Connection{pi},
		[]gobot.Device{button},
		work,
	)

	robot.Start()
}

func shutdown() error {
	fmt.Println("正在關機...")
	_, err := exec.Command("sudo", "shutdown", "now", "-h").Output()
	if err != nil {
		return err
	}
	return nil
}

func turnOnLED(led *gpio.LedDriver) {
	err := led.On()
	if err != nil {
		fmt.Println("fail to turn LED on", err)
	}
}

func turnOffLED(led *gpio.LedDriver) {
	err := led.Off()
	if err != nil {
		fmt.Println("fail to turn LED on", err)
	}
}
