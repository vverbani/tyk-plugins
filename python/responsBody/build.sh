#!/bin/bash

# the working directory mounted into the container with -v
WORKDIR=/plugin

# my user id and group id so that we can set the ownership
myUID=1000
myGID=1000

# Bundle name to build
bundle=responsBody

cd $WORKDIR
/opt/tyk-gateway/tyk bundle build -m $WORKDIR/$bundle.json -y -o $WORKDIR/$bundle.zip

# set the ownership of the output file or it will be owned by root
chmod 664 $WORKDIR/$bundle.zip
chown $myUID:$myGID $WORKDIR/$bundle.zip
