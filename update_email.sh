#!/bin/bash

yab ./idl/code.uber.internal/wonsoh/hello-world/update_email.yab -A id:$1 -A email:$2 | jq .
