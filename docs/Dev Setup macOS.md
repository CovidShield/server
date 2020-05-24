# Developer Setup for macOS

If you're at Shopify, it's probably best that you follow the setup instructions in [CONTRIBUTING.md](../CONTRIBUTING.md). For external folks that wish to build/run locally on macOS, here's instructions.

This has been tested with the following environment:

- macOS 10.15.4
- [Homebrew](https://brew.sh/)

## Setup Dependencies

### Install Go Lang

The application has been tested with the Go programming language at version 1.14.x. This installs that version explicitly.

`brew install go@1.14`

### Install Ruby

This sets up your Mac to run version 2.6.5 of Ruby, which is what the application has been tested at as of writing.

1. Install OpenSSL: `brew install openssl`
1. Install Rbenv: `brew install rbenv ruby-build`
1. Add the following to your `.zshrc` or `.bashrc` file:
    ``` 
    RUBY_CONFIGURE_OPTS="--with-openssl-dir=$(brew --prefix openssl)"`
    `eval "$(rbenv init -)"
    ```
1. Restart your terminal
1. Install Ruby 2.6.5: `rbenv install 2.6.5`
1. Move into the clone of this repo
1. Set the version of Ruby to 2.6.5: `rbenv local 2.6.5` (this will drop the file `.ruby-version`)
1. Install Bundler: `gem install bundler`

### Install Protocol Buffers

`brew install protobuf`

### Install Protocol Buffer Generator for Go

`brew install protoc-gen-go`

### Install Sodium

`brew install libsodium`

### Install Docker

(Note: The latest release of Docker for Desktop has some issues on macOS 10.15.4. These instructions install a prior version.)

1. Go to the [Docker for Mac Stable release notes](https://docs.docker.com/docker-for-mac/release-notes/)
1. Download Docker Desktop Community 2.2.0.5.
1. Copy the `Docker.app` file to `/Applications`
1. Run `Docker.app`

### Install MySQL Build Dependencies

These are needed for automated tests.

1. Install MySQL client (for its header files): `brew install mysql`
1. Install Ruby driver: `gem install mysql2 -- --with-cflags=\"-I/usr/local/opt/openssl@1.1/include\" --with-ldflags=\"-L/usr/local/opt/openssl@1.1/lib\"`

## Run Database (in Docker)

This runs a MySQL database at `localhost` on port `3306` with a DB root user password of `somepasswordhere`. As of writing, this is MySQL 8.x.

1. Go into the cloned copy of this repo
1. Create file for MySQL data files: `mkdir mysql-data`
1. Run MySQL: `docker run -it --rm -p 3306:3306 --name covidshield-db -v $PWD/mysql-data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=somepasswordhere -d mysql`

## Build CovidShield Server

1. Go into the cloned copy of this repo
1. Run the build: `make`

This will place built files into `./build/release`.

## Setup the Database

1. Setup some dev environment variables: `source scripts/dev_env_vars.sh`
1. Setup the database: `./build/release/key-retrieval migrate-db`

This will leave up the Key Retrieval server running at `http://localhost:8000`. You should see something like:

```
... (more stuff) ...

INFO[0000] released table lock on schema_migrations      component=db
INFO[0000] migrations done                               component=db
INFO[0000] running                                       component=expiration uuid=40dc0a34-8909-4ddd-44f2-40d2fd3d3cc4
INFO[0000] deleted old diagnosis keys                    component=expiration count=0 uuid=40dc0a34-8909-4ddd-44f2-40d2fd3d3cc4
INFO[0000] deleted old encryption keys                   component=expiration count=0 uuid=40dc0a34-8909-4ddd-44f2-40d2fd3d3cc4
INFO[0000] starting server                               bind="0.0.0.0:8001" component=srvutil
INFO[0000] started server                                addr="[::]:8001" bind="0.0.0.0:8001" component=srvutil
INFO[0030] running                                       component=expiration uuid=5d988d3d-a4d9-4ca5-7342-9604b1ac109c
INFO[0030] deleted old diagnosis keys                    component=expiration count=0 uuid=5d988d3d-a4d9-4ca5-7342-9604b1ac109c
INFO[0030] deleted old encryption keys                   component=expiration count=0 uuid=5d988d3d-a4d9-4ca5-7342-9604b1ac109c
```

Kill it with `CTRL+C` after you see the above.

## Run the Servers

### Run Key Submission

This runs the Key Submission server at `http://localhost:8000`.

1. Open a new terminal and move into the root of this cloned repo
1. Setup some dev environment variables: `source scripts/dev_env_vars.sh`
1. Run the server `PORT=8000 ./build/release/key-submission`

### Run Key Retrieval

This runs the Key Retrieval server at `http://localhost:8001`.

1. Open a new terminal and move into the root of this cloned repo
1. Setup some dev environment variables: `source scripts/dev_env_vars.sh`
1. Run the server: `PORT=8001 ./build/release/key-retrieval`

## Run Tests

1. Setup the following environment variables:
    ```
    export DB_USER=<username>
    export DB_PASS=<password>
    export DB_HOST=<hostname>
    export DB_NAME=<test database name>
    ```
1. Install dependencies for the tests: `bundle install`
1. Run the tests: `make test`