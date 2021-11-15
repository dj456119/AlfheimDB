module github.com/AlfheimDB

go 1.16

require (
	github.com/BBVA/raft-badger v1.1.0 // indirect
	github.com/google/pprof v0.0.0-20211108044417-e9b028704de0 // indirect
	github.com/hashicorp/raft v1.3.2
	github.com/hashicorp/raft-boltdb v0.0.0-20210422161416-485fa74b0b01
	github.com/jinzhu/configor v1.2.1
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/tidwall/redcon v1.4.2 // indirect
	golang.org/x/sys v0.0.0-20211112164355-7580c6e521dc // indirect
	google.golang.org/grpc v1.42.0
	k8s.io/client-go v0.22.3
)

//replace github.com/hashicorp/raft => ../raft
