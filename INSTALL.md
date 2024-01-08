
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

- Golang > 1.21
- Pandoc > 3
- Postgres > 15
- PostgREST > 12
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

## Getting recent versions of Pandoc and PostgREST on M1/M2 Macs

Newt is intended to work along side Pandoc and PostgREST. I usually install these from source.  I've had the best results on both Linux and macOS using GHCup to provide my Haskell compile and build environment.  Below are quick recipes for building and installing PostgREST and Pandoc, these are not a replacement for their respective project documentation instruction. They reflect what I did to get Pandoc/PostgREST current releases installed on a Mac Mini with a M1 processor. Your mileage may vary.

The basic recipe is

1. Install Haskell with [GHCup](https://www.haskell.org/ghcup/) (I accept the defaults and use the "recommended" versions set via `ghcup tui`)
2. Make sure the GHCup environment is available, `source $HOME/.ghcup/env`
3. Make sure I am using the "recommended" of GHC, Cabal, Stack, etc.
4. Clone the GitHub [PostgREST](https://github.com/PostgREST/postgrest) and  [Pandoc](https://github.com/jgm/pandoc) repositories to your machine
5. Change to the cloned repository directory
6. Checkout the version you want to build
7. Run the usual Haskell/Cabal build process

Here is the steps I typed at the command line to install PostgREST on my M1 Mac Mini.

~~~
curl --proto '=https' --tlsv1.2 -sSf https://get-ghcup.haskell.org | sh
source $HOME/.ghcup/env
ghcup tui
git clone git@github.com:PostgREST/postgrest
cd postgrest
git checkout v11.1.0
cabal clean
cabal update
cabal build
cabal install
cd ..
~~~

Here are the steps I typed at the command line to install Pandoc on my M1 Mac Mini. The basic recipe is modified because of a potential link conflict in which libiconv to use I've encountered on macOS running on my M1 Mac Mini.  I need the Mac Ports version for the Mac Ports installed of Git to work. This means I need to do all my Git commands before I removing libiconv. I then invoking Cabal and put libiconv back afterward Pandoc installation is complete. If you're not using the libiconv installed as part of Mac Ports applications you can skip those steps.

NOTE: I've skipped installing GHCup because I assume you've already installed it when you compiled.

1. Clone [Pandoc](https://github.com/jgm/pandoc) repository from GitHub
2. Change into the Pandoc directory
3. Checkout the version of Pandoc you want to build (e.g. 3.1.4)
4. Remove the Mac Ports libiconv library using Ports command
5. Run the usual Haskell/Cabal build process
6. Put the Mac Ports libiconv back using the Ports command


~~~
git clone git@github.com:jgm/pandoc
cd pandoc
git checkout 3.1.4
sudo port uninstall libiconv
cabal clean
cabal update
cabal install pandoc-cli
sudo port install libiconv
cd ..
~~~

Even though you are installing “pandoc-cli” it can function as the Pandoc web service by invoking pandoc server command option.
