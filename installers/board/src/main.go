package main

import (
	_ "embed"
	"path/filepath"

	"github.com/siderolabs/go-copy/copy"
	"github.com/siderolabs/talos/pkg/machinery/overlay"
	"github.com/siderolabs/talos/pkg/machinery/overlay/adapter"
)

func main() {
	adapter.Execute(&BoardInstaller{})
}

type BoardInstaller struct{}

type boardExtraOptions struct {
	Console    []string `json:"console"`
	ConfigFile string   `json:"configFile"`
}

func (i *BoardInstaller) GetOptions(extra boardExtraOptions) (overlay.Options, error) {
	kernelArgs := []string{
		"console=tty0",
		"sysctl.kernel.kexec_load_disabled=1",
		"talos.dashboard.disabled=1",
	}

	kernelArgs = append(kernelArgs, extra.Console...)

	return overlay.Options{
		Name:       "board",
		KernelArgs: kernelArgs,
	}, nil
}

func (i *BoardInstaller) Install(options overlay.InstallOptions[boardExtraOptions]) error {
	// allows to copy a directory from the overlay to the target
	// err := copy.Dir(filepath.Join(options.ArtifactsPath, "arm64/firmware/boot"), filepath.Join(options.MountPrefix, "/boot/EFI"))
	// if err != nil {
	// 	return err
	// }

	// allows to copy a file from the overlay to the target
	err := copy.File(filepath.Join(options.ArtifactsPath, "arm64/u-boot/board/u-boot.bin"), filepath.Join(options.MountPrefix, "/boot/EFI/u-boot.bin"))
	if err != nil {
		return err
	}

	if options.ExtraOptions.ConfigFile != "" {
		// do something with the config file
	}

	return nil
}
