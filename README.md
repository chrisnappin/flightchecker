# Flight Checker

Using Go 1.13 (via brew)
Created project by running: go mod init gitlab.com/chrisnappin/flightchecker

todo - makefile and package structure

## Data Sources
Provides airports with IATA codes, downloaded from https://ourairports.com/data/ (free, public domain)

Need to save as:
* `data/airports/airports.csv` (8,374,699 bytes, last modified Oct 6, 2019)
Large file, containing information on all airports on this site.

* `data/airports/countries.csv` (20,386 bytes, last modified Oct 6, 2019)
A list of the world's countries. You need this spreadsheet to interpret the country codes in the airports and navaids files.

* `data/airports/regions.csv` (365,815 bytes, last modified Oct 6, 2019)
A list of all countries' regions (provinces, states, etc.). You need this spreadsheet to interpret the region codes in the airport file.


## How to build
* Check the repo out to anywhere outside of $GOROOT
* Set your API Host and Key in `arguments.json`
* `go build ./...`
* `~/go/bin/flightchecker`


Intention of the Go tool is not to need Makefiles!

Run `go task ./...` where `task` can be
* clean => removes object files from source dirs, mostly un-needed!
* list => lists all packages and imported modules
* get => downloads and install code or modules (to the $GOPATH)
  * e.g. `go get github.com/sirupsen/logrus` (uses latest version) - updates `go.mod`
* build => compiles all code, but throws away the results
* test => compiles and runs all tests, caches results from unchanged tests
* install => compiles all code, writes them to `~/go/bin`
