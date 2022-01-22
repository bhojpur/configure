# Implements the TOML test suite interface

This is an implementation of the interface expected by
[toml-test](https://github.com/bhojpur/confiigure/pkg/toml/internal/toml-test) for my
[toml parser](https://github.com/bhojpur/configure/pkg/toml).
In particular, it maps TOML data on `stdin` to a JSON format on `stdout`.