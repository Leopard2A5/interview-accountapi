#!/usr/bin/env sh

while true; do
    curl "${BASEURL}/health" &> /dev/null
    if [ "$?" == "0" ]; then
        break
    fi
    echo "waiting for accountapi to come up..."
    sleep 3
done

go test -v .
