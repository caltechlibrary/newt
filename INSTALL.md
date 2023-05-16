
Installation
------------

Quick install with curl
-----------------------

The following shell command should work on most POSIX systems where
curl is installed.

~~~
curl https://caltechlibrary.github.io/newt/installer.sh | sh
~~~

Install from source
-------------------

## Requirements

- Golang >= 1.20
- Pandoc >= 3
- Postgres >= 15
- PostgREST >= 11
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


