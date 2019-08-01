#!/bin/bash

yab ./idl/code.uber.internal/wonsoh/hello-world/get.yab -A id:$1 | jq .
