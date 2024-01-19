# arduino-util

Build your Arduino projects with a Makefile and pinned, local libraries.

## Introduction

There are many things that make me uncomfortable in the various standard Arduino development processes. I don’t like typing in the Arduino IDE. I don’t like constantly having to redefine which port my Arduino is available at because I don’t always plug the same device into the same USB port of the same computer. I don’t like having to install libraries globally, and then having to guess which project was working with which version of which library, based only on the `#include` directives. I like having a `Makefile` with all the common tasks for my project, so I don’t need to remember them, and I like the dependencies to be local to the project directory, in the exact version I developed my code with, possibly, but not necessarily, as a git submodule.

The new `arduino-cli` command line tool is a big improvement for me, as I can now compile and upload my projects from a `Makefile`, but it doesn’t solve the port or libraries problems. Build profiles let you specify the libraries your project depends on, and their version, but it’s limited to those installable from the library manager: if you need to fork a library, or if you developed one for your own personal use, you not only can’t specify them in a build profile, you can’t use build profiles altogether. Moreover, even if your project only uses libraries from the library manager, you’re still dependent on the constant availability of the libraries you use there (I know it’s the same thing with `pip` or `npm`, but there are much more people working on those to prevent such problems from happening.)

`arduino-util` is a small tool that solves these problems for me. It helps with two things: it can create a `Makefile` that compiles and uploads an Arduino project with local copies of its dependent libraries, and it can automatically find the port where an Arduino is plugged in for uploading or launching the serial monitor.

## Usage

Download the right version of `arduino-util` for your OS and architecture, and make sure it is executable and available in your path. You can rename it or symlink it, it will use the name you use to call it in the `Makefile` it generates. Personally, I symlinked it to `_arduino-util` so it doesn’t compete with `arduino-cli` in my shell completion.

### Generating a Makefile

Type `arduino-util makefile`, and it will print a `Makefile` to compile an Arduino project:


```Makefile
# Generated by 'arduino-util makefile'

# Path to the folder containing the libraries for your project.
LIB_PATH = libraries
# Space-delimited list of libraries from LIB_PATH.
LIBRARIES =
BUILD_PATH = build

BOARD_FQBN = # e.g. arduino:avr:nano
BOARD_OPTIONS = # e.g. cpu=atmega328old

# Adapt this regex to match the range of names your Arduino might appear in /dev/,
# or leave empty for the default: cu\.usb(serial|modem)
PORT_REGEX =
BAUD_RATE = 115200

### End of configuration ###

LIB_OPTS = $(addprefix --library $(LIB_PATH)/,$(LIBRARIES))
COMPILE_OPTS = \
	--fqbn "$(BOARD_FQBN)" \
	--board-options "$(BOARD_OPTIONS)" \
	$(LIB_OPTS) \
	--build-path "$(BUILD_PATH)"
UPLOAD_OPTS = \
	--fqbn "$(BOARD_FQBN)" \
	--board-options "$(BOARD_OPTIONS)" \
	--input-dir "$(BUILD_PATH)"

PORT_ARG = $(if $(PORT_REGEX),-regex '$(PORT_REGEX)',)

compile:
	arduino-cli compile $(COMPILE_OPTS)

find-board:
	$(eval PORT := $(shell arduino-util find-board $(PORT_ARG)))
	$(if $(PORT),,$(error "Stopping"))

upload: find-board
	arduino-cli upload $(UPLOAD_OPTS) --port "$(PORT)"

monitor: find-board
	arduino-cli monitor --config baudrate=$(BAUD_RATE) --port "$(PORT)"
```

Save this to as your project’s `Makefile` and configure the options:

* `LIB_PATH`: relative or absolute path to a folder containing the libraries used by your project;
* `LIBRARIES`: space-delimited list of libraries used by your project, found at the root of the `LIB_PATH` folder; libraries cannot have spaces in their folder name;
* `BUILD_PATH`: path to the build folder, so you can easily go peek at what the compiler does;
* `BOARD_FQBN`: the board FQBN of your Arduino;
* `BOARD_OPTIONS`: optional board options, passed to `compile` and `upload` via the `--board-options` option;
* `PORT_REGEX`: regular expression that the port of the Arduino you want to program should match, for automatic port detection; leave blank for the default;
* `BAUD_RATE`: the baud rate at which your sketch communicates on the serial port, for the serial monitor.

To automatically find the port for `arduino-cli upload` or `arduino-cli monitor`, the `Makefile` calls `arduino-util find-board` with an optional `-regex` option. `arduino-util find-board` then lists `/dev/` and filters devices that match the regex. If there is exactly one result, the search is successful and it passes it to `arduino-cli` via the `--port` option. If there are no or several results, the `Makefile` exits with an error.

The default value for `PORT_REGEX` is `cu\.usb(serial|modem)`, because I own Arduino Nano clones and an Arduino Nano Every, and they appear as `/dev/cu.usbserial-nnn` and `/dev/cu.usbmodem-nnn`, respectively.

### Compiling the Project

Type `make compile`. The project is built in the `build` directory by default.

### Uploading to the Arduino

Plug in your Arduino and type `make upload`. If it’s the only device that matches the `PORT_REGEX`, automatic board detection will be successful and the sketch will be uploaded. Otherwise, it will stop with the appropriate error message:

```
No device matching 'cu\.usb(serial|modem)' found in /dev/.
Makefile:36: *** "Stopping".  Stop.
```

or

```
More than one device matching 'cu' found in /dev/:
  cu.Bluetooth-Incoming-Port
  cu.usbserial-42
  cu.wlan-debug
Makefile:36: *** "Stopping".  Stop.
```

### Launching the Serial Monitor

Plug in your Arduino and type `make monitor`. Automatic board detection will be performed the same way as for uploading, and if successful, `arduino-cli monitor` will be started with the appropriate options.

## License

<p xmlns:cc="http://creativecommons.org/ns#" xmlns:dct="http://purl.org/dc/terms/"><a property="dct:title" rel="cc:attributionURL" href="https://github.com/oscherler/arduino-util">arduino-util</a> by <a rel="cc:attributionURL dct:creator" property="cc:attributionName" href="https://github.com/oscherler">Olivier Scherler</a> is licensed under <a href="http://creativecommons.org/licenses/by-sa/4.0/?ref=chooser-v1" target="_blank" rel="license noopener noreferrer" style="display:inline-block;">CC BY-SA 4.0<img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;" src="https://mirrors.creativecommons.org/presskit/icons/cc.svg?ref=chooser-v1"><img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;" src="https://mirrors.creativecommons.org/presskit/icons/by.svg?ref=chooser-v1"><img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;" src="https://mirrors.creativecommons.org/presskit/icons/sa.svg?ref=chooser-v1"></a></p>