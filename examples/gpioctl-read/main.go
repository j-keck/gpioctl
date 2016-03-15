package main

import "strconv"
import "os"
import "fmt"
import "github.com/j-keck/gpioctl"

func main(){
  var err error

  //
  // parse / interpret args
  //

  if len(os.Args) != 2 {
    fmt.Printf("usage: %s <PIN>\n", os.Args[0])
    os.Exit(1)
  }


  var nr int
  if nr, err = strconv.Atoi(os.Args[1]); err != nil {
    fmt.Printf("invalid argument for 'PIN' - %s\n", err.Error())
    os.Exit(1)
  }


  //
  // action
  //

  var gpio gpioctl.Gpio
  if gpio, err = gpioctl.Open(); err != nil {
    fmt.Printf("unable to open gpio device - %s\n", err.Error())
    os.Exit(1)
  }
  defer gpio.Close()

  var pin gpioctl.Pin
  if pin, err = gpio.Pin(nr, gpioctl.GPIO_PIN_INPUT); err != nil {
    fmt.Printf("unable to open gpio device - %s\n", err.Error())
    os.Exit(1)
  }

  var value int
  if value, err = pin.Read(); err != nil {
    fmt.Printf("unable to read from pin - %s\n", err.Error())
    os.Exit(1)
  }

  fmt.Printf("pin: %d, current value: %d\n", nr, value)
}
