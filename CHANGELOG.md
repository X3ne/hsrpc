# Changelog

## [1.1.0](https://github.com/X3ne/hsrpc/compare/v1.0.3...v1.1.0) (2024-02-19)


### Features

* added update checking & a modal to automatically download and update the .exe file ([8a42d43](https://github.com/X3ne/hsrpc/commit/8a42d43fc6de2def74c19269613d573da0737e4d))

## [1.0.3](https://github.com/X3ne/hsrpc/compare/v1.0.2...v1.0.3) (2024-02-18)


### Bug Fixes

* fixed a bug with static file import names (only when the application is built on linux) ([802dd00](https://github.com/X3ne/hsrpc/commit/802dd00213a7381d3e8302622bf5b3f38e085794))


### Miscellaneous Chores

* release 1.0.3 ([b2b5916](https://github.com/X3ne/hsrpc/commit/b2b5916cd3e7708399d70b8863d6e09ee4e4c5f8))

## [1.0.2](https://github.com/X3ne/hsrpc/compare/v1.0.1...v1.0.2) (2024-02-18)


### Bug Fixes

* moved go-winres to go:generate ([af29c3b](https://github.com/X3ne/hsrpc/commit/af29c3ba2fb96b809508dc72e3dbbb2373fe52e4))


### Miscellaneous Chores

* release 1.0.2 ([f5bfe8a](https://github.com/X3ne/hsrpc/commit/f5bfe8a48e1d8fa4d1e657c914948c787e967ec3))

## [1.0.1](https://github.com/X3ne/hsrpc/compare/v1.0.0...v1.0.1) (2024-02-18)


### Bug Fixes

* fixed an issue with the build workflow ([7fe1802](https://github.com/X3ne/hsrpc/commit/7fe1802366a804515d52bbf528005a9583fe1f5b))


### Miscellaneous Chores

* release 1.0.1 ([6bd5c64](https://github.com/X3ne/hsrpc/commit/6bd5c647fcd363e23cd4b64f9d05c41baeeb91c6))

## 1.0.0 (2024-02-18)


### Features

* add "Real-Time Combat View" menu ([506bfbf](https://github.com/X3ne/hsrpc/commit/506bfbf7372b92bd36134c5dffd4a119e39bfeeb))
* add battle pause menu ([91f8d01](https://github.com/X3ne/hsrpc/commit/91f8d017c5a26581b6454177d32d02f471413dfb))
* add bundle package to bundle static files in .exe ([e3461af](https://github.com/X3ne/hsrpc/commit/e3461af0bc963f23a7e5725955dd3efade811dae))
* add combat presence with combat duration ([a8ce998](https://github.com/X3ne/hsrpc/commit/a8ce998d55422799df673085a0a28686e9e91883))
* add error logging to the logs directory ([99a5e8a](https://github.com/X3ne/hsrpc/commit/99a5e8ad67265038bd44c5c3a999d1a957d76f2b))
* add more character menu tabs ([8de0aa8](https://github.com/X3ne/hsrpc/commit/8de0aa823c5cd906d2d905ea021ffa33eab0b5bd))
* add panic recovery to create a crash report ([faef744](https://github.com/X3ne/hsrpc/commit/faef744819d50bfee402d97d811b1c6dc2ac6b1f))
* added persistent configuration state, redesigned loop time, added a few fields in the GUI application ([dfff321](https://github.com/X3ne/hsrpc/commit/dfff321586e4d4e710d9b65a5168398f36efd660))
* center window after update ([19382f5](https://github.com/X3ne/hsrpc/commit/19382f51988a6e2b799e7be62c099255f722c32d))
* change combat paused icon ([7663751](https://github.com/X3ne/hsrpc/commit/766375169ea902d07754463bd95a115a506f92a4))
* pre-processing added to ocr image to improve text recognition ([34460c2](https://github.com/X3ne/hsrpc/commit/34460c214d15a7ee0179930f9b8207988c8a32e9))
* started working on update system ([27fc8eb](https://github.com/X3ne/hsrpc/commit/27fc8ebf15bceb1ad8f5b4fb0df5eddd5dbf4558))


### Bug Fixes

* change characters coords to improve text detection ([fa7820b](https://github.com/X3ne/hsrpc/commit/fa7820b744d130c73bac0d2d1764f37ea1f1eb49))
* fixed some gosec errors ([54119e8](https://github.com/X3ne/hsrpc/commit/54119e865b9fb454bb10dffbd14c8a7c602cb8a6))
* now combat presence should not be reset when a character uses an ult ([8227433](https://github.com/X3ne/hsrpc/commit/82274336f995f8b58d6e5019c08559f8d851c871))
