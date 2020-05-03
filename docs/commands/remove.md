remove
==========

Remove a bucket by its id

## Syntax
```
hog remove {id} [flags]
```

This command takes only one id as argument. It will check if the bucket exists inside `hog.yml` and will remove the 
reference from the file. By default this command runs immediately, if you which to delay the deletion use one of the 
flags described below.

## Examples

```
hog remove Xgo8gKM
hog remove 3lWO7rn --deadline 9:41AM
hog remove 4i7RAkR --ttl 10m
hog remove 2iez0Wa --ttl 10m -d
hog remove RM0UoHD --ttl 10m --detached
```

## Flags
| Flag                   | Short code | Description                            | 
| -------                | ------     | -------                                | 
| [ttl](#ttl)            |            | Remove a bucket after a period of time |
| [deadline](#deadline)  |            | Remove a bucket at a particular time   |
| [detached](#detached)  | d          | Run the command in detached mode       |

### ttl

This flag will set a `ttl` [Time To live][time-to-live]. Which is a duration that the file is going to be available
in `hog`. This flag will delay the execution of the program until the duration set in the flag is passed. To avoid
waiting for the delay run with `detached`.


```
hog remove 4i7RAkR --ttl 10m
hog remove 4i7RAkR --ttl 10m -d
hog remove 4i7RAkR --ttl 10m --detached
```

### deadline

This flag will set a particular time when the file must be deleted. This flag will delay the execution of the program 
until the duration set in the flag is passed. To avoid waiting for the delay run with `detached`.


```
hog remove 3lWO7rn --deadline 9:41AM
hog remove 3lWO7rn --deadline 9:41AM -d 
hog remove 3lWO7rn --deadline 9:41AM --detached
```

### detached
Execute the command in detached mode.

Example:
```
hog remove 9R59TwD -d
hog remove 4i7RAkR --ttl 10m -d
hog remove 4i7RAkR --ttl 10m --detached
hog remove 3lWO7rn --deadline 9:41AM -d
hog remove 3lWO7rn --deadline 9:41AM --detached
```

[time-to-live]: https://www.cloudflare.com/learning/cdn/glossary/time-to-live-ttl/