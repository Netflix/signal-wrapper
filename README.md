[![CircleCI](https://circleci.com/gh/Netflix/signal-wrapper.svg?style=svg)](https://circleci.com/gh/Netflix/signal-wrapper)
# Signal Wrapper
Signal wrapper is a very simple tool. It wraps programs to help them deal with signals. 

`Usage: ./bin/signal-wrapper-darwin-amd64 [shutdown command] [cmd] [args]`

Let's say you have a program in a container. Let's say that this container is in a load balancing system, and it relies on that container's health to decide whether to route traffic to it. 

For now, we can say the healthcheck is `test ! -f /tmp/unhealthy`, and it runs every 15 seconds. Now, service discovery systems can be slow, so we might want to cause our healthcheck to fail, and wait another little bit to ensure that everyone has caught up with our state.

You might have a script: `shutdown.sh`:

```
#/bin/bash
# Let's make ourselves unhealthy
touch /tmp/unhealthy
# And wait for everyone to catch up.
sleep 60
```

You'd invoke it something like this:
`./signal-wrapper-darwin-amd64 ./shutdown.sh my_real_program -f -p 80`

It will call shutdown.sh, and wait until it exits before sending the signal to the "real" program.

## Known issue
This program only catches SIGTERM, or SIGINT. All other signals are handled as if they were delivered directly to the signal-wrapper. We may want to whitelist a set of signals which are forwarded.
