#/bin/bash

rm meetingSrv
cd ../cmd
go build -o ../bin/meetingSrv

cd ../bin

if [ "$1" = "run" ];then
./meetingSrv --http.savepath=/tmp/meetingBucket/ --wx.state=developer
fi