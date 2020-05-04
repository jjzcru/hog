get
==========

Get hog configuration values.

_This values comes from the file `.hog.yml`_

## Syntax
```
hog get [command]
```

## Example
```
hog get protocol
hog get domain
hog get port
```

## Command

### protocol

Get protocol value. Should be `http` or `https`.

#### Syntax
```
hog get protocol
```

### domain

Get domain value. 

_This is the domain from which the devices are going to reach the server, can be an ip like `192.168.1.101` or a 
public domain like `example.com`. This value is used for sharing functionalities_

#### Syntax
```
hog get domain
```

### port

Get port value. 

_This value should be a number_

#### Syntax
```
hog get port
```