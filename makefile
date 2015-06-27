
ifeq ($(OS),Windows_NT)
  TARGET = filewatcher.exe
else
  TARGET = filewatcher
endif

$(TARGET): src/filewatcher.go
	go build -o $(TARGET) src/filewatcher.go

