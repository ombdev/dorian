#!/bin/bash


prog="Redis 3 Server"
REDIS_EXE="redis-server"
REDIS_CONF="$HOME/redis/redis.conf"
REDIS_SERVER="$HOME/redis/src/$REDIS_EXE $REDIS_CONF"

start() {

    echo "Starting $prog"
    $REDIS_SERVER

}

stop() {

    echo "Stopping $prog"
    pkill $REDIS_EXE

}

restart() {

    echo "Restarting $prog"
    pkill $REDIS_EXE

    echo "Starting $prog"
    nohup $REDIS_SERVER &
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
