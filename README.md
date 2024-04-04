DNSctl
=======

Tool made to ease and simplify DNS records updates from the cli.

TODO:

* Add records.
* Copy between providers.
  * Exclude non-standard records (Alias, etc)
* Delete Records.
* Export to json.
* Export to Terraform.
* Add support for different providers.


Compile
-------


`make`
Will compile dnsctl and the plugins without tests nor cleanup.

`make test`
Runs unit tests.

`make all`
Will compile dnsctl and the plugins with cleanup and tests

`make clean`
cleans up build and bin directories.

`make plugin`
Will compile plugins folder only.

`make install`
Will run `make all` and then will put the dnsctl binary at `/usr/local/bin` and the plugins at `/usr/local/lib/dnsctl` if root is used; for local installs the setup will use `$HOME/.conf/dnsctl` and `$GOPATH/bin` will be used.

Usage
-----
Each plugin defines the criteria for different DNS providers. Current development is for CloudFlare, GCP and AWS

```
dnsctl add donrudo.com provider   
dnsctl add bleh.donrudo.com cname endpoint.com
dnsctl record delete bleh.donrudo.com
```
