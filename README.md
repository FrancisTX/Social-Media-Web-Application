## Hashicorp raft demo
Although we use etcd/raft in our main branch, we still try to implement the simple key-value storage in this branch.
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
Now, we have already opened the three raft nodes in a cluster.

### Sample request
(Note that any write operation can be only done in leader node, and the read operation can be applied for any nodes.)
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
