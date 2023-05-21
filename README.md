# Go bunrouter + bun application example

[![build workflow](https://github.com/go-bun/bun-realworld-app/actions/workflows/build.yml/badge.svg)](https://github.com/go-bun/bun-realworld-app/actions)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/cristianuser/go-bun-webserver)](https://pkg.go.dev/github.com/cristianuser/go-bun-webserver)

## Introduction

This project implements RealWorld JSON API as specified in the
[spec](https://github.com/gothinkster/realworld). It was created to demonstrate how to use:

- [Bun Router](https://github.com/uptrace/bunrouter).
- [Bun DB](https://github.com/uptrace/bun).
- [bun/migrate](https://bun.uptrace.dev/guide/migrations.html).
- [bun/dbfixture](https://bun.uptrace.dev/guide/fixtures.html).

## Project structure

The project uses Bun [starter kit](https://bun.uptrace.dev/guide/starter-kit.html) and consists of
the following packages:

- [bunapp](bunapp) package parses configs, establishes DB connections etc.
- [org](org) package manages users and sessions.
- [blog](blog) package manages articles and comments.
- [cmd/bun](cmd/bun) provides CLI commands to run HTTP server and work with DB.
- [cmd/bun/migrations](cmd/bun/migrations) contains database migrations.

The most interesting part for Bun users is probably [article filter](blog/article_filter.go).

## Project bootstrap

Project comes with a `Makefile` that contains following recipes:
- `make start` start docker containers, run migrations and run HTTP server.
- `make db_reset` drops existing database and creates a new one.
- `make test` runs unit tests.
- `make api_test` runs API tests provided by
  [RealWorld](https://github.com/gothinkster/realworld/tree/master/api).

After checking that tests are passing you can run HTTP server:

```shell
go run cmd/bun/main.go -env=dev api
```
