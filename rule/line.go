package rule

var lineRulesJSON = []byte(`
[
	{
		"type": "regex",
		"pattern": "-----BEGIN \\S+ PRIVATE KEY-----",
		"caption": "Possible private key data",
		"description": null
	},
	{
		"type": "regex",
		"pattern": "(xox[p|b|o|a]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32})",
		"caption": "Possible Slack API token",
		"description": null
	},
	{
		"type": "regex",
		"pattern": "AKIA[0-9A-Z]{16}",
		"caption": "Possible AWS Access Key",
		"description": null
	}
]	
`)
