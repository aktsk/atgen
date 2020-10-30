# Changelog

## [v0.7.0](https://github.com/aktsk/atgen/compare/v0.6.0...v0.7.0) (2020-10-30)

* [Feature] Add Atgen prefix before functions and variables replaced by atgen #18 #19 #24 ([p1ass](https://github.com/p1ass))
* [Feature] Generate http request body in Atgen #21 ([p1ass](https://github.com/p1ass))
* [Feature] Access nested data with register #27 ([takanakahiko](https://github.com/takanakahiko))
* [Internal] Change to use format.Node when writing generated code #23 ([p1ass](https://github.com/p1ass))
* [Internal] Test on GitHub Actions #26 ([takanakahiko](https://github.com/takanakahiko))

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
