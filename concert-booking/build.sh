# /bin/bash

CURRENT_DIR="$PWD"

# build userservice
cd userservice
go build -v -o user

cd $CURRENT_DIR

# build bookingservice
cd bookingservice
go build -v -o booking

cd $CURRENT_DIR

# build concertservice
cd concertservice
go build -v -o concert
