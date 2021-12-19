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


## Protobuf Generator

Put this shell script in the root of every project I make that generates code from protobuf files.

It scans the caller's directory for folders ending in `pb`, and generates the protobufs therein.

It's only been tested on mac, but it should work on any nixy system.  The only non-standard tool it uses is `tree`, and that's nonessential.  If you want it anyway, ond you're on n a mac, do `brew install tree`.

Obviously, this script also requires protobufs, grpc, and the go-bindings thereof.  For installation instructions there, I'll defer to [this page](https://grpc.io/docs/quickstart/go.html).

If you're able to drop this gen proto script into the root of the example directory they provide, rename `helloworld` to `helloworldpb`, and the script executes without complain, then you're probably good to use this everywhere!
