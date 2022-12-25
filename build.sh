#!/bin/bash

usage()
{
    echo "===========================HELP============================"
    echo "build.sh [-b build configuration (full/stripped)]"
    exit 1
}

while getopts 'b:h?' opt; do
    case "$opt" in
        b)
            echo "running application go build (${OPTARG})"

            cd src

            if [[ $OPTARG = "stripped" ]]
            then
               go build -o ../build/im-next -ldflags '-s'
            fi

            if [[ $OPTARG = "full" ]]
            then
                go build -o ../build/im-next
            fi

            echo "done building."
            echo
            exit 0
            ;;
        *)
            usage
            ;;
    esac
done
shift "$(($OPTIND -1))"
usage