# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2018-02-01T15:32:06Z
# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('25a94268-02f9-4274-8ab9-59baf16d0adb')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 71.532µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 25a94268-02f9-4274-8ab9-59baf16d0adb

    + Body

            {
                "author_id": "56972471-7b12-4984-a27f-338baf3e7362",
                "created_on": "2018-02-01T15:32:06Z",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "25a94268-02f9-4274-8ab9-59baf16d0adb",
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
            resource_id ('25a94268-02f9-4274-8ab9-59baf16d0adb')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 26.831µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 25a94268-02f9-4274-8ab9-59baf16d0adb

    + Body

            [
                {
                    "author_id": "56972471-7b12-4984-a27f-338baf3e7362",
                    "created_on": "2018-02-01T15:32:06Z",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "25a94268-02f9-4274-8ab9-59baf16d0adb",
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
            X-Duration: 55.469µs

    + Body

            {
                "resource_id": "25a94268-02f9-4274-8ab9-59baf16d0adb"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('25a94268-02f9-4274-8ab9-59baf16d0adb')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 39.242µs
            X-Resource-Id: 25a94268-02f9-4274-8ab9-59baf16d0adb

    + Body

            {
                "resource_id": "25a94268-02f9-4274-8ab9-59baf16d0adb"
            }

# PUT /fork/

+ Request
    + Parameters

            resource_id ('25a94268-02f9-4274-8ab9-59baf16d0adb')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 53.307µs
            X-Resource-Id: 25a94268-02f9-4274-8ab9-59baf16d0adb

    + Body

            {
                "resource_id": "8f166bc0-f481-4192-a94d-ddd738aff09d"
            }

