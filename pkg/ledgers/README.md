# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2017-09-04T21:07:10+01:00
# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('80a4772b-339f-4a4e-aefb-4be1b96f0595')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 92.69µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 80a4772b-339f-4a4e-aefb-4be1b96f0595

    + Body

            {
                "author_id": "0734e39b-5272-452b-b989-5ef4720f7ee3",
                "created_on": "2017-09-04T21:07:10+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "80a4772b-339f-4a4e-aefb-4be1b96f0595",
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
            resource_id ('80a4772b-339f-4a4e-aefb-4be1b96f0595')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 38.973µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 80a4772b-339f-4a4e-aefb-4be1b96f0595

    + Body

            [
                {
                    "author_id": "0734e39b-5272-452b-b989-5ef4720f7ee3",
                    "created_on": "2017-09-04T21:07:10+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "80a4772b-339f-4a4e-aefb-4be1b96f0595",
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
            X-Duration: 88.814µs

    + Body

            {
                "resource_id": "80a4772b-339f-4a4e-aefb-4be1b96f0595"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('80a4772b-339f-4a4e-aefb-4be1b96f0595')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 56.121µs
            X-Resource-Id: 80a4772b-339f-4a4e-aefb-4be1b96f0595

    + Body

            {
                "resource_id": "00000000-0000-0000-0000-000000000000"
            }

