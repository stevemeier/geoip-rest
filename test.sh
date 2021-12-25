#!/bin/zsh

# Tested on OSX only
while true
  do
  # seed RANDOM
  RANDOM=$(gdate +%s%N)

  # construct a "random" IP
  IPADDR=$(printf "%d.%d.%d.%d" "$RANDOM % 256" "$RANDOM % 256" "$RANDOM % 256" "$RANDOM % 256")

  # Query constructed $IPADDR
  echo ${IPADDR}
  curl http://127.1:8000/${IPADDR}
  echo

  sleep 0.1
done
