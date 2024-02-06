package commands

import "strings"

const MinMasterKeyLength = 8
const MaxMasterKeyLength = 16
const MinEntropyLength = 16
const MaxEntropyLength = 32

const SecretTypeCreds = "creds"
const SecretTypeNotes = "notes"
const SecretTypeFiles = "files"
const SecretTypeCards = "cards"
const SecretTypeAll = "all"

const SecretDelimiterWidth = 50
const SecretDelimiterChar = "-"

const CardNumberMinLength = 13
const CardNumberMaxLength = 19

const MonthMin = 1
const MonthMax = 12
const YearMin = 20
const MonthMinChars = 1
const MonthMaxChars = 2
const YearMinChars = 2
const YearMaxChars = 2
const CSCMinChars = 3
const CSCMaxChars = 3
const NameMinChars = 1

var SupportedTypes = []string{
	SecretTypeCreds,
	SecretTypeNotes,
	SecretTypeFiles,
	SecretTypeCards,
	SecretTypeAll,
}

func ListSupportedTypes() string {
	return strings.TrimSpace(strings.Join(SupportedTypes, ", "))
}
