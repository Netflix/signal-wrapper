#!/bin/bash
export RUN=${RANDOM}
echo "Running test: ${RUN}"
rm -f log.${RUN} || true
echo "Deleting log"

echo "Starting signal wrapper"
echo "start_test=$(date +%s)" >> log.${RUN}
./signal-wrapper test/shutdown.sh test/trap.sh & PID=$!
echo "Started signal wrapper with PID: ${PID}"

sleep 5 # This is to wait until everything is scaffolded
kill -SIGTERM $PID
wait $PID

echo "Verifying test"
. log.${RUN}

export start_test
export start_trap
export start_shutdown
export finish_shutdown
export finish_trap

./verify.py || exit 1

rm log.${RUN}
