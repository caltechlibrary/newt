
# Installing Pandoc from source

I've had the best results compiling Pandoc on macOS using GHCup to provide my Haskell build environment.  This doesn't replace the recommendations at the [Pandoc Website](https://pandoc.org). It's the my notes about how I got things up on my Mac Mini with an M1 process while using Mac Ports. The basic recipe is modified because of a potential link conflict in which libiconv if you've also installed Mac Ports to get reliable development tools like Git. Pandoc needs to link against the macOS supplied libiconv but Mac Ports also provides this and the Git that comes with Mac Ports relies on the Mac Ports libiconv. The trick then is to do all the Git commands before 
before temporarily uninstalling the libiconv that comes with Mac Ports. You can then invoking Cabal and put libiconv back afterward Pandoc installation is complete. If you're not using the libiconv installed as part of Mac Ports applications you can skip those steps.

The basic recipe

1. Install Haskell with [GHCup](https://www.haskell.org/ghcup/) (I accept the defaults and use the "recommended" versions set via `ghcup tui`)
2. Make sure the GHCup environment is available, `source $HOME/.ghcup/env`
3. Make sure you are using the "recommended" of GHC, Cabal, Stack, etc.
4. Clone [Pandoc](https://github.com/jgm/pandoc) repository from GitHub
5. Change into the Pandoc directory
6. Checkout the version of Pandoc you want to build (e.g. 3.1.4)
7. Remove the Mac Ports libiconv library using Ports command
8. Run the usual Haskell/Cabal build process
9. Put the Mac Ports libiconv back using the Ports command


~~~
curl --proto '=https' --tlsv1.2 -sSf https://get-ghcup.haskell.org | sh
source $HOME/.ghcup/env
ghcup tui
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

NOTE: Even though you are using cabal with the "install pandoc-cli" option the resulting Pandoc can function as the Pandoc web service by invoking Pandoc server command option.

