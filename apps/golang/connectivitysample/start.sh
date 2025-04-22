#!/bin/bash -e
update-ca-certificates
./connectivitysample "$@" &
pid=$!
trap "kill -s INT ${pid}; wait -n ${pid}; exit $?" SIGINT
trap "kill -s TERM ${pid}; wait -n ${pid}; exit $?" SIGTERM
trap "kill -s KILL ${pid}; wait -n ${pid}; exit $?" SIGKILL
wait ${pid}
