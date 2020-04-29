/*
Package getzap provides a collection of preset configurations for uber-go/zap logger.

For each config, the following aspects are considered:

	| Setting              | Option #1      | Option #2 |
	| -------------------- | -------------- | --------- |
	| Environment          | Dev            | Prod      |
	| Format               | TSV with color | JSON      |
	| Output Target        | Console        | File      |
	| Separate Errors      | Yes            | No        |
	| Compress Legacy Logs | Yes            | No        |
*/
package getzap
