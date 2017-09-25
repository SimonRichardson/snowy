# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2017-09-25T14:38:41+01:00
# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('1a9c2a4b-b19b-468d-86db-6fd18e8c58c0')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 48.615µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 1a9c2a4b-b19b-468d-86db-6fd18e8c58c0

    + Body

            {
                "author_id": "9aaecc9c-b7ac-45e7-8014-a849baaf2c3b",
                "created_on": "2017-09-25T14:38:41+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "1a9c2a4b-b19b-468d-86db-6fd18e8c58c0",
                "resource_size": 10,
                "tags": [
                    "abc",
                    "def",
                    "g"
                ]
            }

# GET /revisions/

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('1a9c2a4b-b19b-468d-86db-6fd18e8c58c0')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 54.883µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 1a9c2a4b-b19b-468d-86db-6fd18e8c58c0

    + Body

            [
                {
                    "author_id": "9aaecc9c-b7ac-45e7-8014-a849baaf2c3b",
                    "created_on": "2017-09-25T14:38:41+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "1a9c2a4b-b19b-468d-86db-6fd18e8c58c0",
                    "resource_size": 10,
                    "tags": [
                        "abc",
                        "def",
                        "g"
                    ]
                }
            ]

# POST /

+ Request
    + Parameters


    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 66.357µs

    + Body

            {
                "resource_id": "1a9c2a4b-b19b-468d-86db-6fd18e8c58c0"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('1a9c2a4b-b19b-468d-86db-6fd18e8c58c0')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 100.346µs
            X-Resource-Id: 1a9c2a4b-b19b-468d-86db-6fd18e8c58c0

    + Body

            {
                "resource_id": "1a9c2a4b-b19b-468d-86db-6fd18e8c58c0"
            }

# PUT /fork/

+ Request
    + Parameters

            resource_id ('1a9c2a4b-b19b-468d-86db-6fd18e8c58c0')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 79.189µs
            X-Resource-Id: 1a9c2a4b-b19b-468d-86db-6fd18e8c58c0

    + Body

            {
                "resource_id": "58613ba7-a969-47fb-8251-0d428e85a170"
            }

