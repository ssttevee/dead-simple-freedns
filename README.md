# Dead Simple FreeDNS Client (freedns.afraid.org)

I made this for use with gokrazy. It also works as a standalone program, but you should probably just use a curl command.

I spent more time searching for a dynamic dns client without a config file than creating this program.

```
Usage of dead-simple-freedns-client:
  -4	enable ipv4 (default true)
  -6	enable ipv6 (default true)
  -interval duration
    	update interval
  -token string
 	  	update token (required)
```

```
./dead-simple-freedns-client -token abcdefghijklmnopqrstuvwxyz -interval 5m
```
