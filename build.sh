mkdir -p ~/.pppsrv/logs
mkdir -p ~/.pppsrv/sockets
rm ~/.pppsrv/game
cd game
make
cd ..
mv ./game/game ~/.pppsrv/game
go get ./server/...
go install ./server/...
