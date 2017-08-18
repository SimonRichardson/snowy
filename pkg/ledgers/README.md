# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2017-08-18T15:43:11+01:00
# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('56ef9ef0-34d2-4509-9e00-fa3b21245033')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 74.008µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 56ef9ef0-34d2-4509-9e00-fa3b21245033

    + Body

            {
                "author_id": "ab37eb40-8363-4950-b8d7-4f25ecd4351a",
                "created_on": "2017-08-18T15:43:11+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "56ef9ef0-34d2-4509-9e00-fa3b21245033",
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
            resource_id ('56ef9ef0-34d2-4509-9e00-fa3b21245033')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 35.855µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 56ef9ef0-34d2-4509-9e00-fa3b21245033

    + Body

            [
                {
                    "author_id": "ab37eb40-8363-4950-b8d7-4f25ecd4351a",
                    "created_on": "2017-08-18T15:43:11+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "56ef9ef0-34d2-4509-9e00-fa3b21245033",
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
            X-Duration: 91.211µs

    + Body

            {
                "resource_id": "56ef9ef0-34d2-4509-9e00-fa3b21245033"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('56ef9ef0-34d2-4509-9e00-fa3b21245033')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 75.602µs
            X-Resource-Id: 56ef9ef0-34d2-4509-9e00-fa3b21245033

    + Body

            {
                "resource_id": "00000000-0000-0000-0000-000000000000"
            }

