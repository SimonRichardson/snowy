# snowy

![snowy](images/snowy.jpg)

Documents and contents repository

## Badges

[![travis-ci](https://travis-ci.org/trussle/snowy.svg?branch=master)](https://travis-ci.org/trussle/snowy)
[![Coverage Status](https://coveralls.io/repos/github/trussle/snowy/badge.svg?branch=master)](https://coveralls.io/github/trussle/snowy?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/trussle/snowy)](https://goreportcard.com/report/github.com/trussle/snowy)

## Introduction

Snowy is a append only ledger for document contents, that allow you to associate
tags to a piece of content that can be queryable from the rest end point. The
snowy application is split into two distinct parts, the documents (ledger) and
the associated content for that ledger entry.

Modification of documents and contents is not possible, instead new entities are
required to be inserted in an append only fashion, where a full revision and
audit trail can be viewed for each document.

## Setup

### Local development

Snowy expects that you have a `$GOPATH` configured correctly and that you've
installed it using `go get ...`. Once these are done, it should be as simple as
`make install`, which will get all the correct dependencies for you to start 
working with the code.

### Integration development

Integration development (testing) requires both `docker` and `docker-compose` to
be installed, if testing against the remote filesystem (`Amazon s3`), then 
you'll need to set some environmental variables. There are two ways to do this:

  1. Export the env vars, so they can be read
  2. Put a `.env` file in the root of the snowy project with the following 
  values filled in:

```
AWS_ID=
AWS_SECRET=
AWS_BUCKET=
AWS_REGION=
AWS_TOKEN=
```

Then it should be a case of running the following:

```bash
docker-compose up --build -d
make integration-tests
```

This will fetch an instance of `postgresdb` for you and setup all the right 
schema information, along with any other dependencies, before running the 
integration tests.

## API Endpoints

The following contains the documentation for the API end points for Snowy.

### Contents

Contents API is for retrieving files from the underlying storage. The API allows
insertion and selection of files, it is not possible to update an image as all
files are immutable, so new copies of the file are always stored in the storage.

 - [API](pkg/contents/README.md)

### Documents

Documents API is for retrieving the ledgers of all the contents found with in
Snowy. Each new content is appended to a table, then added to the ledger so
querying becomes possible. Modifying a document or content is not possible and
to add new revision to the ledger a new document, content can be be appended
using the same resource_id.

 - [API](pkg/documents/README.md)
