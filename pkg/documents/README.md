# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2017-08-10T12:47:57+01:00
# GET /multiple

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('6085672b-f109-4f97-8183-19e49926f601')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 41.591µs
            X-Query-Tags: abc,def,g
            X-Resourceid: 6085672b-f109-4f97-8183-19e49926f601

    + Body

            [
                {
                    "author_id": "9e44a3d5-f01e-4337-bec7-1576e7afc7f4",
                    "created_on": "2017-08-10T12:47:57+01:00",
                    "deleted_on": "0001-01-01T00:00:00Z",
                    "name": "document-name",
                    "resource_id": "6085672b-f109-4f97-8183-19e49926f601",
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
            X-Duration: 66.926µs

    + Body

            {
                "resource_id": "6085672b-f109-4f97-8183-19e49926f601"
            }

# PUT /

+ Request
    + Parameters

            resource_id ('6085672b-f109-4f97-8183-19e49926f601')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 100
            Content-Type: application/json
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 50.303µs
            X-Resourceid: 6085672b-f109-4f97-8183-19e49926f601

    + Body

            {
                "resource_id": "00000000-0000-0000-0000-000000000000"
            }

# GET /

+ Request
    + Parameters

            query.tags ('abc,def,g')
            resource_id ('6085672b-f109-4f97-8183-19e49926f601')

    + Headers

            Accept-Encoding: gzip
            User-Agent: Go-http-client/1.1

+ Response 200
    + Headers

            Content-Type: application/json
            X-Duration: 69.048µs
            X-Query-Tags: abc,def,g
            X-Resourceid: 6085672b-f109-4f97-8183-19e49926f601

    + Body

            {
                "author_id": "9e44a3d5-f01e-4337-bec7-1576e7afc7f4",
                "created_on": "2017-08-10T12:47:57+01:00",
                "deleted_on": "0001-01-01T00:00:00Z",
                "name": "document-name",
                "resource_id": "6085672b-f109-4f97-8183-19e49926f601",
                "tags": [
                    "abc",
                    "def",
                    "g"
                ]
            }

