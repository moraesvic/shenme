#!/usr/bin/env bash

set -ex

go build .

if [ ! -d ~/bin/ ] ; then
    mkdir -p ~/bin/
fi

cp ./shenme ~/bin/shenme
