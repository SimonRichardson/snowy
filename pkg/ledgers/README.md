# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2017-08-23T09:35:16+01:00
# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('b78d6ddc-bfee-4694-9f93-7a2040e12baa')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 73.968µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: b78d6ddc-bfee-4694-9f93-7a2040e12baa

    + Body

            {
                "author_id": "346f58b2-60c4-42b6-8e55-0acb9262cca8",
                "created_on": "2017-08-23T09:35:16+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "b78d6ddc-bfee-4694-9f93-7a2040e12baa",
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
            resource_id ('b78d6ddc-bfee-4694-9f93-7a2040e12baa')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 62.655µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: b78d6ddc-bfee-4694-9f93-7a2040e12baa

    + Body

            [
                {
                    "author_id": "346f58b2-60c4-42b6-8e55-0acb9262cca8",
                    "created_on": "2017-08-23T09:35:16+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "b78d6ddc-bfee-4694-9f93-7a2040e12baa",
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
            X-Duration: 142.74µs

    + Body

            {
                "resource_id": "b78d6ddc-bfee-4694-9f93-7a2040e12baa"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('b78d6ddc-bfee-4694-9f93-7a2040e12baa')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 66.582µs
            X-Resource-Id: b78d6ddc-bfee-4694-9f93-7a2040e12baa

    + Body

            {
                "resource_id": "00000000-0000-0000-0000-000000000000"
            }

