add
==========

Group files in a bucket

## Syntax
```
hog add [files] [flags]
```

This command takes at least one file as an argument, it will group all the files and group them inside a bucket.
Then it will return the `id` for the bucket.

## Examples

```
hog add test.jpg
hog add test_1.jpg ./test_2.jpb
hog add test.jpg /home/example/download/file.pdf
hog add test.jpg /home/example/download/file.pdf --url
hog add test.jpg /home/example/download/file.pdf -u
```

## Flags
| Flag         | Short code | Description                            | 
| -------      | ------     | -------                                | 
| [ttl](#ttl)  |            | Remove a bucket after a period of time |
| [url](#url)  | u          | Return a share url as response         |

### ttl

This flag will set a `ttl` [Time To live][time-to-live]. Which is a duration that the file is going to be available
in `hog`. Under the hood is passing the value of the flag to the command `remove`.

```
hog add test.jpg --ttl 10s
hog add test.jpg test_1.png --ttl 1m
```

### url

This flag if enable will return an url to share, created by using the configuration inside `.hog.yml`. 

```
hog add test.jpg --url
hog add test.jpg test_1.png -u
```

[time-to-live]: https://www.cloudflare.com/learning/cdn/glossary/time-to-live-ttl/