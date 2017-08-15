# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2017-08-15T17:04:51+01:00
# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('63a0f8e8-20b6-4715-a9ea-52ef5dd4d047')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 66.805µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 63a0f8e8-20b6-4715-a9ea-52ef5dd4d047

    + Body

            {
                "author_id": "1a821d3c-1baf-4fcf-8720-4b0a13bed786",
                "created_on": "2017-08-15T17:04:51+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "63a0f8e8-20b6-4715-a9ea-52ef5dd4d047",
                "resource_size": 10,
                "tags": [
                    "abc",
                    "def",
                    "g"
                ]
            }

# GET /multiple/

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('63a0f8e8-20b6-4715-a9ea-52ef5dd4d047')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 35.025µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 63a0f8e8-20b6-4715-a9ea-52ef5dd4d047

    + Body

            [
                {
                    "author_id": "1a821d3c-1baf-4fcf-8720-4b0a13bed786",
                    "created_on": "2017-08-15T17:04:51+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "63a0f8e8-20b6-4715-a9ea-52ef5dd4d047",
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
            X-Duration: 81.122µs

    + Body

            {
                "resource_id": "63a0f8e8-20b6-4715-a9ea-52ef5dd4d047"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('63a0f8e8-20b6-4715-a9ea-52ef5dd4d047')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 62.453µs
            X-Resource-Id: 63a0f8e8-20b6-4715-a9ea-52ef5dd4d047

    + Body

            {
                "resource_id": "00000000-0000-0000-0000-000000000000"
            }

