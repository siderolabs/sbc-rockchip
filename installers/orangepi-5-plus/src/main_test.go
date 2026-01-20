// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"reflect"
	"testing"
)

func TestGetOptions_WithoutExtraKernelArgs(t *testing.T) {
	installer := &opi5PlusInstaller{}
	extra := opi5PlusExtraOptions{}

	opts, err := installer.GetOptions(extra)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedArgs := []string{
		"console=tty1",
		"console=ttyS2,1500000",
		"sysctl.kernel.kexec_load_disabled=1",
		"talos.dashboard.disabled=1",
	}

	if !reflect.DeepEqual(opts.KernelArgs, expectedArgs) {
		t.Errorf("expected kernel args %v, got %v", expectedArgs, opts.KernelArgs)
	}
}

func TestGetOptions_WithExtraKernelArgs(t *testing.T) {
	installer := &opi5PlusInstaller{}
	extra := opi5PlusExtraOptions{
		ExtraKernelArgs: []string{
			"cpufreq.default_governor=performance",
			"mitigations=off",
		},
	}

	opts, err := installer.GetOptions(extra)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedArgs := []string{
		"console=tty1",
		"console=ttyS2,1500000",
		"sysctl.kernel.kexec_load_disabled=1",
		"talos.dashboard.disabled=1",
		"cpufreq.default_governor=performance",
		"mitigations=off",
	}

	if !reflect.DeepEqual(opts.KernelArgs, expectedArgs) {
		t.Errorf("expected kernel args %v, got %v", expectedArgs, opts.KernelArgs)
	}
}

func TestGetOptions_ExtraKernelArgsAppendedAtEnd(t *testing.T) {
	installer := &opi5PlusInstaller{}
	extra := opi5PlusExtraOptions{
		ExtraKernelArgs: []string{"custom.arg=value"},
	}

	opts, err := installer.GetOptions(extra)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify that custom arg is at the end
	lastArg := opts.KernelArgs[len(opts.KernelArgs)-1]
	if lastArg != "custom.arg=value" {
		t.Errorf("expected last arg to be 'custom.arg=value', got '%s'", lastArg)
	}

	// Verify hardcoded args are still at the beginning
	if opts.KernelArgs[0] != "console=tty1" {
		t.Errorf("expected first arg to be 'console=tty1', got '%s'", opts.KernelArgs[0])
	}
}
