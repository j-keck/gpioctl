// Package gpioctl to manage GPIO pins under bsd
//
//
// gpioctl use the 'ioctl()' system call to manage
// the gpio pins.
//
//
// EXAMPLE:
//
//   open pin 40 in output mode and set it to high:
//
//     if gpio, err := gpioctl.Open(); err == nil {
//       defer gpio.Close()
//
//       if pin, err := gpio.Pin(40, gpioctl.GPIO_PIN_OUTPUT); err == nil {
//         pin.Write(1)
//       } else {
//         panic(err)
//       }
//     } else {
//       panic(err)
//     }
//
package gpioctl


// #include <sys/cdefs.h>
// #include <paths.h>
// #include <stdlib.h>
// #include <unistd.h>
// #include <sys/gpio.h>
import "C"
import "syscall"
import "unsafe"


// Gpio represents the GPIO bus
type Gpio struct {
  fd int;
}

// Pin represents one GPIO Pin
type Pin struct {
  nr   uint32;
  gpio Gpio;
}


// GPIO pin configuration flags
type PinFlags uint32

const (
  // input direction
  GPIO_PIN_INPUT       PinFlags    = 0x0001;
  // output direction
  GPIO_PIN_OUTPUT      PinFlags    = 0x0002;
  // open-drain output
  GPIO_PIN_OPENDRAIN   PinFlags    = 0x0004;
  // push-pull output
  GPIO_PIN_PUSHPULL    PinFlags    = 0x0008;
  // output disabled
  GPIO_PIN_TRISTATE    PinFlags    = 0x0010;
  // internal pull-up enabled
  GPIO_PIN_PULLUP      PinFlags    = 0x0020;
  // internal pull-down enabled
  GPIO_PIN_PULLDOWN    PinFlags    = 0x0040;
  // invert input
  GPIO_PIN_INVIN       PinFlags    = 0x0080;
  // invert output
  GPIO_PIN_INVOUT      PinFlags    = 0x0100;
  // pulsate in hardware
  GPIO_PIN_PULSATE     PinFlags    = 0x0200;
)



// open the GPIO bus
func Open() (Gpio, error) {
  path := C._PATH_DEVGPIOC + "0"
  if fd, err := syscall.Open(path, syscall.O_RDONLY, 0); err != nil {
    return Gpio{}, err
  } else {
    return Gpio{fd}, nil
  }
}


// close the GPIO bus
func (gpio Gpio) Close() error {
  return syscall.Close(gpio.fd)
}


// configure a GPIO pin
func (gpio Gpio) Pin(nr int, flags ...PinFlags) (Pin, error) {
 var gpiopin gpio_pin
  gpiopin.gp_pin = uint32(nr)
  for _, f := range flags {
    gpiopin.gp_flags |= uint32(f)
  }

  if err := ioctl(uintptr(gpio.fd), C.GPIOSETCONFIG, uintptr(unsafe.Pointer(&gpiopin))); err != nil {
    return Pin{}, err
  }

  return Pin{uint32(nr), gpio}, nil
}


// set the name for the given pin
func (gpio Gpio) SetPinName(nr int, name string) error {
  var gpiopin gpio_pin
  gpiopin.gp_pin = uint32(nr)

  cname := unsafe.Pointer(C.CString(name))
  defer C.free(cname)
  gpiopin.gp_name = *(*[C.GPIOMAXNAME]C.char)(cname)

  return ioctl(uintptr(gpio.fd), C.GPIOSETNAME, uintptr(unsafe.Pointer(&gpiopin)))
}


// set the name for the given pin
func (pin Pin) SetName(name string) error {
  return pin.gpio.SetPinName(int(pin.nr), name)
}


// get the pin name
func (gpio Gpio) GetPinName(nr int) (string, error) {
  var gpiopin gpio_pin
  gpiopin.gp_pin = uint32(nr)

  if err := ioctl(uintptr(gpio.fd), C.GPIOGETCONFIG, uintptr(unsafe.Pointer(&gpiopin))); err == nil {
    return C.GoString(&gpiopin.gp_name[0]), nil
  } else {
    return "", err
  }
}


// get the pin name
func (pin Pin) GetName() (string, error) {
  return pin.gpio.GetPinName(int(pin.nr))
}


// write the given value. 0 -> LOW, >= 1 -> HIGH
func (pin Pin) Write(value int) error {
  var gpioreq gpio_req
  gpioreq.gp_pin = pin.nr
  gpioreq.gp_value = uint32(value)
  return ioctl(uintptr(pin.gpio.fd), C.GPIOSET, uintptr(unsafe.Pointer(&gpioreq)))
}


// toggle the current pin state
func (pin Pin) Toggle() error {
  var gpioreq gpio_req
  gpioreq.gp_pin = pin.nr
  return ioctl(uintptr(pin.gpio.fd), C.GPIOTOGGLE, uintptr(unsafe.Pointer(&gpioreq)))
}


// read the current pin state
func (pin Pin) Read() (int, error) {
  var gpioreq gpio_req
  gpioreq.gp_pin = pin.nr
  if err := ioctl(uintptr(pin.gpio.fd), C.GPIOGET, uintptr(unsafe.Pointer(&gpioreq))); err != nil {
    return 0, err
  } else {
    return int(gpioreq.gp_value), nil
  }
}


type gpio_pin struct {
  gp_pin   uint32;
  gp_name  [64]C.char;
  gp_caps  uint32;
  gp_flags uint32;
}


type gpio_req struct {
  gp_pin   uint32;
  gp_value uint32;
}


func ioctl(fd, req, arg uintptr) error {
  if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, req, arg); err != 0 {
    return syscall.Errno(err)
  }
  return nil
}
