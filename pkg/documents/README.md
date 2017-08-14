# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2017-08-14T13:06:39+01:00
# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('bb71fd48-62f2-42ba-9a23-567ebb30d80b')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 42.002µs
            X-Query-Tags: abc,def,g
            X-Resourceid: bb71fd48-62f2-42ba-9a23-567ebb30d80b

    + Body

            {
                "author_id": "dd745b70-9ff9-4213-87a7-220ebf558034",
                "created_on": "2017-08-14T13:06:39+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "bb71fd48-62f2-42ba-9a23-567ebb30d80b",
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
            resource_id ('bb71fd48-62f2-42ba-9a23-567ebb30d80b')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 33.488µs
            X-Query-Tags: abc,def,g
            X-Resourceid: bb71fd48-62f2-42ba-9a23-567ebb30d80b

    + Body

            [
                {
                    "author_id": "dd745b70-9ff9-4213-87a7-220ebf558034",
                    "created_on": "2017-08-14T13:06:39+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "bb71fd48-62f2-42ba-9a23-567ebb30d80b",
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
            X-Duration: 71.226µs

    + Body

            {
                "resource_id": "bb71fd48-62f2-42ba-9a23-567ebb30d80b"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('bb71fd48-62f2-42ba-9a23-567ebb30d80b')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 52.4µs
            X-Resourceid: bb71fd48-62f2-42ba-9a23-567ebb30d80b

    + Body

            {
                "resource_id": "00000000-0000-0000-0000-000000000000"
            }

