
Installation
------------

Newt is an experimental set of programs for rapid application development for libraries, archives and museums. It provides a data router, code generator and Mustache template engine. At this time you must install it from source code. The programs are likely broken, poorly documented and certainly contain bugs. It's a prototype!

Quick install with curl or irm
------------------------------

If you are running macOS or Linux you can install released versions of newt
with the following curl command.

~~~
curl https://caltechlibrary.github.io/newt/installer.sh
~~~

In Powershell on Windows you can use the following command

~~~
irm https://caltechlibrary.github.io/newt/installer.ps1 | iex
~~~


Install from source
-------------------

## Installation Requirements

- Go >= 1.22.5
- Deno >= 1.44
- Pandoc > 3.1 (to render Markdown content and some Go code)
- Git compatible tool for retrieving the GitHub repository of Newt
- GNU Make > 4
- GNU grep > 3
- GNU cut  > 8

### Building from source code

1. Make sure you have the required development software installed
2. Clone the Newt GitHub repository locally
3. Use the Go tool to build, test and run the application

##### step 1

Make sure you have the necessary software installed first. The list is below.  If you are building on Windows you'll need to use the LSW to follow along and the tools will need to be available in the WSL. Otherwise the steps are generally similar for macOS, Windows and Linux.

- Go >= 1.22.5
- Deno >= 1.44
- Pandoc > 3.1 (to render Markdown content and some Go code)
- Git cli or compatible tool for retrieving the GitHub repository of Newt
- GNU Make > 4
- GNU grep > 3
- GNU cut  > 8

#### step 2

1. Start in your home directory
2. Clone Newt's [repository](https://github.com/caltechlibrary/newt).
3. Change into the Newt directory

From the terminal application running a POSIX shell I type the following.

~~~shell
cd
git clone \
   https://github.com/caltechlibrary/newt \
   src/github.com/caltechlibrary/newt
cd src/github.com/caltechlibrary/newt
~~~

#### step 3

For step three tasks you have two choices. You can use the Go command Or you can run the Make command. The latter has the advantage that it'll build help docuemntaiton, man pages and install the man pages in your home directory.

It assumed you are still in the directory where you completed step 2.

1. Build the newt application
2. Test the newt application
3. If the tests are successful then you can install Newt tools

Here's the Go command version.

~~~shell
go build
go test
go install
~~~

Here's the Make version.

~~~
make
make test
make install
~~~

<!--

## Precompiled binaries

You can find prep-compiled binaries for some versions of Newt. They are available at <https://github.com/caltechlibrary/newt/releases>.

The binaries are available in a Zip archive file for download. The name uses the form `newt-<VERSION_NO>-<OS_TYPE>-<CPU_TYPE>.zip`. VERSION_NO will be a semver associated with the release (e.g. "v0.0.2"), the OS_TYPE will be either "Linux" (including Raspberry Pi OS), "Windows" and "macOS". The CPU_TYPE will vary based on how the CPU type is reported on the POSIX system (e.g. what is returned by `uname -m`) or in the case of Windows either "x86_64" for Intel CPU types or "arm64" for ARM CPU (e.g. those in selected Surface tablets or Microsoft's ARM Development Kit 2023).

- macOS example
  - `newt-v0.0.2-macOS-arm64.zip` (M1, M2 CPU)
  - `newt-v0.0.2-macOS-x86_64.zip` (older Intel based Macs)
- Windows example
  - `newt-v0.0.2-Windows-x86_64.zip` (for Most Windows machines)
- Linux example (including Raspberry Pi OS)
  - `newt-v0.0.2-Linux-x86_64.zip` (Linux on Intel CPUs)
  - `newt-v0.0.2-Linux-aarch64.zip` (Linux on ARM 64)

-->

## Getting help

**The Newt Project is an experiment!!**. The source code for the project is supplied "as is". Newt is most likely broken. At a stretch it could be considered a working prototype. You should not use it for production systems.  However if you'd like to ask a question or have something you'd like to contribute please feel free to file a GitHub issue, see <https://github.com/caltechlibrary/newt/issues>. Just keep in mind it remains an **experiment** as of February 2024.
