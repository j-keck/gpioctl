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

  if len(os.Args) < 2 {
    fmt.Printf("usage: %s <PIN> [<PIN_NAME>]\n", os.Args[0])
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


  if len(os.Args) == 2 {
    if name, err := gpio.GetPinName(nr); err == nil {
      fmt.Println(name)
    } else {
      fmt.Printf("unable to get pin name - %s\n", err.Error())
    }
  } else {
    if err := gpio.SetPinName(nr, os.Args[2]); err != nil {
      fmt.Printf("unable to set pin name - %s\n", err.Error())
    }
  }
}
