# Yeti - YANG Entry Tree Inspector

A web-based tool for browsing and inspecting YANG models.  
Yeti parses YANG modules and presents them as interactive, navigable trees, making it easy to explore container hierarchies and leaf types across different model collections (e.g. Cisco IOS-XR releases).

## Helm Chart

The chart is in `charts/yeti/`. It includes an init container that fetches YANG models from [YangModels/yang](https://github.com/YangModels/yang) on startup.

### Install

```bash
helm install yeti oci://ghcr.io/terjelafton/charts/yeti
```

### Configuration

#### Collections

Each collection maps to a path in the [YangModels/yang](https://github.com/YangModels/yang) repository. The init container sparse-clones only the specified paths.

```yaml
collections:
  - name: xr-7112
    display: "Cisco IOS-XR 7.11.2"
    path: vendor/cisco/xr/7112
```

#### Resources

Parsing YANG models is memory-intensive. The resource requirements per collection depend on the number of models. As a reference, Cisco IOS-XR 7.11.2 and 24.4.2 each use between 1.5 and 2.0 GB of memory. Adjust based on the number of collections you deploy.

```yaml
resources:
  requests:
    cpu: 100m
    memory: 2Gi
  limits:
    cpu: "1"
    memory: 4Gi
```

## Local Development

Prerequisites: Go 1.25+, [just](https://github.com/casey/just), [templ](https://templ.guide), [Tailwind CSS CLI](https://tailwindcss.com/blog/standalone-cli), [air](https://github.com/air-verse/air)

Fetch YANG models and start the dev server with live reload:

```bash
just fetch-models xr-7112 vendor/cisco/xr/7112
just dev
```

Open http://localhost:8080.

Other useful commands:

```bash
just test    # run tests
just lint    # run linter
just fmt     # format code
just check   # run all CI checks
```

## Kind (Local Testing)

With a Kind cluster named `yeti`:

```bash
just install
```

This builds the Docker image, loads it into Kind, and installs the chart in the `yeti` namespace. To use a different namespace:

```bash
just install other-namespace
```
