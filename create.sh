#!/bin/bash

yab ./idl/code.uber.internal/wonsoh/hello-world/create.yab -A name:$1 | jq .
