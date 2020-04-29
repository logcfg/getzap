/*
Package getzap provides a collection of sophisticated uber-go/zap loggers.

For each logger, the following aspects are considered:

	| Setting              | Option #1      | Option #2 |
	| -------------------- | -------------- | --------- |
	| Environment          | Dev            | Prod      |
	| Format               | TSV with color | JSON      |
	| Output Target        | Console        | File      |
	| Separate Errors      | Yes            | No        |
	| Compress Legacy Logs | Yes            | No        |
*/
package getzap
