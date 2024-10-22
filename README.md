# 🤫 `secrecy` - A simple secret-keeping library for Go.

<p>
	<a href="https://pkg.go.dev/github.com/negrel/secrecy">
		<img alt="PkgGoDev" src="https://pkg.go.dev/badge/github.com/negrel/secrecy">
	</a>
	<a href="https://goreportcard.com/report/github.com/negrel/secrecy">
		<img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/negrel/secrecy">
	</a>
	<img alt="Go version" src="https://img.shields.io/github/go-mod/go-version/prismelabs/analytics">
</p>

`secrecy` is a simple library which provides wrapper type for secret management
in Go. It is inspired from the excellent
[`secrecy`](https://github.com/iqlusioninc/crates/tree/main/secrecy) Rust crate.

It provides a `Secret[T]` type for wrapping another value in a secret cell which
attempts to limit exposure (only available via the special `ExposeSecret()`
function).

Each secret has a [`finalizer`](https://pkg.go.dev/runtime#SetFinalizer) attached
that recusively zeroize secret memory when garbage collected. You must not share
memory contained within a secret and share the secret itself.

This helps to ensure secrets aren't accidentally copied, logged, or otherwise
exposed (as much as possible), and also ensures secrets are securely wiped from
memory when garbage collected.

## Contributing

If you want to contribute to `secrecy` to add a feature or improve the code contact
me at [alexandre@negrel.dev](mailto:alexandre@negrel.dev), open an
[issue](https://github.com/negrel/secrecy/issues) or make a
[pull request](https://github.com/negrel/secrecy/pulls).

## :stars: Show your support

Please give a :star: if this project helped you!

[![buy me a coffee](https://github.com/negrel/.github/blob/master/.github/images/bmc-button.png?raw=true)](https://www.buymeacoffee.com/negrel)

## :scroll: License

MIT © [Alexandre Negrel](https://www.negrel.dev/)
