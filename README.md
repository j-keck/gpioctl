# gpioctl - manage GPIO pins under freebsd

  * simple go library to manage GPIO pins under freebsd
  * it's uses the 'ioctl()' system call to manage the gpio pins.
  * tested on 'FreeBSD rpi2 11.0-CURRENT FreeBSD 11.0-CURRENT #0 r296485'
    - doesn't work under FreeBSD 10 - missing ioctl to name a pin (GPIOSETNAME)
    - 'GPIOSETNAME' was implemented in r279761 (git sha: 04b6b93)
  * no interrupts - currently not supported under freebsd
  * check the 'examples' directory for examples
