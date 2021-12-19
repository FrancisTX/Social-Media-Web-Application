# Distributed Finaly Project Course Materials

## Usage
```
git clone https://github.com/os3224/final-project-0b5a2e16-FrancisTX-hannnnk1231.git
cd final-project-0b5a2e16-FrancisTX-hannnnk1231
goreman start
```
## Test
### Web
```
go test -v web_test.go web.go -cover
```

### Client
```
cd client
go test -v -cover
```
### Server
```
cd server
go test -v -cover
```
### Storage
```
cd server/storage
go test -v -cover
```

## Hashicorp raft demo
Although we use etcd/raft in our main branch, we till try to implement the simple key-value storage in hashicorp/raft in the branch ``` hashicorp_raft```.
### Usage
```
cd hashicorp_demo
./start.sh 1
```
```
./start.sh 2
```
```
./start.sh 3
```
Now, we have already opened the three raft node in a cluster.

### Sample request
For set the "key", "val" into the storage:
```
curl -X POST http://127.0.0.1:8000/key -d '{"key": "val"}' -L
```

For get the "key" from the storage:
```
curl http://127.0.0.1:8001/key/"key" -L
```

For delete the "key" from the storage:
```
curl -X DELETE http://127.0.0.1:8000/key/"key" -L
```

## Protobuf Generator

Put this shell script in the root of every project I make that generates code from protobuf files.

It scans the caller's directory for folders ending in `pb`, and generates the protobufs therein.

It's only been tested on mac, but it should work on any nixy system.  The only non-standard tool it uses is `tree`, and that's nonessential.  If you want it anyway, ond you're on n a mac, do `brew install tree`.

Obviously, this script also requires protobufs, grpc, and the go-bindings thereof.  For installation instructions there, I'll defer to [this page](https://grpc.io/docs/quickstart/go.html).

If you're able to drop this gen proto script into the root of the example directory they provide, rename `helloworld` to `helloworldpb`, and the script executes without complain, then you're probably good to use this everywhere!
