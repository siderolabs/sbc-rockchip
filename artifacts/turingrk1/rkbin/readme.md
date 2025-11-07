# DDR Binary Patching for Turing RK1

The RKBIN DDR initialization binary for RK3588 is patched to configure the internal UART
for compatibility with the Turing Pi 2 backplane. Without this patch, the `tpi uart` command
will not function properly — attempting to read the console from the Turing Pi 2's management
interface will fail when accessing the RK1 module.

## Changes Applied

- **UART ID**: Changed from 2 to 9 (internal UART blade port)
- **Baud Rate**: Changed from 1,500,000 to 115,200 (expected by the internal UART blade port)

These modifications aligns the UART configuration with the console kernel parameters
defined in the installer, enabling proper communication through the Turing Pi 2's management system.