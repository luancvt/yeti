# Changelog

## [0.4.0](https://github.com/luancvt/yeti/compare/v0.3.0...v0.4.0) (2026-06-26)


### Features

* add air live-reload and fetch-models script ([b78794c](https://github.com/luancvt/yeti/commit/b78794c41ae243e798782b2316ba01e3d2c7df75))
* add collection config loader with validation ([71c7174](https://github.com/luancvt/yeti/commit/71c71748754bb35e1981938b6dc6baf8549642d9))
* add collection/model selection UI with shareable URLs ([f7b566b](https://github.com/luancvt/yeti/commit/f7b566b41db53e82d433d511c2c417840bd3bd78))
* add dark theme with OKLCH color scheme ([fb8eaf1](https://github.com/luancvt/yeti/commit/fb8eaf1abcb17d491f08c578bad7854e998c209e))
* add dark-themed scrollbar styling ([5cb9813](https://github.com/luancvt/yeti/commit/5cb981321304945eee3de313c92c203ec6a7bcab))
* add draggable panel divider between tree and detail views ([81516d0](https://github.com/luancvt/yeti/commit/81516d0334b312931c1762685682b30a13c4df20))
* add fade-in transitions for HTMX content swaps ([b1c6599](https://github.com/luancvt/yeti/commit/b1c65995b43ab5387a88300baf492ff1157a4343))
* add fetch-models script for downloading YANG models ([99ba971](https://github.com/luancvt/yeti/commit/99ba971cb4de0c6af0ad2ba4b2864e233c5ec634))
* add GitHub Actions for CI, release-please, and GHCR publishing ([f545286](https://github.com/luancvt/yeti/commit/f5452865bbc775769290e73daa30290f51db81e9))
* add golangci-lint config with fmt/lint justfile commands and fix all lint issues ([43cf24a](https://github.com/luancvt/yeti/commit/43cf24af0df6388368cc2f78044fbe059374c032))
* add health probes and resource limits to Helm chart ([a159615](https://github.com/luancvt/yeti/commit/a1596155366b66c6410b81d11614c23475d180d7))
* add HTTP handlers for tree browsing and detail panel ([ef524a6](https://github.com/luancvt/yeti/commit/ef524a6ba46902bfc1e83c720e5474b2e9ddac1e))
* add Kubernetes deployment support with Helm chart and config-driven collections ([bc149fd](https://github.com/luancvt/yeti/commit/bc149fdbf6b5b547a08830ac542e5d7c1bf4a2ce))
* add labels to collection and model picker controls ([2070080](https://github.com/luancvt/yeti/commit/20700804940b009c33280fcb2be9ff26f17c679b))
* add loading pulse animation on tree node expand ([3c8448d](https://github.com/luancvt/yeti/commit/3c8448d550522856dce209497f81915a52eccc21))
* add new icon ([f1320af](https://github.com/luancvt/yeti/commit/f1320af9f25e66c6922e440d243e8834da321b02))
* add self-hosted Fira Code font for monospace elements ([16fca47](https://github.com/luancvt/yeti/commit/16fca479373e0ba24ea6b97b18be00b449e7c527))
* add styled empty states with icons for tree and detail panels ([5e9c808](https://github.com/luancvt/yeti/commit/5e9c8085b7d4605dfd688be903de007fcee435fc))
* add templ components for layout, tree, and detail panel ([202daa2](https://github.com/luancvt/yeti/commit/202daa214f8b464a7119d7a746ea26e2d72efa49))
* add tree connection lines and row-click expand/collapse with detail loading ([86b972b](https://github.com/luancvt/yeti/commit/86b972ba421f8cc994b3c0d6aae56e5cacd2c1ac))
* add YANG parser that loads collections from fs.FS ([31e79e3](https://github.com/luancvt/yeti/commit/31e79e3cbba4868cfd9109d8a4c8add7da63df97))
* close model picker on click outside via CSS :focus-within ([fffc5ae](https://github.com/luancvt/yeti/commit/fffc5ae42d008e5087b19deefa8eb23e2102ee89))
* generate new templates with icon ([6095a9d](https://github.com/luancvt/yeti/commit/6095a9d93c64538da6808f6e37af45fc1d1c1f0d))
* highlight selected node in tree view ([fd9adf7](https://github.com/luancvt/yeti/commit/fd9adf72a02dfe6bc4ca0d870babf636cb639093))
* redesign header with brand bar and move selectors to tree panel toolbar ([d43afea](https://github.com/luancvt/yeti/commit/d43afeaec1655125229306a3d3b853c6a0edc80c))
* replace CDN assets with embedded local files via go:embed ([94eaf89](https://github.com/luancvt/yeti/commit/94eaf894125fd69c4697492e6cfab928c59ea737))
* return ok body from health endpoint ([7d5dfa8](https://github.com/luancvt/yeti/commit/7d5dfa856ce2ec5281cd65aeaba706ab3b4e0e69))
* surface goyang Extra fields in detail view and tree ([d6a366c](https://github.com/luancvt/yeti/commit/d6a366c53d0c5548f31e213fe67688b2806f2ef9))
* switch to Ginkgo/Gomega, add NodeKind type safety, add tree query tests ([def17d7](https://github.com/luancvt/yeti/commit/def17d734052d7ba444f033175ccefc6078e5a2f))
* use goyang fork with YANG 1.1 must-in-input fix, parse top-level only ([106d302](https://github.com/luancvt/yeti/commit/106d30274495b7031f14539ad4df95b9585b0ade))
* wire up HTTP server for YANG tree browser ([11dd439](https://github.com/luancvt/yeti/commit/11dd439f856d0bc7dfb33e6016a1d2b7475161ce))


### Bug Fixes

* add missing gitignore file ([a980220](https://github.com/luancvt/yeti/commit/a9802208049620acd063008c8f87163977b8d9e5))
* highlight leaf nodes on click by using Alpine scope for selectNode ([dc59251](https://github.com/luancvt/yeti/commit/dc59251a568285fa8d76e897c5027bfdd09fa2ab))
* ignore some parser warnings ([8445398](https://github.com/luancvt/yeti/commit/8445398d524c7de86526cb05ea2474253ccd5e1b))
* improve helm chart defaults and add kind install recipe ([3a5e31e](https://github.com/luancvt/yeti/commit/3a5e31e79380b0dc6d4b81a2340db8cb0a2e1ce4))
* install tailwindcss package before running CLI in CI and Docker ([b7c90cc](https://github.com/luancvt/yeti/commit/b7c90ccb8b0b18009138bdc811f2fce109d3a412))
* just check ([59376f8](https://github.com/luancvt/yeti/commit/59376f8f5470dc740800b91f9e1895fa4785b031))
* merge release jobs into release-please workflow ([4cdf51e](https://github.com/luancvt/yeti/commit/4cdf51e247bffc86145f0b97de6b515f6c63bf9e))
* remove default list bullets from tree panel nodes ([540263a](https://github.com/luancvt/yeti/commit/540263a3ecd697a58c7b43ddc4171fbbbb7015ff))
* remove nonexistent templ --check flag from CI and justfile ([fc6e0ce](https://github.com/luancvt/yeti/commit/fc6e0ce339b0d4ae5beb99da1f693e449ee1770c))
* reset tree and detail panels on collection/model change using HTMX OOB swaps ([71abb11](https://github.com/luancvt/yeti/commit/71abb11c09c6567be7d237b30eab78ba2c9e2f40))
* revert health endpoint body to avoid lint error ([7fff755](https://github.com/luancvt/yeti/commit/7fff7559b448f1808b9abb5a04853a5ae862a48e))
* treat Process() errors as warnings for real-world model compatibility ([7901918](https://github.com/luancvt/yeti/commit/7901918532b4a49ccabbc0f8387d077cef525b0d))
* use go release type and fix golangci-lint version in CI ([51fe712](https://github.com/luancvt/yeti/commit/51fe7126ce06eb15c96a4e789f901ac79933ecab))
* use npx for tailwindcss in CI to resolve [@import](https://github.com/import) ([16fae93](https://github.com/luancvt/yeti/commit/16fae935c9b3d56aff5b6913cac0d14d3f5acf9d))

## [0.3.0](https://github.com/TerjeLafton/yeti/compare/v0.2.0...v0.3.0) (2026-03-03)


### Features

* return ok body from health endpoint ([7d5dfa8](https://github.com/TerjeLafton/yeti/commit/7d5dfa856ce2ec5281cd65aeaba706ab3b4e0e69))


### Bug Fixes

* merge release jobs into release-please workflow ([4cdf51e](https://github.com/TerjeLafton/yeti/commit/4cdf51e247bffc86145f0b97de6b515f6c63bf9e))

## [0.2.0](https://github.com/TerjeLafton/yeti/compare/v0.1.0...v0.2.0) (2026-03-03)


### Features

* add air live-reload and fetch-models script ([b78794c](https://github.com/TerjeLafton/yeti/commit/b78794c41ae243e798782b2316ba01e3d2c7df75))
* add collection config loader with validation ([71c7174](https://github.com/TerjeLafton/yeti/commit/71c71748754bb35e1981938b6dc6baf8549642d9))
* add collection/model selection UI with shareable URLs ([f7b566b](https://github.com/TerjeLafton/yeti/commit/f7b566b41db53e82d433d511c2c417840bd3bd78))
* add dark theme with OKLCH color scheme ([fb8eaf1](https://github.com/TerjeLafton/yeti/commit/fb8eaf1abcb17d491f08c578bad7854e998c209e))
* add dark-themed scrollbar styling ([5cb9813](https://github.com/TerjeLafton/yeti/commit/5cb981321304945eee3de313c92c203ec6a7bcab))
* add draggable panel divider between tree and detail views ([81516d0](https://github.com/TerjeLafton/yeti/commit/81516d0334b312931c1762685682b30a13c4df20))
* add fade-in transitions for HTMX content swaps ([b1c6599](https://github.com/TerjeLafton/yeti/commit/b1c65995b43ab5387a88300baf492ff1157a4343))
* add fetch-models script for downloading YANG models ([99ba971](https://github.com/TerjeLafton/yeti/commit/99ba971cb4de0c6af0ad2ba4b2864e233c5ec634))
* add GitHub Actions for CI, release-please, and GHCR publishing ([f545286](https://github.com/TerjeLafton/yeti/commit/f5452865bbc775769290e73daa30290f51db81e9))
* add golangci-lint config with fmt/lint justfile commands and fix all lint issues ([43cf24a](https://github.com/TerjeLafton/yeti/commit/43cf24af0df6388368cc2f78044fbe059374c032))
* add health probes and resource limits to Helm chart ([a159615](https://github.com/TerjeLafton/yeti/commit/a1596155366b66c6410b81d11614c23475d180d7))
* add HTTP handlers for tree browsing and detail panel ([ef524a6](https://github.com/TerjeLafton/yeti/commit/ef524a6ba46902bfc1e83c720e5474b2e9ddac1e))
* add Kubernetes deployment support with Helm chart and config-driven collections ([bc149fd](https://github.com/TerjeLafton/yeti/commit/bc149fdbf6b5b547a08830ac542e5d7c1bf4a2ce))
* add labels to collection and model picker controls ([2070080](https://github.com/TerjeLafton/yeti/commit/20700804940b009c33280fcb2be9ff26f17c679b))
* add loading pulse animation on tree node expand ([3c8448d](https://github.com/TerjeLafton/yeti/commit/3c8448d550522856dce209497f81915a52eccc21))
* add self-hosted Fira Code font for monospace elements ([16fca47](https://github.com/TerjeLafton/yeti/commit/16fca479373e0ba24ea6b97b18be00b449e7c527))
* add styled empty states with icons for tree and detail panels ([5e9c808](https://github.com/TerjeLafton/yeti/commit/5e9c8085b7d4605dfd688be903de007fcee435fc))
* add templ components for layout, tree, and detail panel ([202daa2](https://github.com/TerjeLafton/yeti/commit/202daa214f8b464a7119d7a746ea26e2d72efa49))
* add tree connection lines and row-click expand/collapse with detail loading ([86b972b](https://github.com/TerjeLafton/yeti/commit/86b972ba421f8cc994b3c0d6aae56e5cacd2c1ac))
* add YANG parser that loads collections from fs.FS ([31e79e3](https://github.com/TerjeLafton/yeti/commit/31e79e3cbba4868cfd9109d8a4c8add7da63df97))
* close model picker on click outside via CSS :focus-within ([fffc5ae](https://github.com/TerjeLafton/yeti/commit/fffc5ae42d008e5087b19deefa8eb23e2102ee89))
* highlight selected node in tree view ([fd9adf7](https://github.com/TerjeLafton/yeti/commit/fd9adf72a02dfe6bc4ca0d870babf636cb639093))
* redesign header with brand bar and move selectors to tree panel toolbar ([d43afea](https://github.com/TerjeLafton/yeti/commit/d43afeaec1655125229306a3d3b853c6a0edc80c))
* replace CDN assets with embedded local files via go:embed ([94eaf89](https://github.com/TerjeLafton/yeti/commit/94eaf894125fd69c4697492e6cfab928c59ea737))
* surface goyang Extra fields in detail view and tree ([d6a366c](https://github.com/TerjeLafton/yeti/commit/d6a366c53d0c5548f31e213fe67688b2806f2ef9))
* switch to Ginkgo/Gomega, add NodeKind type safety, add tree query tests ([def17d7](https://github.com/TerjeLafton/yeti/commit/def17d734052d7ba444f033175ccefc6078e5a2f))
* use goyang fork with YANG 1.1 must-in-input fix, parse top-level only ([106d302](https://github.com/TerjeLafton/yeti/commit/106d30274495b7031f14539ad4df95b9585b0ade))
* wire up HTTP server for YANG tree browser ([11dd439](https://github.com/TerjeLafton/yeti/commit/11dd439f856d0bc7dfb33e6016a1d2b7475161ce))


### Bug Fixes

* add missing gitignore file ([a980220](https://github.com/TerjeLafton/yeti/commit/a9802208049620acd063008c8f87163977b8d9e5))
* highlight leaf nodes on click by using Alpine scope for selectNode ([dc59251](https://github.com/TerjeLafton/yeti/commit/dc59251a568285fa8d76e897c5027bfdd09fa2ab))
* ignore some parser warnings ([8445398](https://github.com/TerjeLafton/yeti/commit/8445398d524c7de86526cb05ea2474253ccd5e1b))
* improve helm chart defaults and add kind install recipe ([3a5e31e](https://github.com/TerjeLafton/yeti/commit/3a5e31e79380b0dc6d4b81a2340db8cb0a2e1ce4))
* install tailwindcss package before running CLI in CI and Docker ([b7c90cc](https://github.com/TerjeLafton/yeti/commit/b7c90ccb8b0b18009138bdc811f2fce109d3a412))
* remove default list bullets from tree panel nodes ([540263a](https://github.com/TerjeLafton/yeti/commit/540263a3ecd697a58c7b43ddc4171fbbbb7015ff))
* remove nonexistent templ --check flag from CI and justfile ([fc6e0ce](https://github.com/TerjeLafton/yeti/commit/fc6e0ce339b0d4ae5beb99da1f693e449ee1770c))
* reset tree and detail panels on collection/model change using HTMX OOB swaps ([71abb11](https://github.com/TerjeLafton/yeti/commit/71abb11c09c6567be7d237b30eab78ba2c9e2f40))
* treat Process() errors as warnings for real-world model compatibility ([7901918](https://github.com/TerjeLafton/yeti/commit/7901918532b4a49ccabbc0f8387d077cef525b0d))
* use go release type and fix golangci-lint version in CI ([51fe712](https://github.com/TerjeLafton/yeti/commit/51fe7126ce06eb15c96a4e789f901ac79933ecab))
* use npx for tailwindcss in CI to resolve [@import](https://github.com/import) ([16fae93](https://github.com/TerjeLafton/yeti/commit/16fae935c9b3d56aff5b6913cac0d14d3f5acf9d))
