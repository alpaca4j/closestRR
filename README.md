# closestRR

## Name

*closestRR* - if there are multiple A records in a response, the plugin will return only those which fit inside the subnet calculated from requestors IP/23.
If non of them fit, it returns all records.



## Description

It uses /23 to calculate the network addres for the requestor IP and checks if there is an IP in the DNS anwser which fits in that subnet.
 
 

## Compilation

This package will always be compiled as part of CoreDNS and not in a standalone way. It will require you to use `go get` or as a dependency on [plugin.cfg](https://github.com/coredns/coredns/blob/master/plugin.cfg).

The [manual](https://coredns.io/manual/toc/#what-is-coredns) will have more information about how to configure and extend the server with external plugins.

A simple way to consume this plugin, is by adding the following on [plugin.cfg](https://github.com/coredns/coredns/blob/master/plugin.cfg), and recompile it as [detailed on coredns.io](https://coredns.io/2017/07/25/compile-time-enabling-or-disabling-plugins/#build-with-compile-time-configuration-file).

~~~
closestRR:github.com/alpaca4j/closestAaddress/
~~~

After this you can compile coredns by:

```shell script
go generate
go build
```

Or you can instead use make:

```shell script
make
```

## Syntax

~~~ txt
closestRR [ZONES...]
~~~

## Ready

This plugin reports readiness to the ready plugin. It will be immediately ready.

## Examples


~~~ corefile
. {
  forward . 9.9.9.9
  closestRR example.org example2.org
}
~~~

## Also See

See the [manual](https://coredns.io/manual).
