mkdir ipc
OUT_PKG="ipc/"
PROTO_FILE="*.proto"
protoc -I=. --go_out=${OUT_PKG} ${PROTO_FILE}
rm ../publisher/ipc/*.pb.go
rm ../subscriber/ipc/*pb.go

cp ./ipc/*.pb.go ../publisher/ipc/ 
cp ./ipc/*.pb.go ../subscriber/ipc/
rm -rf ipc