// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/siderolabs/go-copy/copy"
	"github.com/siderolabs/talos/pkg/machinery/overlay"
	"github.com/siderolabs/talos/pkg/machinery/overlay/adapter"
	"golang.org/x/sys/unix"
)

const (
	off int64 = 512 * 64
	dtb       = "rockchip/rk3588-orangepi-5-max.dtb"
)

func main() {
	adapter.Execute(&opi5MaxInstaller{})
}

type opi5MaxInstaller struct{}

type opi5MaxExtraOptions struct {
	SPIBoot         bool     `yaml:"spi_boot,omitempty"`
	ExtraKernelArgs []string `yaml:"extraKernelArgs,omitempty"`
}

func (i *opi5MaxInstaller) GetOptions(extra opi5MaxExtraOptions) (overlay.Options, error) {
	kernelArgs := []string{
		"console=tty1",
		"console=ttyS2,1500000",
		"sysctl.kernel.kexec_load_disabled=1",
		"talos.dashboard.disabled=1",
	}

	kernelArgs = append(kernelArgs, extra.ExtraKernelArgs...)

	return overlay.Options{
		Name:       "orangepi-5-max",
		KernelArgs: kernelArgs,
		PartitionOptions: overlay.PartitionOptions{
			Offset: 2048 * 10,
		},
	}, nil
}

func (i *opi5MaxInstaller) Install(options overlay.InstallOptions[opi5MaxExtraOptions]) error {
	if !options.ExtraOptions.SPIBoot {
		uBootBin := filepath.Join(options.ArtifactsPath, "arm64/u-boot/orangepi-5-max/u-boot-rockchip.bin")

		if err := uBootLoaderInstall(uBootBin, options.InstallDisk); err != nil {
			return err
		}
	}

	src := filepath.Join(options.ArtifactsPath, "arm64/dtb", dtb)
	dst := filepath.Join(options.MountPrefix, "boot/EFI/dtb", dtb)

	if err := copyFileAndCreateDir(src, dst); err != nil {
		return err
	}

	return nil
}

func copyFileAndCreateDir(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o600); err != nil {
		return err
	}

	return copy.File(src, dst)
}

func uBootLoaderInstall(uBootBin, installDisk string) error {
	f, err := os.OpenFile(installDisk, os.O_RDWR|unix.O_CLOEXEC, 0o666)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", installDisk, err)
	}

	defer f.Close() //nolint:errcheck

	uboot, err := os.ReadFile(uBootBin)
	if err != nil {
		return err
	}

	if _, err = f.WriteAt(uboot, off); err != nil {
		return err
	}

	// NB: In the case that the block device is a loopback device, we sync here
	// to ensure that the file is written before the loopback device is
	// unmounted.
	return f.Sync()
}
