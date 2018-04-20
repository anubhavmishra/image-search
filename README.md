# image-search
An example microservice that fetches gifs from giphy.

## Usage

Set environment variables

```bash
export GIPHY_API_KEY="API_KEY_HERE"
```

*Optionally* set feature flag

```bash
export APRIL_FOOLS="true"
```

Run the application

```bash
make run
```

Build project

```bash
make build
```

Build the docker image and push it

```bash
make build-docker-image
```
