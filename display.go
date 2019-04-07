package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	evdev "github.com/gvalkov/golang-evdev"
)

// #cgo LDFLAGS: -lX11
// #include <stdlib.h>
// #include <stdio.h>
// #include <err.h>
// #include <X11/Xlib.h>
// #include <X11/XKBlib.h>
import "C"

type Event struct {
	//Device *evdev.InputDevice
	Items []evdev.InputEvent
}

func openDisplay() (*C.Display, error) {
	var xkbEventType, xkbError, xkbReason C.int
	var majorVers, minorVers C.int

	majorVers = C.XkbMajorVersion
	minorVers = C.XkbMinorVersion

	display := C.XkbOpenDisplay(
		nil, &xkbEventType, &xkbError, &majorVers, &minorVers, &xkbReason,
	)
	if display == nil {
		switch xkbReason {
		case C.XkbOD_BadServerVersion:
		case C.XkbOD_BadLibraryVersion:
			return nil, fmt.Errorf("incompatible versions of client and server XKB libraries")
		case C.XkbOD_ConnectionRefused:
			return nil, fmt.Errorf("connection to X server refused")
		case C.XkbOD_NonXkbServer:
			return nil, fmt.Errorf("XKB extension is not present")
		default:
			return nil, fmt.Errorf("unknown error from XkbOpenDisplay: %d", xkbReason)
		}
	}

	return display, nil
}

func closeDisplay(display *C.Display) {
	C.XCloseDisplay(display)
}

func listenDevice(
	//name string,
	device *evdev.InputDevice,
	inbox chan Event,
) {
	for {
		events, err := device.Read()
		if err != nil || len(events) == 0 {
			// device lost
			return
		}

		inbox <- Event{Items: events}
	}
}

func getInputDevices() map[string]*evdev.InputDevice {
	inputDevices := make(map[string]*evdev.InputDevice)

	devicePaths, err := evdev.ListInputDevicePaths("/dev/input/event*")
	if err == nil && len(devicePaths) > 0 {
		for _, devicePath := range devicePaths {
			device, err := evdev.Open(devicePath)
			if err != nil {
				log.Printf("unable to open device %s: %s", devicePath, err)

				continue
			}

			inputDevices[devicePath] = device
		}
	}

	return inputDevices
}

func watchKeyPress() (up, down chan string, err error) {
	display, err := openDisplay()
	if err != nil {
		panic(err)
	}

	defer closeDisplay(display)

	events := make(chan Event, 8)

	devices := getInputDevices()
	if len(devices) == 0 {
		return nil, nil, errors.New("unable to open devices (requires root privileges)")
	}

	for _, device := range devices {
		go listenDevice(
			device,
			events,
		)
	}

	up = make(chan string)
	down = make(chan string)

	go func() {
		for {
			select {
			case event := <-events:
				for _, item := range event.Items {
					if item.Type != evdev.EV_KEY {
						continue
					}

					if item.Value == 2 {
						continue
					}

					key := strings.TrimPrefix(evdev.KEY[int(item.Code)], "KEY_")

					switch item.Value {
					case 0:
						up <- key
					case 1:
						down <- key
					}
				}
			}
		}
	}()

	return up, down, nil
}
