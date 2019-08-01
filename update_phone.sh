#!/bin/bash

yab ./idl/code.uber.internal/wonsoh/hello-world/update_phone.yab -A id:$1 -A phone:$2 | jq .
