
## Getting recent versions of PostgREST on M1/M2 Macs

Newt is intended to work along side Postgres and PostgREST. I usually install PostgREST from source.  I've had the best results on both Linux and macOS using GHCup to provide my Haskell build environment.  Below is a quick recipe for building PostgREST. It is not intended to replace the PostgREST project documentation or installation instructions. It is offered here for convenience. This recipe reflects what I did to get PostgREST current releases installed on a Mac Mini with a M1 processor running macOS Sonoma. Your mileage may vary considerably.

The basic recipe is

1. Install Haskell with [GHCup](https://www.haskell.org/ghcup/) (I accept the defaults and use the "recommended" versions set via `ghcup tui`)
2. Make sure the GHCup environment is available, `source $HOME/.ghcup/env`
3. Make sure I am using the "recommended" of GHC, Cabal, Stack, etc.
4. Clone the GitHub [PostgREST](https://github.com/PostgREST/postgrest)
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

