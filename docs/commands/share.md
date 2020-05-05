share
==========

Share bucket by its id

## Syntax
```
hog share {id} [flags]
```

This command takes only one id as argument. It will check if the bucket exists and return a share url for that bucket 
using the configuration in `.hog.yml`.

## Examples

```
hog share Xgo8gKM
hog share Xgo8gKM -q
hog share Xgo8gKM --qr
hog share Xgo8gKM -p 80
hog share Xgo8gKM --port 80
hog share Xgo8gKM --protocol http
hog share Xgo8gKM --protocol https
hog share Xgo8gKM --domain 192.168.1.1
```

## Flags
| Flag                   | Short code | Description                               | 
| -------                | ------     | -------                                   | 
| [qr](#qr)              | q          | Return a qr code with the url as response |
| [protocol](#protocol)  |            | Overwrite the protocol in the url         |
| [domain](#domain)      |            | Overwrite the domain in the url           |
| [port](#port)          | p          | Overwrite the port in the url             |

### qr

This flag will make the command return a `qr` code with the value of the url.

```
hog share Xgo8gKM -q
hog share Xgo8gKM --qr
```

### protocol

This flag will overwrite the protocol in the share url. 

_This do not overwrite the value in the file `.hog.yml`_

The only valid values are:
- `http`
- `https`

```
hog share Xgo8gKM --protocol http
hog share Xgo8gKM --protocol https
```

### domain

This flag will overwrite the domain in the share url. 

_This do not overwrite the value in the file `.hog.yml`_

```
hog share Xgo8gKM --domain 192.168.1.101
hog share Xgo8gKM --domain example.com
```

### port

This flag will overwrite the port in the share url. 

_This do not overwrite the value in the file `.hog.yml`_

```
hog share Xgo8gKM -p 5000
hog share Xgo8gKM --port 5000
```