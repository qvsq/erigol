# How to build

To build this you need go whatever version. First let's download dependencies:

```bash
make deps
```

Then build binaries:

```bash
make client
make server
```

# How to run

You need to run rabbitmq. The easest way do it via docker:

```bash
docker run -d -p 5672:5672 --hostname my-rabbit --name some-rabbit rabbitmq:3
```

This will expose 5672 port.

Now you are ready to run server:

```bash
bin/server
```

To see what are the paremeters you could use help flag for both client and server, for example:

```bash
bin/client --help
```

# Some QA

## Why not logrus for logging ?

Yes, it is probably better but I think standart log here is enough.

## Why you used another library in your order map implementation ?

I used container/list which is double linked list. I can explain how it is works and how ordered map works.

## Is order map safe for parallel usage ?

Not yet, probably mutex is required.
