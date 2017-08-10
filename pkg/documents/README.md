# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2017-08-10T14:20:59+01:00
# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('2bbf74c7-1444-42aa-9a65-fc4da29e340b')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 67.884µs
            X-Query-Tags: abc,def,g
            X-Resourceid: 2bbf74c7-1444-42aa-9a65-fc4da29e340b

    + Body

            {
                "author_id": "f22449d3-090e-41d4-b7bf-95bd207914e0",
                "created_on": "2017-08-10T14:20:59+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_id": "2bbf74c7-1444-42aa-9a65-fc4da29e340b",
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
            resource_id ('2bbf74c7-1444-42aa-9a65-fc4da29e340b')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 37.253µs
            X-Query-Tags: abc,def,g
            X-Resourceid: 2bbf74c7-1444-42aa-9a65-fc4da29e340b

    + Body

            [
                {
                    "author_id": "f22449d3-090e-41d4-b7bf-95bd207914e0",
                    "created_on": "2017-08-10T14:20:59+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_id": "2bbf74c7-1444-42aa-9a65-fc4da29e340b",
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
            X-Duration: 91.979µs

    + Body

            {
                "resource_id": "2bbf74c7-1444-42aa-9a65-fc4da29e340b"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('2bbf74c7-1444-42aa-9a65-fc4da29e340b')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 65.854µs
            X-Resourceid: 2bbf74c7-1444-42aa-9a65-fc4da29e340b

    + Body

            {
                "resource_id": "00000000-0000-0000-0000-000000000000"
            }

