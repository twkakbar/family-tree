package domain

import "time"

// TreeTransferFormatVersion is the JSON export schema version.
const TreeTransferFormatVersion = 1

// TreeTransfer is the portable JSON document for export/import.
type TreeTransfer struct {
	FormatVersion int                  `json:"format_version"`
	ExportedAt    time.Time            `json:"exported_at"`
	Family        TreeTransferFamily   `json:"family"`
	Persons       []TreeTransferPerson `json:"persons"`
	Relationships []TreeTransferEdge   `json:"relationships"`
}

// TreeTransferFamily carries family metadata (name/description only on import).
type TreeTransferFamily struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TreeTransferPerson is a portable person row. ID is the export-time UUID used
// to remap relationship endpoints on import.
type TreeTransferPerson struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Gender    Gender  `json:"gender"`
	BirthDate *string `json:"birth_date"` // YYYY-MM-DD
	DeathDate *string `json:"death_date"`
	Bio       string  `json:"bio"`
}

// TreeTransferEdge is a portable relationship using export-time person IDs.
type TreeTransferEdge struct {
	FromPersonID string           `json:"from_person_id"`
	ToPersonID   string           `json:"to_person_id"`
	Type         RelationshipType `json:"type"`
}

// ImportResult summarizes an additive import.
type ImportResult struct {
	PersonsCreated       int `json:"persons_created"`
	RelationshipsCreated int `json:"relationships_created"`
}
