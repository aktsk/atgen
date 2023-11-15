# Changelog

## [v0.8.0](https://github.com/aktsk/atgen/compare/v0.7.0...v0.8.0) (2023-11-15)

* exampleに```/v1/people/all```を追加する。 [#55](https://github.com/aktsk/atgen/pull/55) ([seipan](https://github.com/seipan))
* 配列型のjsonのレスポンスチェックに対応させる [#53](https://github.com/aktsk/atgen/pull/53) ([seipan](https://github.com/seipan))
* Bump golang.org/x/text from 0.3.2 to 0.3.8 [#31](https://github.com/aktsk/atgen/pull/31) ([dependabot[bot]](https://github.com/apps/dependabot))
* Update module golang.org/x/tools to v0.13.0 [#40](https://github.com/aktsk/atgen/pull/40) ([renovate[bot]](https://github.com/apps/renovate))
* Update module github.com/spf13/afero to v1.9.5 [#37](https://github.com/aktsk/atgen/pull/37) ([renovate[bot]](https://github.com/apps/renovate))
* 依存関係の更新 [#48](https://github.com/aktsk/atgen/pull/48) ([takanakahiko](https://github.com/takanakahiko))
* Update module github.com/urfave/cli to v2 [#44](https://github.com/aktsk/atgen/pull/44) ([renovate[bot]](https://github.com/apps/renovate))
* Update module github.com/urfave/cli to v1.22.14 [#39](https://github.com/aktsk/atgen/pull/39) ([renovate[bot]](https://github.com/apps/renovate))
* Update actions/checkout action to v4 [#42](https://github.com/aktsk/atgen/pull/42) ([renovate[bot]](https://github.com/apps/renovate))
* Update module gopkg.in/yaml.v2 to v3 [#45](https://github.com/aktsk/atgen/pull/45) ([renovate[bot]](https://github.com/apps/renovate))
* Update actions/setup-go action to v4 [#43](https://github.com/aktsk/atgen/pull/43) ([renovate[bot]](https://github.com/apps/renovate))
* Update module github.com/gorilla/mux to v1.8.0 [#36](https://github.com/aktsk/atgen/pull/36) ([renovate[bot]](https://github.com/apps/renovate))
* Configure Renovate [#35](https://github.com/aktsk/atgen/pull/35) ([renovate[bot]](https://github.com/apps/renovate))
* Bump gopkg.in/yaml.v2 from 2.2.2 to 2.2.8 [#30](https://github.com/aktsk/atgen/pull/30) ([dependabot[bot]](https://github.com/apps/dependabot))
* Update go modules [#29](https://github.com/aktsk/atgen/pull/29) ([mizzy](https://github.com/mizzy))

## [v0.7.0](https://github.com/aktsk/atgen/compare/v0.6.0...v0.7.0) (2020-10-30)

* Access nested data with register [#27](https://github.com/aktsk/atgen/pull/27) ([takanakahiko](https://github.com/takanakahiko))
* Test on GitHub Actions [#26](https://github.com/aktsk/atgen/pull/26) ([takanakahiko](https://github.com/takanakahiko))
* Support form and raw type request body [#22](https://github.com/aktsk/atgen/pull/22) ([p1ass](https://github.com/p1ass))
* Fix variable name from vars to atgenVars [#24](https://github.com/aktsk/atgen/pull/24) ([p1ass](https://github.com/p1ass))
* Change to use format.Node when writing generated code [#23](https://github.com/aktsk/atgen/pull/23) ([p1ass](https://github.com/p1ass))
* Generate http request body in Atgen [#21](https://github.com/aktsk/atgen/pull/21) ([p1ass](https://github.com/p1ass))
* Fix to replace status to atgenStatus for placeholder [#19](https://github.com/aktsk/atgen/pull/19) ([p1ass](https://github.com/p1ass))
* Add Atgen prefix before functions and variables replaced by atgen [#18](https://github.com/aktsk/atgen/pull/18) ([p1ass](https://github.com/p1ass))

## [v0.6.0](https://github.com/aktsk/atgen/compare/v0.5.0...v0.6.0) (2019-05-13)

* Cache results of packges.Load() to shorten execution time [#14](https://github.com/aktsk/atgen/pull/14) ([mizzy](https://github.com/mizzy))
* Remove RouterFuncs of Generator struct [#13](https://github.com/aktsk/atgen/pull/13) ([mizzy](https://github.com/mizzy))

## [v0.5.0](https://github.com/aktsk/atgen/compare/v0.4.0...v0.5.0) (2019-04-08)

* Replace TrimLeft and TrimRight with TrimPrefix and TrimSuffix [#12](https://github.com/aktsk/atgen/pull/12) ([mizzy](https://github.com/mizzy))

## [v0.4.0](https://github.com/aktsk/atgen/compare/v0.3.0...v0.4.0) (2019-04-01)

* Use go/packages instead of go/loader [#11](https://github.com/aktsk/atgen/pull/11) ([mizzy](https://github.com/mizzy))
* Disable golint temporary [#10](https://github.com/aktsk/atgen/pull/10) ([mizzy](https://github.com/mizzy))

## [v0.3.0](https://github.com/aktsk/atgen/compare/v0.2.0...v0.3.0) (2019-03-12)

* Support register variables [#8](https://github.com/aktsk/atgen/pull/8) ([mizzy](https://github.com/mizzy))
* Implement safe copy [#6](https://github.com/aktsk/atgen/pull/6) ([objectx](https://github.com/objectx))

## [v0.2.0](https://github.com/aktsk/atgen/compare/v0.1.0...v0.2.0) (2019-01-29)

* Inject router function [#5](https://github.com/aktsk/atgen/pull/5) ([sachaos](https://github.com/sachaos))
