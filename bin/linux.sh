#/bin/bash

rm meetingSrv
cd ../cmd
GOOS=linux GOARCH=amd64 go build -o ../bin/meetingSrv
cd ../bin
scp meetingSrv dev:~/
ssh dev "./meeting/launch.sh"
