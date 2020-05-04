[![Build Status](https://travis-ci.com/jjzcru/hog.svg?branch=master)](https://travis-ci.com/jjzcru/hog)
[![Coverage Status](https://coveralls.io/repos/github/jjzcru/hog/badge.svg?branch=master)](https://coveralls.io/github/jjzcru/hog?branch=master)
![Release](https://github.com/jjzcru/hog/workflows/Release/badge.svg?branch=master)

Hog
==========

Hog 🐗 is a file sharing tool, it enables you to share your files directly to other people without the need to upload
the files to a server and then download them from there. 

You can single files, directory or multiple file paths without any configuration in your file system. Just select say
which files you would like to share and it will share it for you, is that simple.

Since it's written in [Go][go], most of the commands runs across multiple operating systems (`Linux`, `macOS`, 
`Windows`).

*Why should i use this?* You can watch some [Use Cases](#use-cases)

## Table of contents
  * [Getting Started](#getting-started)
    + [Installation](#installation)
    + [Usage](#usage)
  * [Syntax](#syntax)
  * [Use Cases](#use-cases)
  * [Commands](#commands)
  * [Roadmap](#roadmap)
  * [Changelog][changelog]
  * [Releases][releases]

## Getting Started
The main use case for `hog` is that you are able to share your files with other people directly from your machine
without any intermediary _(Dropbox, Google Drive, WeTransfer)_.

`hog` uses a file called `.hog.yml` to store the references from your files, by default this file is located in your
home directory but you can change the directory used by `hog` by setting the `env` variable `HOG_PEN`. 

`hog` groups your files in `buckets` which are just a list of paths inside your file system for the files you would
like to share, and each `bucket` is reference by an unique id. 

If the `bucket` only contains one path and the path is a file, `hog` will serve that file directly from your file 
system, but if the path is a directory, it will create a temporary `.zip` file on your `temp` directory.

If the `bucket` has multiples path, `hog` will create a temporary `.zip` file with the reference for all the paths.

_In case you want to share multiple files that are big in size, is suggested that you shared them as different buckets 
to avoid the duplication of the files while creating the `.zip` files._

### Installation

#### Bash
Installation with `cURL` and `sh` thanks to the project [Go Binaries][gobinaries].
```
curl -sf https://gobinaries.com/jjzcru/hog | sh
```

#### Download 
1. Grab the latest binary of your platform from the [Releases](https://github.com/jjzcru/hog/releases) page.
2. If you are running on `macOS` or `Linux`, run `chmod +x hog` to give `executable` permissions to the binary. If you
are on `windows` you can ignore this step.
3. Add the binary to `$PATH`.
4. Run `hog version` to make sure that the binary is installed.

### Usage

#### Start Server

First you need to start the server so the application is able to serve the files.

```
hog start -d
```

By default the application will run at port `1618` and the endpoint for download is `/download/{BucketID}`. For more 
information about how the `start` command works go to [start][start] documentation.

#### Add files

Now you need to add files to serve, for this go to your terminal and navigate to file or directory that you want to
share and run the `add` command. 
```
hog add file.png ./file.jpg /home/root/file.pdf
```

You can add multiple files which are going to be group together as a bucket, you can add relative path, name of files
or directories or absolute paths. This command will return the `BucketID` that was created.

#### Share your bucket

To share, after you run the `add` command, you need to use the `BucketID` to generate the link. Lets say that the server 
is running on port `1618` and the `BucketID` generated is `2iez0Wa`. Now you just need to create a url that targets
your computer in that port.

```
// For localhost
http://localhost:1618/download/2iez0Wa 

// For IP
http://192.168.1.101:1618/download/2iez0Wa

// For domain
http://my.domain.com:1618/download/2iez0Wa
```

## Syntax
The syntax consists of:
- `domain`: The domain that is use to reference the address of the machine, this value is used by some sharing 
functionalities. The default value is `localhost`.
- `protocol`: This is the protocol that is going to be used to serve the files, this value is used by some sharing 
functionalities. The only valid values are `http` or `https`. The default value is `http`.
- `port`: Port in which the server is running. The default value is `1618`.
- `buckets`: This is a hash map where the keys are unique values that reference a `bucket` and the values are a list of 
paths.

### Example

```yml
domain: localhost
protocol: http
port: 1618
buckets:
  2iez0Wa:
  - /home/root/Downloads/file.pdf
  3lWO7rn:
  - /home/root/Documents
  4i7RAkR:
  - /home/root/Documents/file.pdf
  - /home/root/Documents/file.epub

```

## Use Cases

The goal of `hog` is to share files directly from your machine, if you have a file that you want to share to someone 
on the same network, you could share a link like this one `http://192.168.1.101:1618/download/4i7RAkR` and anyone on 
the network will be able to download that file.

Another example would be, say you want to share a file with someone remotly over the internet, you could use a tool 
like [ngrok][ngrok] to share your traffic and they will be able to download the file directly from your computer.

Another use case is that you want to share a file from your computer to your smarthphone or tablet, you don't want 
to go search for a cable, to then connect your PC, go to your file system and copy the file. Instead just send you the 
link and the file will automatically download.

## Commands

| Command           | Description                            | Syntax                           |
| -------           | ------                                 | -------                          |
| [add][add]        | Group files in a bucket                | `hog add [files] [flags]`        |
| [bucket][bucket]  | Display the buckets and their files    | `hog bucket [flags]`             |
| [get][get]        | Get hog configuration values           | `hog get [command]`              |
| [remove][remove]  | Remove a bucket by its id              | `hog remove {id} [flags]`        |
| [set][set]        | Set hog configuration values           | `hog set [command]`              |
| [start][start]    | Start hog service                      | `hog start [flags]`              |
| [update][update]  | Update the files in a bucket by its id | `hog update {id} [files] [flags]`|
| [version][version]| Display version number                 | `hog version [flags]`            |


## Roadmap
Each release has a particular idea in mind and the tasks inside that release are focusing on that main idea.

To learn more about the progress and what is being planned go to [Projects][projects].

[go]: https://golang.org/
[gobinaries]: https://github.com/tj/gobinaries
[ngrok]: https://ngrok.com/

[releases]: https://github.com/jjzcru/hog/releases
[changelog]: CHANGELOG.md
[projects]: https://github.com/jjzcru/hog/projects

[add]: docs/commands/add.md
[bucket]: docs/commands/bucket.md
[get]: docs/commands/get.md
[remove]: docs/commands/remove.md
[set]: docs/commands/set.md
[start]: docs/commands/start.md
[update]: docs/commands/update.md
[version]: docs/commands/version.md
