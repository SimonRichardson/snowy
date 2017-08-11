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
