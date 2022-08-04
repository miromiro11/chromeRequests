# 32-bit
GOOS=windows CGO_ENABLED=1 GOARCH=386 go build -buildmode=c-shared -o ./compiled/library32-windows.so 

# 64-bit
GOOS=windows CGO_ENABLED=1 GOARCH=amd64 go build -buildmode=c-shared -o ./compiled/library64-windows.so 

darwin 
GOOS=darwin CGO_ENABLED=1 GOARCH=arm64 go build -buildmode=c-shared -o ./compiled/library-Darwin-64.so