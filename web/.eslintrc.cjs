module.exports = {
	root: true,
	extends: [
		"eslint:recommended",
		"plugin:@typescript-eslint/recommended",
		"plugin:svelte/recommended",
	],
	parser: "@typescript-eslint/parser",
	plugins: [ "@typescript-eslint" ],
	parserOptions: {
		sourceType: "module",
		ecmaVersion: 2020,
		extraFileExtensions: [ ".svelte" ]
	},
	env: {
		browser: true,
		es2017: true,
		node: true
	},
	overrides: [
		{
			files: [ "*.svelte", ".ts" ],
			parser: "svelte-eslint-parser",
			parserOptions: { parser: "@typescript-eslint/parser" }
		}
	],
	"rules": {
		"@typescript-eslint/no-explicit-any": "off",
		"no-mixed-spaces-and-tabs": "off",
		"@typescript-eslint/no-unused-vars": "warn",
		"indent": [
			"error",
			"tab",
			{ "SwitchCase": 1 }
		],
		"linebreak-style": [
			"error",
			"unix"
		],
		"quotes": [
			"error",
			"double"
		],
		"semi": [
			"error",
			"always"
		],
		"max-len": [ "error", { "code": 120 } ],
		"space-infix-ops": [ "error", { "int32Hint": false } ],
		"space-in-parens": [ "error", "never" ],
		"space-before-blocks": [ "error", "always" ],
		"array-bracket-spacing": [ "error", "always" ],
		"object-curly-spacing": [ "error", "always" ],
		"object-curly-newline": [ "error", {
			"multiline": true,
			"minProperties": 5,
			"consistent": false
		} ],
		"brace-style": [ "error", "1tbs" ],
		"object-property-newline": [ "error", { "allowMultiplePropertiesPerLine": false } ],
		"keyword-spacing": [ "error", {
			"before": true,
			"after": true
		} ],
		"comma-spacing": [ "error", { "after": true } ],
		"key-spacing": [ "error", {
			"beforeColon": false,
			"afterColon": true,
			"mode": "strict"
		} ],
		"arrow-spacing": [ "error", {
			"before": true,
			"after": true
		} ]
	}
};
