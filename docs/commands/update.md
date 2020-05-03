update
==========

Update the files in a bucket by its id

## Syntax
```
hog update {id} [files] [flags]
```

This command takes only one id as argument. It will check if the bucket exists inside `.hog.yml` and will overwrite
the files inside the bucket with the files provided.

## Examples

```
hog update Xgo8gKM file.jpg
hog update Xgo8gKM file.jpg file.png
hog update Xgo8gKM file.jpg /home/root/file.png
```
