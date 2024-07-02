
# Debug Go CLI Application
go/debug:
	cd ./geo && dlv debug --headless --listen :4040 main.go -- load --dir ../assets/example_data

# Build Go CLI Application and add to $GOBIN
go/build:
	cd ./geo && go build -o /home/ubuntu/go/bin/geo main.go