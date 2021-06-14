module github.com/client

go 1.13

replace github.com/utils => ../utils

replace github.com/crdt => ../crdt

require (
	github.com/crdt v0.0.0-00010101000000-000000000000
	github.com/utils v0.0.0-00010101000000-000000000000
)
