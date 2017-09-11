# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2017-09-10T19:32:41+01:00
# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('13077bb5-d804-448a-8079-ffa0bc448d39')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 68.492µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 13077bb5-d804-448a-8079-ffa0bc448d39

    + Body

            {
                "author_id": "7dd0b064-226f-4f25-96c6-bc1d0f9501fb",
                "created_on": "2017-09-10T19:32:41+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "parent_id": "00000000-0000-0000-0000-000000000000",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "13077bb5-d804-448a-8079-ffa0bc448d39",
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
            resource_id ('13077bb5-d804-448a-8079-ffa0bc448d39')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 42.629µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 13077bb5-d804-448a-8079-ffa0bc448d39

    + Body

            [
                {
                    "author_id": "7dd0b064-226f-4f25-96c6-bc1d0f9501fb",
                    "created_on": "2017-09-10T19:32:41+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "parent_id": "00000000-0000-0000-0000-000000000000",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "13077bb5-d804-448a-8079-ffa0bc448d39",
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
            X-Duration: 89.739µs

    + Body

            {
                "resource_id": "13077bb5-d804-448a-8079-ffa0bc448d39"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('13077bb5-d804-448a-8079-ffa0bc448d39')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 66.83µs
            X-Resource-Id: 13077bb5-d804-448a-8079-ffa0bc448d39

    + Body

            {
                "resource_id": "00000000-0000-0000-0000-000000000000"
            }

