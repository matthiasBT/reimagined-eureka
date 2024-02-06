package commands

import "strings"

const minMasterKeyLength = 8
const maxMasterKeyLength = 16
const minEntropyLength = 16
const maxEntropyLength = 32

const secretTypeCreds = "creds"
const secretTypeNotes = "notes"
const secretTypeFiles = "files"
const secretTypeCards = "cards"
const secretTypeAll = "all"

const secretDelimiterWidth = 50
const secretDelimiterChar = "-"

const cardNumberMinLength = 13
const cardNumberMaxLength = 19

const monthMin = 1
const monthMax = 12
const yearMin = 20
const monthMinChars = 1
const monthMaxChars = 2
const yearMinChars = 2
const yearMaxChars = 2
const cscMinChars = 3
const cscMaxChars = 3
const nameMinChars = 1

var supportedTypes = []string{
	secretTypeCreds,
	secretTypeNotes,
	secretTypeFiles,
	secretTypeCards,
	secretTypeAll,
}

func listSupportedTypes() string {
	return strings.TrimSpace(strings.Join(supportedTypes, ", "))
}
