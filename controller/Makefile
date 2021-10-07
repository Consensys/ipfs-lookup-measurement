build controller : cmd/controller/*.go pkg/*/*.go
	go build cmd/controller/*

agent: cmd/agent/*.go pkg/*/*.go
	go build cmd/agent/*

nodes-list.out :
	printf " http://localhost:3030 \n " > nodes-list.out

itest: nodes-list.out controller agent
	./controller -l nodes-list.out

utest:
	go test -v ./...
