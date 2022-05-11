# kubectl-suspender

A `kubectl` plugin to (un)suspend namespaces managed by [kube-ns-suspender](https://github.com/govirtuo/kube-ns-suspender).

- [kubectl-suspender](#kubectl-suspender)
  - [Getting started](#getting-started)
  - [Usage](#usage)
    - [`list`](#list)
    - [`suspend`](#suspend)
    - [`unsuspend`](#unsuspend)
  - [Known limitations](#known-limitations)
  - [License](#license)

## Getting started

Grab the latest release [here](https://github.com/govirtuo/kubectl-suspender/releases) or build it yourself:

```
git clone https://github.com/govirtuo/kubectl-suspender.git
cd kubectl-suspender/
make build
```

Once you have a working binary, you have to put it somewhere in your path, and `kubectl` will find it automatically.

## Usage

There is three main commands: `list`, `suspend` and `unsuspend`.

### `list`

This command simply lists the namespaces managed by kube-ns-suspender and displays their states.

Usage example:

```
kubectl suspender list
```

### `suspend`

The suspend command allows a user to suspend a namespace managed by kube-ns-suspender. If no namespace is provided in arg, it will prompt a list of namespaces to the user.

Usage examples:

```
kubectl suspender suspend
kubectl suspender suspend my-namespace
kubectl suspender suspend my-namespace my-other-namespace
```

### `unsuspend`

The unsuspend command allows a user to unsuspend a namespace managed by kube-ns-suspender. If no namespace is provided in arg, it will prompt a list of namespaces to the user.

Usage example:

```
kubectl suspender unsuspend
kubectl suspender unsuspend my-namespace
kubectl suspender unsuspend my-namespace my-other-namespace
```

## Known limitations

This project is in its early days, and a few things are missing:

* it is impossible to set the context using a flag (i.e. `--context staging-cluster`);
* the completion only works if you use the full binary name (`kubectl-suspender` instead of `kubectl suspender`);
* the [controller name](https://github.com/govirtuo/kube-ns-suspender#controllername) cannot be set using flag, it uses `kube-ns-suspender` as a default. This should be easily fixable.

## License

[MIT](./LICENSE)