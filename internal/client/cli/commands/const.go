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
