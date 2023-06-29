
Installation
------------

newt is an experimental web service for working with Pandoc+PostgREST and Pandoc. It also provides a simple static file web service. At this time you must install it from source code. It is probably broken and certainly will contain many bugs. It's a prototype!

Quick install with curl
-----------------------

If you are running macOS or Linux you can install released versions of newt
with the following curl command.

~~~
curl https://caltechlibrary.github.io/newt/installer.sh
~~~

Install from source
-------------------

## Requirements

- Golang > 1.20
- Pandoc > 3
- Postgres > 15
- PostgREST > 11
- GNU Make
- Git

## Steps

1. Clone the Git repository for the project
2. change directory into the cloned project
3. Run `make`, `make test` and `make install`

Here's what that looks like for me.

~~~
git clone https://github.com/caltechlibrary/newt src/github.com/caltechlibrary/newt
cd src/github.com/caltechlibrary/newt
make
make test
make install
~~~

By default it will install the programs in `$HOME/bin`. `$HOME/bin` needs
to be included in your `PATH`. E.g.

~~~
export PATH="$HOME/bin:$PATH"
~~~

Can be added to your `.profile`, `.bashrc` or `.zshrc` file depending on your system's shell.

## Precompiled binaries

You can find pre-compiled binaries for some versions of Newt. They are available at <https://github.com/caltechlibrary/newt/releases>.

The binaries are available in a Zip archive file for download. The name uses the form `newt-<VERSION_NO>-<OS_TYPE>-<CPU_TYPE>.zip`. VERSION_NO will be a semver associated with the release (e.g. "v0.0.2"), the OS_TYPE will be either "Linux" (including Raspbery Pi OS), "Windows" and "macOS". The CPU_TYPE will vary based on how the CPU type is reported on the POSIX system (e.g. what is returned by `uname -m`) or in the case of Windows either "x86_64" for Intel CPU types or "arm64" for ARM CPU (e.g. those in selected Surface tablets or Microsoft's ARM Development Kit 2023).

- macOS example
    - `newt-v0.0.2-macOS-arm64.zip` (M1, M2 CPU) or `newt-v0.0.2-macOS-x86_64.zip` (older Intel based Macs)
- Windows example
    - `newt-v0.0.2-Windows-x86_64.zip` (for Most Windows machines)
- Linux example (including Raspberry Pi OS)
    - `newt-v0.0.2-Linux-x86_64.zip` (Linux on Intel CPUs)
    - `newt-v0.0.2-Linux-armv7l.zip` (Raspberry Pi OS, 32bit)
    - `newt-v0.0.2-Linux-aarch64.zip` (Linux on ARM 64)

