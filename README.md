# Flight Checker

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
* Install sqlite3 locally
  * e.g. on MacOSX `brew install sqlite3`
* Install gcc locally
  * e.g. on MacOSX install xcode
* Build the sqlite 3 driver using GCC
  * `CGO_ENABLED=1; go install --tags "darwin" github.com/mattn/go-sqlite3`   
* Set your API Host and Key in `arguments.json`
* Build the code, doesn't need GCC
  * `go build ./...`
* Run the flight checker
  * `~/go/bin/flightchecker`
* Or run the airport code finder
  * `~/go/bin/airports`


Intention of the Go tool is not to need Makefiles!

Run `go task ./...` where `task` can be
* clean => removes object files from source dirs, mostly un-needed!
* list => lists all packages and imported modules
* get => downloads and install code or modules (to the $GOPATH)
  * e.g. `go get github.com/sirupsen/logrus` (uses latest version) - updates `go.mod`
* build => compiles all code, but throws away the results
* test => compiles and runs all tests, caches results from unchanged tests (add `-cover` for coverage stats)
* install => compiles all code, writes them to `~/go/bin`
