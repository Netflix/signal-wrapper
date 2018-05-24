#!/bin/bash
echo "start_shutdown=$(date +%s)" >> ./log.${RUN}
sleep 9
echo "finish_shutdown=$(date +%s)" >> ./log.${RUN}

