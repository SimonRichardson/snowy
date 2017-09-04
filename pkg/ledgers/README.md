# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2017-09-04T20:36:42+01:00
# GET /multiple/

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('196df4a7-fb6f-4e41-b2c9-5c2f2cbf558d')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 45.424µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 196df4a7-fb6f-4e41-b2c9-5c2f2cbf558d

    + Body

            [
                {
                    "author_id": "40271c71-c061-4a4c-bd24-7023e70e7d34",
                    "created_on": "2017-09-04T20:36:42+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_address": "abcdefghij",
                    "resource_content_type": "application/octet-stream",
                    "resource_id": "196df4a7-fb6f-4e41-b2c9-5c2f2cbf558d",
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
            X-Duration: 131.571µs

    + Body

            {
                "resource_id": "196df4a7-fb6f-4e41-b2c9-5c2f2cbf558d"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('196df4a7-fb6f-4e41-b2c9-5c2f2cbf558d')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 71.303µs
            X-Resource-Id: 196df4a7-fb6f-4e41-b2c9-5c2f2cbf558d

    + Body

            {
                "resource_id": "00000000-0000-0000-0000-000000000000"
            }

# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('196df4a7-fb6f-4e41-b2c9-5c2f2cbf558d')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 77.478µs
            X-Query-Author-Id: 
            X-Query-Tags: abc,def,g
            X-Resource-Id: 196df4a7-fb6f-4e41-b2c9-5c2f2cbf558d

    + Body

            {
                "author_id": "40271c71-c061-4a4c-bd24-7023e70e7d34",
                "created_on": "2017-09-04T20:36:42+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_address": "abcdefghij",
                "resource_content_type": "application/octet-stream",
                "resource_id": "196df4a7-fb6f-4e41-b2c9-5c2f2cbf558d",
                "resource_size": 10,
                "tags": [
                    "abc",
                    "def",
                    "g"
                ]
            }

