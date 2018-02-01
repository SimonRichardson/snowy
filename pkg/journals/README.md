# Snowy

The following was automatically generated via [Betwixt](https://github.com/simonrichardson/betwixt).
Date generated on: 2018-02-01T16:30:35Z
# POST /

+ Request
    + Parameters


    + Headers

            Accept-Encoding: gzip
            Content-Length: 730
            Content-Type: multipart/form-data; boundary=d5397af16385a891ebff3707057852692c75f78e0bb5ec831b4d2e09d073
            User-Agent: Go-http-client/1.1

    + Body

            --d5397af16385a891ebff3707057852692c75f78e0bb5ec831b4d2e09d073
Content-Disposition: form-data; name="content"; filename="content"
Content-Length: 176
Content-Type: application/octet-stream

TxY_Xw-aYh1ylWbHTRADfE17uwQH0eLGSYGFWthoHQ2G0ekeABZ5OctmlNLEIqzSCKAHKTlIf2mZ650YpEeEBF2H88Z88idG6ZWvWiU2eVG6ov9s1HHEg_FfuQuts3xYIbbZVSakGpUEaAtOfIt2OhsdSdSVXISGIWMlJT_sc43XqeI=
--d5397af16385a891ebff3707057852692c75f78e0bb5ec831b4d2e09d073
Content-Disposition: form-data; name="document"; filename="document"
Content-Length: 100
Content-Type: application/json

{"name":"document-name","author_id":"46b917ba-997b-4782-8314-90bb96100581","tags":["abc","def","g"]}
--d5397af16385a891ebff3707057852692c75f78e0bb5ec831b4d2e09d073--


+ Response 200
    + Headers

            Content-Type: multipart/form-data
            X-Duration: 164.352µs

    + Body

            {"resource_id":"46b917ba-997b-4782-8314-90bb96100581"}


# PUT /

+ Request
    + Parameters

            resource_id ('46b917ba-997b-4782-8314-90bb96100581')

    + Headers

            Accept-Encoding: gzip
            Content-Length: 730
            Content-Type: multipart/form-data; boundary=0e9e5617716925588db2c51f42e7753f8e3c9b8ad2af218b8f4d9d9610f3
            User-Agent: Go-http-client/1.1

    + Body

            --0e9e5617716925588db2c51f42e7753f8e3c9b8ad2af218b8f4d9d9610f3
Content-Disposition: form-data; name="content"; filename="content"
Content-Length: 176
Content-Type: application/octet-stream

TxY_Xw-aYh1ylWbHTRADfE17uwQH0eLGSYGFWthoHQ2G0ekeABZ5OctmlNLEIqzSCKAHKTlIf2mZ650YpEeEBF2H88Z88idG6ZWvWiU2eVG6ov9s1HHEg_FfuQuts3xYIbbZVSakGpUEaAtOfIt2OhsdSdSVXISGIWMlJT_sc43XqeI=
--0e9e5617716925588db2c51f42e7753f8e3c9b8ad2af218b8f4d9d9610f3
Content-Disposition: form-data; name="document"; filename="document"
Content-Length: 100
Content-Type: application/json

{"name":"document-name","author_id":"46b917ba-997b-4782-8314-90bb96100581","tags":["abc","def","g"]}
--0e9e5617716925588db2c51f42e7753f8e3c9b8ad2af218b8f4d9d9610f3--


+ Response 200
    + Headers

            Content-Type: multipart/form-data
            X-Duration: 100.359µs
            X-Resource-Id: 46b917ba-997b-4782-8314-90bb96100581

    + Body

            {"resource_id":"46b917ba-997b-4782-8314-90bb96100581"}


