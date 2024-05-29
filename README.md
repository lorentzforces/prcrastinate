# prcrastinate

A convenience CLI utility to ease dealing with lots of PR review

## Building the Project

### Requirements:

- a Golang installation (built & tested on Go v1.22)
- an internet connection to download dependencies (only necessary if dependencies have changed or
  this is the first build)
- to use the provided `Makefile`, a `make` installation which supports `.RECIPEPREFIX` (_e.g._ GNU
  Make >=4.0)

To build the project, simply run `make` in the project's root directory to build the output
executable.

> _Note: running with `make` is not strictly necessary. Reference the provided `Makefile` for
> typical development commands._
