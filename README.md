# Prometheus File Exporter

## Feature(s)
Counts files in given directory

## Usage

### Flags

--observe <path to observed dir> -- sets directory path which should be monitored for file count

--exclude <path to excluded dir> -- works with --observe. If path given in --observe contains path given in --exclude then that path is being ignored in file_count.

--port -- Configures listening port for exporter. Defaults to 9111

### Example
Example with given tree:
```
a
|____b
| |____c
| | |____C
| |____B
|____A
```
Directory 'a' contains file 'A' and subdir 'b'. We can see that there are 3 total files represented by numbers. If we pass plain --observe=/a the result of path_file_count will be 3. In case we pass --exclude=/a/b/c --observe=/a value of path_file_count will be 2  

### Command
prometheus_file_exporter [--observe=<>] [--exclude=<>] [--port=<>]