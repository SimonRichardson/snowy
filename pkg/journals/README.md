# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2018-05-10T20:00:35+01:00
# PUT /

+ Request
    + Parameters

            resource_id ('4ddff772-c511-4ad3-8b20-a1f229f8d9f4')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 730
            Content-Type: multipart/form-data; boundary=a660f2475d1e949ada3f4d587ea3464118c635628b10715976cb6112c09b
            User-Agent: Go-http-client/1.1

    + Body

            --a660f2475d1e949ada3f4d587ea3464118c635628b10715976cb6112c09b
Content-Disposition: form-data; name="content"; filename="content"
Content-Length: 176
Content-Type: application/octet-stream

TxY_Xw-aYh1ylWbHTRADfE17uwQH0eLGSYGFWthoHQ2G0ekeABZ5OctmlNLEIqzSCKAHKTlIf2mZ650YpEeEBF2H88Z88idG6ZWvWiU2eVG6ov9s1HHEg_FfuQuts3xYIbbZVSakGpUEaAtOfIt2OhsdSdSVXISGIWMlJT_sc43XqeI=
--a660f2475d1e949ada3f4d587ea3464118c635628b10715976cb6112c09b
Content-Disposition: form-data; name="document"; filename="document"
Content-Length: 100
Content-Type: application/json

{"name":"document-name","author_id":"4ddff772-c511-4ad3-8b20-a1f229f8d9f4","tags":["abc","def","g"]}
--a660f2475d1e949ada3f4d587ea3464118c635628b10715976cb6112c09b--


+ Response 200
    + Headers

            Content-Type: multipart/form-data
            X-Duration: 214.829µs
            X-Resource-Id: 4ddff772-c511-4ad3-8b20-a1f229f8d9f4

    + Body

            {"resource_id":"4ddff772-c511-4ad3-8b20-a1f229f8d9f4"}


# POST /

+ Request
    + Parameters


    + Headers

            Accept-Encoding: gzip
            Content-Length: 730
            Content-Type: multipart/form-data; boundary=572f47036a62a555332b0f47e65152a6f636a48be1353d2614fcf43e4049
            User-Agent: Go-http-client/1.1

    + Body

            --572f47036a62a555332b0f47e65152a6f636a48be1353d2614fcf43e4049
Content-Disposition: form-data; name="content"; filename="content"
Content-Length: 176
Content-Type: application/octet-stream

TxY_Xw-aYh1ylWbHTRADfE17uwQH0eLGSYGFWthoHQ2G0ekeABZ5OctmlNLEIqzSCKAHKTlIf2mZ650YpEeEBF2H88Z88idG6ZWvWiU2eVG6ov9s1HHEg_FfuQuts3xYIbbZVSakGpUEaAtOfIt2OhsdSdSVXISGIWMlJT_sc43XqeI=
--572f47036a62a555332b0f47e65152a6f636a48be1353d2614fcf43e4049
Content-Disposition: form-data; name="document"; filename="document"
Content-Length: 100
Content-Type: application/json

{"name":"document-name","author_id":"4ddff772-c511-4ad3-8b20-a1f229f8d9f4","tags":["abc","def","g"]}
--572f47036a62a555332b0f47e65152a6f636a48be1353d2614fcf43e4049--


+ Response 200
    + Headers

            Content-Type: multipart/form-data
            X-Duration: 257.09µs

    + Body

            {"resource_id":"4ddff772-c511-4ad3-8b20-a1f229f8d9f4"}


