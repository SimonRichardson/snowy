# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2017-08-11T17:00:47+01:00
# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('4a33ac76-9f81-4d93-8384-b8596ab8ae84')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 143.482µs
            X-Query-Tags: abc,def,g
            X-Resourceid: 4a33ac76-9f81-4d93-8384-b8596ab8ae84

    + Body

            {
                "author_id": "f3316ecc-77fe-4dab-b1bf-491639314f86",
                "created_on": "2017-08-11T17:00:47+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "4a33ac76-9f81-4d93-8384-b8596ab8ae84",
                "resource_size": 10,
                "tags": [
                    "abc",
                    "def",
                    "g"
                ]
            }

# GET /multiple

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('4a33ac76-9f81-4d93-8384-b8596ab8ae84')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 60.469µs
            X-Query-Tags: abc,def,g
            X-Resourceid: 4a33ac76-9f81-4d93-8384-b8596ab8ae84

    + Body

            [
                {
                    "author_id": "f3316ecc-77fe-4dab-b1bf-491639314f86",
                    "created_on": "2017-08-11T17:00:47+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "4a33ac76-9f81-4d93-8384-b8596ab8ae84",
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
            X-Duration: 124.155µs

    + Body

            {
                "resource_id": "4a33ac76-9f81-4d93-8384-b8596ab8ae84"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('4a33ac76-9f81-4d93-8384-b8596ab8ae84')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 155.919µs
            X-Resourceid: 4a33ac76-9f81-4d93-8384-b8596ab8ae84

    + Body

            {
                "resource_id": "00000000-0000-0000-0000-000000000000"
            }

