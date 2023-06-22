#!/bin/bash
curl -XPOST http://localhost:8080/user/upload -H "Content-Type: multipart/form-data" -F "file=@hello.zip"
