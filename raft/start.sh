node=$1

if [ "${node}" == 1 ] ; then
    args="--bootstrap true"
else
    args="--join=127.0.0.1:8001"
fi

raft_port=$((7000 + $node))
args="${args} --raft-port=${raft_port}"

http_port=$((8000 + $node))
args="${args} --http-port=${http_port}"

dir="./store${node}"
args="${args} --store-dir=${dir}"

go build raftserver.go
echo $args
./raftserver ${args}

