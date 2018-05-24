#!/bin/bash
sleep infinity & PID=$!

finish()
{
	kill $PID
	echo "finish_trap=$(date +%s)" >> ./log.${RUN}
}
trap finish SIGTERM


echo "start_trap=$(date +%s)" >> ./log.${RUN}
wait
