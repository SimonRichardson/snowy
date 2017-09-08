# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2017-09-08T10:40:20+01:00
# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('d7030131-03cc-4eb4-8d7c-ec2ebde846e5')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 77.145µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: d7030131-03cc-4eb4-8d7c-ec2ebde846e5

    + Body

            {
                "author_id": "10fcbc8c-7054-41bc-8329-aa31b2f1caed",
                "created_on": "2017-09-08T10:40:20+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "parent_id": "00000000-0000-0000-0000-000000000000",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "d7030131-03cc-4eb4-8d7c-ec2ebde846e5",
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
            resource_id ('d7030131-03cc-4eb4-8d7c-ec2ebde846e5')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 44.301µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: d7030131-03cc-4eb4-8d7c-ec2ebde846e5

    + Body

            [
                {
                    "author_id": "10fcbc8c-7054-41bc-8329-aa31b2f1caed",
                    "created_on": "2017-09-08T10:40:20+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "parent_id": "00000000-0000-0000-0000-000000000000",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "d7030131-03cc-4eb4-8d7c-ec2ebde846e5",
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
            X-Duration: 107.36µs

    + Body

            {
                "resource_id": "d7030131-03cc-4eb4-8d7c-ec2ebde846e5"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('d7030131-03cc-4eb4-8d7c-ec2ebde846e5')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 77.975µs
            X-Resource-Id: d7030131-03cc-4eb4-8d7c-ec2ebde846e5

    + Body

            {
                "resource_id": "00000000-0000-0000-0000-000000000000"
            }

