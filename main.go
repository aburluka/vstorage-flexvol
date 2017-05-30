package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"errors"

	"github.com/jaxxstorm/flexvolume"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "vstorage flexvolume"
	app.Usage = "Mount vstorage volumes in kubernetes using the flexvolume driver"
	app.Commands = flexvolume.Commands(Vstorage{})
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Virtuozzo anonymous programmers crew",
		},
	}
	app.Version = "0.1a"
	app.Run(os.Args)
}

type Vstorage struct{}

func (v Vstorage) Init() flexvolume.Response {
	return flexvolume.Response{
		Status:  flexvolume.StatusSuccess,
		Message: "Vstorage is available",
	}
}

func verifyOptions(options map[string]string) (flexvolume.Response, error) {
	msg := ""
	if options["clusterName"] == "" {
		msg = "Must specify a cluster name"
	}

	if options["clusterPassword"] == "" {
		msg = "Must specify a cluster password"
	}

	if options["hostMountPath"] == "" {
		msg = "Must specify a host mount path path"
	}

	if msg != "" {
		return flexvolume.Response{
			Status:  flexvolume.StatusFailure,
			Message: msg,
		}, errors.New(msg)
	}
	return flexvolume.Response{Status: flexvolume.StatusSuccess}, nil
}

func (v Vstorage) Attach(options map[string]string) flexvolume.Response {
	if resp, err := verifyOptions(options); err != nil {
		return resp
	}

	auth := exec.Command("/usr/bin/vstorage", "-c", options["clusterName"], "auth-node", "-P")
	var b bytes.Buffer
	b.Write([]byte(options["clusterPassword"]))
	auth.Stdout = nil
	auth.Stdin = &b
	auth.Stderr = nil
	if err := auth.Run(); err != nil {
		return flexvolume.Response{
			Status:  flexvolume.StatusFailure,
			Message: "Cannot auth in vstorage",
		}
	}

	return flexvolume.Response{
		Status:  flexvolume.StatusSuccess,
		Message: "Successfully attached the vstorage volume",
		Device: options["hostMountPath"] + "/" + options["clusterName"],
	}
}

func (v Vstorage) Detach(volumeName string) flexvolume.Response {
	return flexvolume.Response{
		Status:  flexvolume.StatusSuccess,
		Message: "Successfully detached the vstorage volume",
		Device:  volumeName,
	}
}

func (v Vstorage) WaitForAttach(volumeName string) flexvolume.Response {
	return flexvolume.Response{
		Status: flexvolume.StatusSuccess,
		Device: volumeName,
	}
}

func (v Vstorage) MountDevice(target string, options map[string]string) flexvolume.Response {
	if resp, err := verifyOptions(options); err != nil {
		return resp
	}

	mount := exec.Command("/usr/bin/vstorage-mount", "-c", "stor1", options["hostMountPath"])
	if err := mount.Run(); err != nil {
		// TODO do not ignore errors
		//return flexvolume.Response{
		//	Status:  flexvolume.StatusFailure,
		//	Message: fmt.Sprintf("Cannot mount vstorage to %s", target),
		//}
	}

	if err:= syscall.Mount(options["hostMountPath"], target, "", syscall.MS_BIND, ""); err != nil {
		return flexvolume.Response{
			Status:  flexvolume.StatusFailure,
			Message: fmt.Sprintf("Cannot bindmount vstorage to %s", target),
		}
	}
	return flexvolume.Response{
		Status:  flexvolume.StatusSuccess,
		Message: "Vstorage volume already mounted",
		Device:  target,
	}
}

func (v Vstorage) Mount(target string, options map[string]string) flexvolume.Response {
	if resp, err := verifyOptions(options); err != nil {
		return resp
	}

	return flexvolume.Response{
		Status:  flexvolume.StatusSuccess,
		Message: "Vstorage mounted",
		Device:  target,
	}
}

func (v Vstorage) Unmount(mount string) flexvolume.Response {
	return flexvolume.Response{
		Status:  flexvolume.StatusSuccess,
		Message: "Successfully unmounted",
		Device:  mount,
	}
}

func (v Vstorage) UnmountDevice(mount string) flexvolume.Response {
	if err := syscall.Unmount(mount, syscall.MNT_DETACH); err != nil {
		return flexvolume.Response{
			Status:  flexvolume.StatusFailure,
			Message: fmt.Sprintf("Cannot unmount %s", mount),
			Device:  mount,
		}
	}
	return flexvolume.Response{
		Status:  flexvolume.StatusSuccess,
		Message: "Successfully unmounted the vstorage volume",
		Device:  mount,
	}
}

func (v Vstorage) GetVolumeName(options map[string]string) flexvolume.Response {
	if resp, err := verifyOptions(options); err != nil {
		return resp
	}

	return flexvolume.Response{
		Status:     flexvolume.StatusSuccess,
		VolumeName: options["hostMountPath"] + "/" + options["clusterName"],
	}
}

func (v Vstorage) IsAttached(options map[string]string) flexvolume.Response {
	if resp, err := verifyOptions(options); err != nil {
		return resp
	}

	return flexvolume.Response{
		Status:   flexvolume.StatusSuccess,
		Attached: true,
	}
}
