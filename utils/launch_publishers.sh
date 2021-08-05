#!/bin/sh

BROKER="localhost:1883"
MSG_RATE="1000"
DURATION="0"
REQUEST_SIZE="1024"
MAX_PUBLISHERS="10"
SPAWN_DELAY="120"
QOS="1"
DOCKER_IMAGE_NAME="claudiolor/latency-mqtt-client"
NETWORK_NAME="publishers-net"

usage(){
	echo "Usage: $PROG_NAME [-options]"
	echo "  -b, --broker         Specify broker address [DEFAULT=$BROKER]"
	echo "  -r, --message-rate   Specify the message rate  [DEFAULT=$MSG_RATE]"
	echo "  -d, --duration       Specify the duration time of the experiment is seconds (0 is unlimited) [DEFAULT=$DURATION]"
	echo "  -s, --size           Specify the size of the messages [DEFAULT=$REQUEST_SIZE]"
	echo "  -p, --max-publishers Specify the maximum number of publishers [DEFAULT=$MAX_PUBLISHERS]"
	echo "  -w, --delay          Specify the delay after which another publisher is spawned in seconds [DEFAULT=$SPAWN_DELAY]"
	echo "  -q, --qos            Specify the qos between [1,2,3] [DEFAULT=$QOS]"
	exit 1
}

parse_args(){
	while [ "${1:-}" != "" ]; do
		case "$1" in
			"-b" | "--broker")
				shift
				BROKER=$1
				;;
			"-r" | "--message-rate")
				shift
				MSG_RATE=$1
				;;
			"-d" | "--duration")
				shift
				DURATION=$1
				;;
			"-s" | "--size")
				shift
				REQUEST_SIZE=$1
				;;
			"-p" | "--max-publishers")
				shift
				MAX_PUBLISHERS=$1
				;;
			"-w" | "--delay")
				shift
				SPAWN_DELAY=$1
				;;
			"-q" | "--qos")
				shift
				QOS=$1
				;;
			*)
				usage
				;;
		esac
		shift
	done
}

cleanup(){
	echo "Cleaning up environment"
	PUBLISHERS="$(docker container ps --filter "name=mqtt-delay-publisher" -q)"

	if [ "$PUBLISHERS" = " " ];
	then
		docker container stop $PUBLISHERS
	fi
	docker network rm $NETWORK_NAME
	exit 0
}

trap "cleanup" 2

parse_args "$@"

# Create a docker network
docker network create $NETWORK_NAME

# Start launching publishers
echo "...LAUNCHING PUBLISHERS..."
i=1
while [ $i -le $MAX_PUBLISHERS ]
do
	 docker run --rm --network $NETWORK_NAME --name mqtt-delay-publishers-$i $DOCKER_IMAGE_NAME --broker="$BROKER" --messageRate="$MSG_RATE" --duration="$DURATION" --requestSize="$REQUEST_SIZE" --qos="$QOS" &
	 i=$((i+1))
	 sleep $SPAWN_DELAY
done

echo "Press a key to exit..."
read _
cleanup