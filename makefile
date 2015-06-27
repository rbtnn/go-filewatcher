
ifeq ($(OS),Windows_NT)
  TARGET = filewatcher.exe
else
  TARGET = filewatcher
endif

$(TARGET): src/main.go
	go build -o $(TARGET) src/main.go

