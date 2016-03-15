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

  if len(os.Args) != 3 {
    fmt.Printf("usage: %s <PIN> <0|1|t>\n", os.Args[0])
    fmt.Println("  0: low")
    fmt.Println("  1: high")
    fmt.Println("  t: toggle pin")
    os.Exit(1)
  }


  var nr int
  if nr, err = strconv.Atoi(os.Args[1]); err != nil {
    fmt.Printf("invalid argument for 'PIN' - %s\n", err.Error())
    os.Exit(1)
  }

  var action func()
  var pin gpioctl.Pin
  if os.Args[2] == "t" {
    action = func(){
      pin.Toggle()
    }
  } else if value, err := strconv.Atoi(os.Args[2]); err == nil {
    action = func(){
      pin.Write(value)
    }
  } else {
    fmt.Printf("invalid argument for 'VALUE' - expected: <0|1> - %s\n", err.Error())
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

  if pin, err = gpio.Pin(nr, gpioctl.GPIO_PIN_OUTPUT); err != nil {
    fmt.Printf("unable to open gpio device - %s\n", err.Error())
    os.Exit(1)
  }

  action()
}
