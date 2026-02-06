@echo off
REM Generate Go gRPC stubs from proto file

echo Installing protoc plugins...
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo.
echo Generating Go gRPC code...
protoc --go_out=internal/grpc/pb --go_opt=paths=source_relative --go-grpc_out=internal/grpc/pb --go-grpc_opt=paths=source_relative proto/recommendation.proto

if %ERRORLEVEL% EQU 0 (
    echo ✓ Go gRPC stubs generated successfully
) else (
    echo ✗ Failed to generate stubs. Make sure protoc is installed:
    echo   Download from: https://github.com/protocolbuffers/protobuf/releases
)

pause
