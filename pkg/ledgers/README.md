# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2018-02-01T16:30:35Z
# GET /revisions/

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('6021617c-16e4-4fa2-9c40-804dba95d6a6')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 26.855µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 6021617c-16e4-4fa2-9c40-804dba95d6a6

    + Body

            [
                {
                    "author_id": "0d20392c-6855-4e03-905a-97fcf9313195",
                    "created_on": "2018-02-01T16:30:35Z",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "6021617c-16e4-4fa2-9c40-804dba95d6a6",
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

    + Body

            {
                "author_id": "6021617c-16e4-4fa2-9c40-804dba95d6a6",
                "name": "document-name",
                "tags": [
                    "abc",
                    "def",
                    "g"
                ]
            }

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 46.426µs

    + Body

            {
                "resource_id": "6021617c-16e4-4fa2-9c40-804dba95d6a6"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('6021617c-16e4-4fa2-9c40-804dba95d6a6')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

    + Body

            {
                "author_id": "6021617c-16e4-4fa2-9c40-804dba95d6a6",
                "name": "document-name",
                "tags": [
                    "abc",
                    "def",
                    "g"
                ]
            }

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 34.823µs
            X-Resource-Id: 6021617c-16e4-4fa2-9c40-804dba95d6a6

    + Body

            {
                "resource_id": "6021617c-16e4-4fa2-9c40-804dba95d6a6"
            }

# PUT /fork/

+ Request
    + Parameters

            resource_id ('6021617c-16e4-4fa2-9c40-804dba95d6a6')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

    + Body

            {
                "author_id": "6021617c-16e4-4fa2-9c40-804dba95d6a6",
                "name": "document-name",
                "tags": [
                    "abc",
                    "def",
                    "g"
                ]
            }

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 32.534µs
            X-Resource-Id: 6021617c-16e4-4fa2-9c40-804dba95d6a6

    + Body

            {
                "resource_id": "409c7128-2214-4b2a-bea5-339adeb20282"
            }

# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('6021617c-16e4-4fa2-9c40-804dba95d6a6')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 56.071µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 6021617c-16e4-4fa2-9c40-804dba95d6a6

    + Body

            {
                "author_id": "0d20392c-6855-4e03-905a-97fcf9313195",
                "created_on": "2018-02-01T16:30:35Z",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "6021617c-16e4-4fa2-9c40-804dba95d6a6",
                "resource_size": 10,
                "tags": [
                    "abc",
                    "def",
                    "g"
                ]
            }

