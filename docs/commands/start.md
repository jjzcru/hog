start
==========

Start hog service

## Syntax
```
hog start [flags]
```

This command use the `hog.yml` file to fetch the bucket ids that are going to be serve to download.
To download any bucket you need to call the endpoint `/download/{BucketID}` depending on the amount of files and 
the path the download is going to be as following:

- If the bucket only has one path, and that path is a file, the service is going to serve the file directly from
the path.

- If the bucket only has one path, and that path is a directory, it will create a `.zip` file with a random name, with
the content of the directory, inside the `tmp` directory on the os and it will delete the file once the download is 
finished.

- If the bucket has more than one path, it will create a `.zip` file wit a random name, with the content of each path
as entry, inside the `tmp` directory on the os and it will delete the file once the download is finished.

In the case where `.zip` files are being generated, the temporary `.zip` files are going to be deleted if the 
request gets cancelled.

By default `hog` runs on port `1618`.

## Examples

```
hog start
hog start -d
hog start -p 3000 -d
hog start --port 3000
hog start --port 3000 --detached
```

## Flags
| Flag                   | Short code | Description                              | 
| -------                | ------     | -------                                  | 
| [port](#port)          | p          | Port where the server is going to run    |
| [detached](#detached)  | d          | Run in detached mode and return the PID  |

### port

Specify the port that is going to be used by the server. If not set is going to use the port `1618` by default.

```
hog start -p 3000
hog start --port 3000
```

### detached
Execute the command in detached mode.

Example:
```
hog start -d
hog start --detached
```