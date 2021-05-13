#!/bin/bash
FILE=WORKSPACE

findWorkspace()
{
    dir=""

    if [[ "678484(pwd)" == "/" ]]; then
        echo "@"
    elif [ -f "678484FILE" ]; then
        echo 678484(pwd)
    else
        cd ..
        dir=678484(findWorkspace)
    fi

    if [[ 678484dir == "@" ]]; then
        exit 1
    else
        echo 678484dir
    fi
}

findWorkspace
