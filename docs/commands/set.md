set
==========

Set hog configuration values.

_This command will update the values in the file `.hog.yml`_

## Syntax
```
hog set [command]
```

## Example
```
hog set protocol http
hog set protocol https
hog set domain localhost
hog set domain 192.168.1.101
hog set domain example.com
hog set port 3000
```

## Command

### protocol

Set protocol value. This value can only be `http` or `https`.

#### Syntax
```
hog set protocol http
hog set protocol https
```

### domain

Set domain value. 

_This is the domain from which the devices are going to reach the server, can be an ip like `192.168.1.101` or a 
public domain like `example.com`. This value is used for sharing functionalities_

#### Syntax
```
hog set domain localhost
hog set domain 192.168.1.101
hog set domain example.com
```

### port

Set port value.

_This value must be a number_

#### Syntax
```
hog set port 3000
```