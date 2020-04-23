# prerequisite

follow the document to run the telegramd in your custom development environment. 
etcd, mysql and redis are required components. for sake of simplicity docker is applied :-)

### install docker
* [Install Docker for Ubuntu.](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
for version 18.04 see [Ubuntu-18-04](https://linuxconfig.org/how-to-install-docker-on-ubuntu-18-04-bionic-beaver)
* [Install Docker for Mac](https://docs.docker.com/docker-for-mac/install/)
* [Install Docker for Windows](https://docs.docker.com/docker-for-windows/install/#start-docker-for-windows)

### run etcd container
to pull and run etcd container enter the following command in the shell:
```
$ docker run --name etcd-docker -d -p 2379:2379 -p 2380:2380 appcelerator/etcd
```

### run mysql container
to create mysql container run the following command:
```
$ docker run --name mysql-docker -p 3306:3306 -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -d mysql:5.7
```
alternatively, if you want create a mysql container with root password use the below command
 - using such a container, before run telegramd modules corresponding config file need to be changed
```
$ docker run --name mysql-docker -p 3306:3306 -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:5.7
```
note that ***my-secret-pw*** is the password to be set for the MySQL root user

to run mysql client for test connection or test a sql command run the following:
```
$ docker exec -it mysql-docker mysql -uroot -p
```
and enter your password defined in previous command:
```
mysql> exit
```

### run redis container
to install redis run the following command
```
$ docker run --name redis-docker -p 6379:6379 -d redis 
```

### start containers
After restart your development environment to start etc, mysql, and redis container run
the following command 
```
$ docker start redis-docker mysql-docker etcd-docker
```

to see current running containers run the following command
```
$ docker ps
```

# build telegramd


### get telegramd

```
$ mkdir $GOPATH/src/github.com/nebulaim/
$ cd $GOPATH/src/github.com/nebulaim/
$ git clone https://github.com/nebulaim/telegramd.git
```

### create DB schema
run the following command to create database
```
$ docker exec -it mysql-docker sh -c 'exec mysql -u root -p -e"CREATE DATABASE nebulaim;"' 
```
 and to create db schema run the following:
 
 1- if root password does not set for mysql container:
 ```
 $ docker exec -i mysql-docker mysql --user=root nebulaim < $GOPATH/src/github.com/nebulaim/telegramd/scripts/nebulaim.sql
 ```
 
 2- if root password is set:
```
$ docker exec -i mysql-docker mysql --user=root --password=my-secret-pw nebulaim < $GOPATH/src/github.com/nebulaim/telegramd/scripts/nebulaim.sql
```
note: ***my-secret-pw*** is the same as defined in run mysql container section

##### 2. set custom password in config files
if password is empty ignore this section otherwise add password to the following files
```
$ $GOPATH/src/github.com/nebulaim/telegramd/server/access/auth_key/auth_key.toml
$ $GOPATH/src/github.com/nebulaim/telegramd/server/sync/sync.toml
$ $GOPATH/src/github.com/nebulaim/telegramd/server/upload/upload.toml
$ $GOPATH/src/github.com/nebulaim/telegramd/service/ocument/document.toml
$ $GOPATH/src/github.com/nebulaim/telegramd/server/biz_server/biz_server.toml
```
set ***my-secret-pw*** in mysql dsn as follow:
```
[[mysql]]
name = "immaster"
dsn = "root:my-secret-pw@/nebulaim?charset=utf8"
...

[[mysql]]
name = "imslave"
dsn = "root:my-secret-pw@/nebulaim?charset=utf8"
...
```

  
 
### build frontend
```
$ cd $GOPATH/src/github.com/nebulaim/telegramd/server/access/frontend
$ go get
$ go build
```

### build session
```
$ cd $GOPATH/src/github.com/nebulaim/telegramd/server/access/session
$ go get
$ go build
```

### build auth_key
```
$ cd $GOPATH/src/github.com/nebulaim/telegramd/server/access/auth_key
$ go get
$ go build
```

### build auth_key
```
$ cd $GOPATH/src/github.com/nebulaim/telegramd/service/auth_session
$ go get
$ go build
```

### build sync
```
$ cd $GOPATH/src/github.com/nebulaim/telegramd/server/sync
$ go get
$ go build
```

### build upload
```
$ cd $GOPATH/src/github.com/nebulaim/telegramd/server/upload
$ go get
$ go build
```

### build document
```
$ cd $GOPATH/src/github.com/nebulaim/telegramd/service/document
$ go get
$ go build
```

### build biz_server
```
$ cd $GOPATH/src/github.com/nebulaim/telegramd/server/biz_server
$ go get
$ go build
```

### set DcOptions
in the following file 
```
$ $GOPATH/src/github.com/nebulaim/telegramd/config.toml
```
replace ipAddress by your IP
```
[[DcOptions]]
ipv6 = false
mediaOnly = false
tcpoOnly = false
cdn = false
static = false
id = 2
ipAddress = "127.0.0.1"
port = 12345
```


### run telegramd modules
```
$ cd $GOPATH/src/github.com/nebulaim/telegramd/server/access/frontend
$ ./frontend

$ cd $GOPATH/src/github.com/nebulaim/telegramd/server/access/auth_key
$ ./auth_key

$ cd $GOPATH/src/github.com/nebulaim/telegramd/service/auth_session
$ ./auth_session

$ cd $GOPATH/src/github.com/nebulaim/telegramd/server/sync
$ ./sync

$ cd $GOPATH/src/github.com/nebulaim/telegramd/server/upload
$ ./upload

$ cd $GOPATH/src/github.com/nebulaim/telegramd/service/document
$ ./document

$ cd $GOPATH/src/github.com/nebulaim/telegramd/server/biz_server
$ ./biz_server

$ cd $GOPATH/src/github.com/nebulaim/telegramd/server/access/session
$ ./session
```

# notes
* if a panic is raised for `http: multiple registrations for /debug/requests` then 
[remove github.com/coreos/etcd/vendor/golang.org/x/net/trace folder](https://github.com/coreos/etcd/issues/9357)


