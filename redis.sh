#!/bin/bash
#
#


prog="Redis 4.0.8 Server"
REDIS_EXE="redis-server"
REDIS_CONF="$HOME/redis-4.0.8/redis.conf"
REDIS_SERVER="$HOME/redis-4.0.8/src/$REDIS_EXE $REDIS_CONF"

start() {

    echo "Starting $prog"
    $REDIS_SERVER &

}

stop() {

    echo "Stopping $prog"
    pkill $REDIS_EXE

}

restart() {

    echo "Stopping $prog"
    pkill $REDIS_EXE

    echo "Starting $prog"
    $REDIS_SERVER &
}

case "$1" in
  start)
        start
        ;;
  stop)
        stop
        ;;
  restart)
        restart
        ;;
  *)
        echo $"Usage: $0 {start|stop|restart}"
        exit 2
esac

exit $?
