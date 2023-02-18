# reKinD 🚀 (a better KinD)

## TL/DR
- this is my first Golang project
- it augments kind, so you are going to have kind installed
- this project is just getting started, don't be too angry yet.

## Why another KinD?

Well I love KinD. It is the best way to do on-machine kubernetes work without much hassle, and no setup.
However, there are some things that should be part of the official KinD, yet they aren't.

This project uses Commandline Augmentation to transparently modify the commandline to fit the regular kind version.

In short: 
1. you call rekind
2. rekind translates its own flags and commands to "regular" commands understood by kind natively.
3. For example: it knows about Kind Nodeimages and picks the right image for the version you asked it for.
4. It adds non-kind commands that involce the use of docker-cli and crictl underneath.

## Why not create a ticket or PR on KinD?
Well i think we might even do this as part of rekind.

One of the drivers for reKinD (other than learning golang) is to make an awesome feature request by showing off cool new flags and commands, showing acceptance by users.

In the longer run, every feature that was introduced by rekind and goes into "normal" kind would be a win.



## What it can do already

Not much, right now yet.
There is:

### Automatic Node Image selection
> rekind create cluster --version=1.24 

which will translate to:

>  rekind create cluster --image=kindest/node:v1.24.7@sha256:577c630ce8e509131eab1aea12c022190978dd2f745aac5eb1fe65c0807eb315

when run with kind version 0.17.

## Roadmap

### New Commands
I intend to add support for the following new commands:

### Effortless interaction with crictl to list and sync images to nodes
> rekind get images

> rekind sync images

### Pick Kubernetes version from another context 

> rekind create cluster --version-from=my-other-cluster

This is quite useful if you want to stay in sync with the big cluster version you are running in the cloud.

### Automatically create namespaces

> rekind create cluster --withNamespaces=foo,bar

Helps setting up you cluster with fewer commands.

Even better, copy namespaces from another context:

> rekind create cluster --withNamespaces-from=my-other-cluster

## Further Ideas

Of course, we can do exactly this to other CLIs.
How about a "rejava" that allows to run java programs with other, 
more opinionated defaults? 

## Collaboration
Contact me on https://www.linkedin.com/in/tonit/ or here on github.

