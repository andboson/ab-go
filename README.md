# ab-go

run ./ab-go for usage

## run options:

#### -test
Testing mode. Run cycled tests with send result to graphics. R
Run with this flag and open `localhost:9999` (default port) in your browser
<img src="tests/screenshot.png" />
Data is updated every second.

#### -port
Embedded web-serser port. Used with `-test` flag


-----------------
#### build:
mkdir $GOPATH/src/ab-go && cd $GOPATH/src/ab-go
git clone git@github.com:andboson/ab-go.git .
make
