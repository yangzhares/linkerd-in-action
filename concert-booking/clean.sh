#/bin/bash 

if [ -f userservice/user ] 
then
    rm -f userservice/user
fi 

if [ -f bookingservice/booking ] 
then
    rm -f bookingservice/booking
fi 

if [ -f concertservice/concert ] 
then
    rm -f concertservice/concert
fi 