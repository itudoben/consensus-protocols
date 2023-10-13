module itudoben.io/node

go 1.20

replace itudoben.io/state => ./state

require (
	github.com/brotherpowers/ipsubnet v0.0.0-20170914094241-30bc98f0a5b1
	gopkg.in/netaddr.v1 v1.5.1
	itudoben.io/state v0.0.0-00010101000000-000000000000
	v.io/x/lib v0.1.15
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
