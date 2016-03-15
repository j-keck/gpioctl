package gpioctl


// #include <sys/cdefs.h>
// #include <paths.h>
// #include <stdlib.h>
// #include <unistd.h>
// #include <sys/gpio.h>
import "C"
import "syscall"
import "unsafe"


type Gpio struct {
  fd int;
}

type Pin struct {
  nr   uint32;
  gpio Gpio;
}
type PinFlags uint32

const (
  GPIO_PIN_INPUT       PinFlags    = 0x0001;
  GPIO_PIN_OUTPUT      PinFlags    = 0x0002;
  GPIO_PIN_OPENDRAIN   PinFlags    = 0x0004;
  GPIO_PIN_PUSHPULL    PinFlags    = 0x0008;
  GPIO_PIN_TRISTATE    PinFlags    = 0x0010;
  GPIO_PIN_PULLUP      PinFlags    = 0x0020;
  GPIO_PIN_PULLDOWN    PinFlags    = 0x0040;
  GPIO_PIN_INVIN       PinFlags    = 0x0080;
  GPIO_PIN_INVOUT      PinFlags    = 0x0100;
  GPIO_PIN_PULSATE     PinFlags    = 0x0200;
)



func Open() (Gpio, error) {
  path := C._PATH_DEVGPIOC + "0"
  if fd, err := syscall.Open(path, syscall.O_RDONLY, 0); err != nil {
    return Gpio{}, err
  } else {
    return Gpio{fd}, nil
  }
}


func (gpio Gpio) Close() error {
  return syscall.Close(gpio.fd)
}


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


func (pin Pin) Write(value int) error {
  var gpioreq gpio_req
  gpioreq.gp_pin = pin.nr
  gpioreq.gp_value = uint32(value)
  return ioctl(uintptr(pin.gpio.fd), C.GPIOSET, uintptr(unsafe.Pointer(&gpioreq)))
}

func (pin Pin) Toggle() error {
  var gpioreq gpio_req
  gpioreq.gp_pin = pin.nr
  return ioctl(uintptr(pin.gpio.fd), C.GPIOTOGGLE, uintptr(unsafe.Pointer(&gpioreq)))
}

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
