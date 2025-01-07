//go:build !linux
// +build !linux

package desko

import (
	"errors"

	"github.com/bearsh/hid"
)

// GetDeviceInfo - returns HID device info
func GetDeviceInfo() (*hid.DeviceInfo, error) {
	for _, d := range hid.Enumerate(deskoUsbVendorID, deskoUsbProductID) {
		if d.Usage == 1 && d.Interface == 2 {
			return &d, nil
		}
	}
	return nil, errors.New("no supported DESKO device found")
}
