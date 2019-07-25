# hello-world

Run tests to make sure everything is looking good:

```
make test
```

Issue an rpc by first running the service:

```
make run
```

And then in another window, run the following:

```
yab -y idl/code.uber.internal/wonsoh/hello-world/hello.yab -A name:Bob
```

Cheers - you have a working Go service! To learn more, read:

* [t.uber.com/fx](http://t.uber.com/fx).
* [t.uber.com/yarpc](http://t.uber.com/yarpc).

## What now?

Some common steps service owners take after scaffolding are:

1. Creating a [new Jenkins job](https://stack.uberinternal.com/questions/1521/how-to-setup-a-new-jenkins-build).
2. Provisioning a new service in [Infraportal](https://wonsoh.uberinternal.com/docker/initiate/).
3. Deploying the service using [uDeploy](https://udeploy.uberinternal.com/).
4. Onboarding a new service to [Muttley](http://t.uber.com/muttley-onboarding).
5. Replacing this README with your own content.

After the service has been deployed, you can verify it in production by:

```
ssh adhoc10-sjc1
uns://sjc1/sjc1-prod01/us1/hello-world:http
git clone gitolite@code.uber.internal:wonsoh/hello-world && cd hello-world
yab -y ./idl/code.uber.internal/wonsoh/hello-world/hello.yab -P uns://sjc1/sjc1-prod01/us1/hello-world:http -A name:Bob
```
