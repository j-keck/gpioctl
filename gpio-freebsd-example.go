package main

import "strconv"
import "os"
import "fmt"

func main(){
  var err error

  //
  // parse args
  //

  if len(os.Args) != 3 {
    fmt.Printf("usage: %s <PIN> <0|1>\n", os.Args[0])
    os.Exit(1)
  }

  var nr int
  if nr, err = strconv.Atoi(os.Args[1]); err != nil {
    fmt.Printf("invalid argument for 'PIN' - %s\n", err.Error())
    os.Exit(1)
  }

  var value int
  if value, err = strconv.Atoi(os.Args[2]); err != nil {
    fmt.Printf("invalid argument for 'VALUE' - expected: <0|1> - %s\n", err.Error())
    os.Exit(1)
  }



  //
  // action
  //

  var gpio Gpio
  if gpio, err = Open(); err != nil {
    fmt.Printf("unable to open gpio device - %s\n", err.Error())
    os.Exit(1)
  }
  defer gpio.Close()


  var pin Pin
  if pin, err = gpio.Pin(nr, GPIO_PIN_OUTPUT); err != nil {
    fmt.Printf("unable to open gpio device - %s\n", err.Error())
    os.Exit(1)
  }

  pin.Write(value)
}
