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
	off   int64 = 512 * 64
	board       = "nanopi-r5s"
	dtb         = "rockchip/rk3568-nanopi-r5s.dtb"
)

func main() {
	adapter.Execute(&nanopir5s{})
}

type nanopir5s struct{}

type nanopir5sExtraOptions struct{}

func (i *nanopir5s) GetOptions(extra nanopir5sExtraOptions) (overlay.Options, error) {
	return overlay.Options{
		Name: "nanopi-r5s",
		KernelArgs: []string{
			"console=tty0",
			"console=ttyS2,1500000n8",
			"sysctl.kernel.kexec_load_disabled=1",
			"talos.dashboard.disabled=1",
		},
		PartitionOptions: overlay.PartitionOptions{
			Offset: 2048 * 10,
		},
	}, nil
}

func (i *nanopir5s) Install(options overlay.InstallOptions[nanopir5sExtraOptions]) error {
	f, err := os.OpenFile(options.InstallDisk, os.O_RDWR|unix.O_CLOEXEC, 0o666)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", options.InstallDisk, err)
	}

	defer f.Close() //nolint:errcheck

	uboot, err := os.ReadFile(filepath.Join(options.ArtifactsPath, "arm64/u-boot", board, "u-boot-rockchip.bin"))
	if err != nil {
		return err
	}

	if _, err = f.WriteAt(uboot, off); err != nil {
		return err
	}

	// NB: In the case that the block device is a loopback device, we sync here
	// to ensure that the file is written before the loopback device is
	// unmounted.
	err = f.Sync()
	if err != nil {
		return err
	}

	src := filepath.Join(options.ArtifactsPath, "arm64/dtb", dtb)
	dst := filepath.Join(options.MountPrefix, "/boot/EFI/dtb", dtb)

	err = os.MkdirAll(filepath.Dir(dst), 0o600)
	if err != nil {
		return err
	}

	if err := copy.File(src, dst); err != nil {
		return err
	}

	return nil
}
