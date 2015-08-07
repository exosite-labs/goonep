About goonep
============
goonep is a Go binding for Exosite One Platform API. The API is exposed over HTTP in a JSON RPC style interface. 

License is BSD, Copyright 2015, Exosite LLC (see LICENSE file)

========================================
Quick Start
========================================
1.) [Install Go](http://golang.org/doc/install) on your system

2.) Set up your workspace in the [manner described in the Go documentation](https://golang.org/doc/code.html#Workspaces)

3.) Clone this repository to the src/github.com/exosite-labs/goonep directory in your Go workspace

4.) Once you've set up your workspace and set your GOPATH, run the command `go get github.com/stretchr/testify/assert` to install a library goonep uses in testing



Getting A CIK
-------------

Access to the Exosite API requires a Client Identification Key (CIK). 
If you're just getting started with the API and have signed up with a 
community portal, here's how you can find a CIK:

1.) Log in: https://portals.exosite.com

2.) Click on "devices" on the menu on the left

3.) Click on a device to open its properties

4.) The device's CIK is displayed on the left

Once you have a CIK, you can substitute it in the tests. Note that any functions that take a parameter called `auth` can take a string CIK directly.


General API Information
-----------------------

For more information on the API, see:

https://github.com/exosite/docs

HTTP Data Interface
-------------------

The HTTP Data Interface is a minimal HTTP API best suited to resource-constrained 
devices or networks. It is limited to reading and writing data one point at a 
time.

The API is documented here:

http://docs.exosite.com/http/