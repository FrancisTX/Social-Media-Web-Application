# Use goreman to run

web: ./web/web

server: ./server/server --userport 10380 --postport 11380 --followport 12380

userstorage1: ./server/storage/storage --storage user --id 1 --cluster http://127.0.0.1:10379,http://127.0.0.1:20379,http://127.0.0.1:30379 --port 10380
userstorage2: ./server/storage/storage --storage user --id 2 --cluster http://127.0.0.1:10379,http://127.0.0.1:20379,http://127.0.0.1:30379 --port 20380
userstorage3: ./server/storage/storage --storage user --id 3 --cluster http://127.0.0.1:10379,http://127.0.0.1:20379,http://127.0.0.1:30379 --port 30380

poststorage1: ./server/storage/storage --storage post --id 1 --cluster http://127.0.0.1:11379,http://127.0.0.1:21379,http://127.0.0.1:31379 --port 11380
poststorage2: ./server/storage/storage --storage post --id 2 --cluster http://127.0.0.1:11379,http://127.0.0.1:21379,http://127.0.0.1:31379 --port 21380
poststorage3: ./server/storage/storage --storage post --id 3 --cluster http://127.0.0.1:11379,http://127.0.0.1:21379,http://127.0.0.1:31379 --port 31380

followstorage1: ./server/storage/storage --storage follow --id 1 --cluster http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379 --port 12380
followstorage2: ./server/storage/storage --storage follow --id 2 --cluster http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379 --port 22380
followstorage3: ./server/storage/storage --storage follow --id 3 --cluster http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379 --port 32380