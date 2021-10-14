module github.com/xiekeyi98/example/dockerpod

go 1.16

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/xiekeyi98/cpuloadhpa v0.0.0-20211014040257-17bab339e2ca
	go.uber.org/automaxprocs v1.4.0 // indirect
)

//replace github.com/xiekeyi98/cpuloadhpa => ../../
