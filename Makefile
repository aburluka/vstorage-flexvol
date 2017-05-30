SOURCES := $(shell find . 2>&1 | grep -E '.*\.(c|h|go)$$')

.DEFAULT: vstorage-fv

vstorage-fv: $(SOURCES)
	go build -o vstorage-fv .
