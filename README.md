About goonep
============

Getting A CIK
-------------

Access to the Exosite API requires a Client Identification Key (CIK). 
If you're just getting started with the API and have signed up with a 
community portal, here's how you can find a CIK:

1.) Log in: https://portals.exosite.com

2.) Click on "devices" on the menu on the left

3.) Click on a device to open its properties

4.) The device's CIK is displayed on the left

Once you have a CIK, you can substitute it in the examples below. Note that any functions that take a parameter called `auth` can take a string CIK directly, or you can pass an auth dictionary as described [here](http://docs.exosite.com/rpc/#authentication).


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

https://github.com/exosite/docs/tree/master/data