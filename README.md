# Golang Grafana API

Grafana HTTP API Client for Go

## Installation

    $ go get github.com/AutogrowSystems/go-grafana-api

## Usage

First create a client, you can use an API key or username:pass auth for the first argument, but be aware that
the API key is not supported for some API endpoints.  See the [documentation](http://docs.grafana.org/http_api/admin/)
for more info.

```golang
client, err := gapi.New("username:pass", "http://localhost:3000")
```

Once you have the client, you can perform various operations:

```golang
org, err := client.NewOrg("freds beaver tanks")
fmt.Println(org.Id)

datasources, err := org.DataSources(client)
ds = datasources[0]
ds.IsDefault = true

if err := client.UpdateDataSource(ds); err != nil {
    panic(err)
}
```

See the documentation for other methods.

## CLI Usage

There is also a CLI tool that comes with the package called `gapi`.  It's experimental and will
probably change a lot in the future (aiming towards outputting more JSON).

List all orgs:

    $ gapi -org -list
    1       Main Org.
    2       freds beaver tanks

Add a new org:

    $ gapi -org -create -name "horse monkey"
    2018/03/01 22:58:36 created new org with ID 3
    3

Note that the log output goes to stderr and the org ID goes to stdout.

Add a new datasource:

    $ cat newdatasource.json | gapi -datasource -create

## Tests

To run the tests:

```
go test
```
