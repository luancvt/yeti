# Vendored goyang

Vendored from [hansthienpondt/goyang](https://github.com/hansthienpondt/goyang)
`v1.6.0-fixrpc` (a fork of [openconfig/goyang](https://github.com/openconfig/goyang)),
with one patch:

- `Uses.Augment` is `[]*Augment` instead of `*Augment`. RFC 7950 allows a
  `uses` block to contain multiple `augment` statements; upstream goyang
  rejects the second one with `augment: already set` and drops the whole
  module. This breaks modules that import `ietf-keystore` (RFC 9642 uses
  two augments in one `uses`), which IOS-XR vendor trees ship.

Only `pkg/yang` and `pkg/indent` are vendored, without tests. The module path
stays `github.com/openconfig/goyang` because the main module replaces that
import path with this directory.

To pick up upstream changes, copy the new `pkg/yang` and `pkg/indent`
sources over these, re-apply the `Uses.Augment` patch, and drop the test
files again.
