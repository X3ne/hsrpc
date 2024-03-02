# Changelog

## [1.3.0](https://github.com/X3ne/hsrpc/compare/v1.2.0...v1.3.0) (2024-03-02)


### Features

* added a crash report GUI ([3dab4e5](https://github.com/X3ne/hsrpc/commit/3dab4e59ca2be5e8f30d5b7f8ed95496f9bb3680))
* create logs.log file to store info logs ([41c59a1](https://github.com/X3ne/hsrpc/commit/41c59a14d60c7fa2cdfa08e7b1b58a9b91269003))
* **GUI:** add theme ([ec5729d](https://github.com/X3ne/hsrpc/commit/ec5729d408aeaaa82e2266c170f7c327e2e41880))
* **GUI:** add theme color based on windows accent color ([1cfda59](https://github.com/X3ne/hsrpc/commit/1cfda59c1c73cb069739a51b65a528c2074816f2))
* hide GUI window by default ([245da89](https://github.com/X3ne/hsrpc/commit/245da892e2c85397511ed59fb2da8e9e27302c73))


### Bug Fixes

* **config:** few corrections to the position of the menus ([b5ac9aa](https://github.com/X3ne/hsrpc/commit/b5ac9aa6977ae88d890e0815da7238e7746cbfb7))
* **config:** sync of config file values with default values when the value is not found to avoid crashes ([8631edd](https://github.com/X3ne/hsrpc/commit/8631eddf9c52d1a3bea43ece6e8671daf9a517fc))
* **game data:** typo fix `Tavel Log` =&gt; `Travel Log` ([5a67388](https://github.com/X3ne/hsrpc/commit/5a67388f92e8c3f9dd0878dedae9eecbab929bdb))
* game OCR results should be better with lower resolutions [#8](https://github.com/X3ne/hsrpc/issues/8) ([50dbff6](https://github.com/X3ne/hsrpc/commit/50dbff6e720dd89eb2188abb29c8225b283457bf))
* **GUI:** add padding on containers ([76d3da8](https://github.com/X3ne/hsrpc/commit/76d3da8c7a2dbd700f4df3f47b1c17317e094040))
* **GUI:** discord gateway connexion no longer blocks the main thread ([5391a40](https://github.com/X3ne/hsrpc/commit/5391a40aea30e620adb1d2d21c7775efa1f839dc))
* **GUI:** discord reconnect button no longer stops main thread ([f7cf8a2](https://github.com/X3ne/hsrpc/commit/f7cf8a24d992982c1505ecabcf3dc207a69c1eab))
* hide tesseract console when the app is built ([8dfeeca](https://github.com/X3ne/hsrpc/commit/8dfeecaa795a94b9b02d577746fbe0e2554e7f9a))
* missing spaces in game data loader logs ([6af9e10](https://github.com/X3ne/hsrpc/commit/6af9e10810410dabf305214041aefe6b5181c255))
* now the PreprocessThreshold value change when changed in the GUI ([7ec90d5](https://github.com/X3ne/hsrpc/commit/7ec90d5981ec6184a5e83272dc74acf86da3afba))
* prevent threshold parameters from being overwritten ([6757223](https://github.com/X3ne/hsrpc/commit/675722384c5502fc7555e38e7dedfff85ee984b5))
* reduce PreprocessThreshold value ([463f11e](https://github.com/X3ne/hsrpc/commit/463f11ecb0be2caced5e09b858de82fa9e284168))
* remove adjustement on width for resolutions &gt; 1920 to reduce ocr errors ([73153af](https://github.com/X3ne/hsrpc/commit/73153af65218e035d068e4744f4adc419bf85604))
* removed brightness treshhold ([1a09648](https://github.com/X3ne/hsrpc/commit/1a09648e4a1cfcd8c2dff4e2a8274036ac2b4ed6))
* some adjustments to coordinate calculations ([fed437c](https://github.com/X3ne/hsrpc/commit/fed437c3a7c6e0918be09d7dccb9eead732a97be))


### Performance Improvements

* add image preprocessing to combat ocr ([232ba5a](https://github.com/X3ne/hsrpc/commit/232ba5a4c743c8c44566dee81ae7a72377db2f80))
* switch from image buffer to path (this seems to greatly improve ocr performance/accuracy) ([c96e8a8](https://github.com/X3ne/hsrpc/commit/c96e8a8d7a18109f415459e3bc61a43abd9b37dd))

## [1.2.0](https://github.com/X3ne/hsrpc/compare/v1.1.0...v1.2.0) (2024-02-20)


### Features

* the program now detects window size and automatically adapts the coordinates of the game's ui, now supports multiple screens ([48ed306](https://github.com/X3ne/hsrpc/commit/48ed306a5b792d4836b0807b27d4d63baddfa1bf))

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
