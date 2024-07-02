
go/debug:
	cd ./geo && dlv debug --headless --listen :4040 main.go

go/build:
	cd ./geo && go build -o /home/ubuntu/go/bin/geo main.go 