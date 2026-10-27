package interpreter

var alterMap = map[TokenType]map[TokenType]struct{}{
	ADD: {
		COLUMN: {},
		CONSTRAINT: {},
		CHECK: {},
		FIELD: {},
	},
	ALTER: {
		COLUMN: {},
		FIELD: {},
	},
	DROP: {
		COLUMN: {},
		CONSTRAINT: {},
		FIELD: {},
	},
	RENAME: {
		COLUMN: {},
		FIELD: {},
		TO: {},
	},
}
