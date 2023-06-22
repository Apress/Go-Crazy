#!/bin/bash
curl -XPOST http://localhost:8080/user/post -H "Content-Type: application/json"  -d '{"name":"nico","email":"hello.nico@gocrazy.com"}'
