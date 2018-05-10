# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2018-05-10T20:00:35+01:00
# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('b8fea624-4231-4ddc-b2cc-4b6a41831b03')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 84.24µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: b8fea624-4231-4ddc-b2cc-4b6a41831b03

    + Body

            {
                "author_id": "463b476f-9f4c-4a36-8bb7-da48802c4105",
                "created_on": "2018-05-10T20:00:35+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "b8fea624-4231-4ddc-b2cc-4b6a41831b03",
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
            resource_id ('b8fea624-4231-4ddc-b2cc-4b6a41831b03')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 72.407µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: b8fea624-4231-4ddc-b2cc-4b6a41831b03

    + Body

            [
                {
                    "author_id": "463b476f-9f4c-4a36-8bb7-da48802c4105",
                    "created_on": "2018-05-10T20:00:35+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "b8fea624-4231-4ddc-b2cc-4b6a41831b03",
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
                "author_id": "b8fea624-4231-4ddc-b2cc-4b6a41831b03",
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
            X-Duration: 171.48µs

    + Body

            {
                "resource_id": "b8fea624-4231-4ddc-b2cc-4b6a41831b03"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('b8fea624-4231-4ddc-b2cc-4b6a41831b03')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

    + Body

            {
                "author_id": "b8fea624-4231-4ddc-b2cc-4b6a41831b03",
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
            X-Duration: 56.737µs
            X-Resource-Id: b8fea624-4231-4ddc-b2cc-4b6a41831b03

    + Body

            {
                "resource_id": "b8fea624-4231-4ddc-b2cc-4b6a41831b03"
            }

# PUT /fork/

+ Request
    + Parameters

            resource_id ('b8fea624-4231-4ddc-b2cc-4b6a41831b03')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

    + Body

            {
                "author_id": "b8fea624-4231-4ddc-b2cc-4b6a41831b03",
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
            X-Duration: 109.797µs
            X-Resource-Id: b8fea624-4231-4ddc-b2cc-4b6a41831b03

    + Body

            {
                "resource_id": "203bf17e-27bb-4c49-b39e-5fc1290c301f"
            }

