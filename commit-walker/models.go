package main

import "time"

type FileDiff struct {
	NumAdditions int
	NumDeletions int
	Renamed      bool
	FileAdded    bool
	FileDeleted  bool
	OldName      string
	FileName     string
}

type WalkerOutput struct {
	RepositoryLocation string         `json:"repository_location"`
	Commits            []CommitOutput `json:"commits"`
}

type CommitOutput struct {
	AuthorName     string        `json:"author_string"`
	AuthorEmail    string        `json:"author_email"`
	AuthorDate     string        `json:"author_date"`
	CommitterName  string        `json:"committer_name"`
	CommitterDate  string        `json:"committer_date"`
	CommitterEmail string        `json:"committer_email"`
	Message        string        `json:"message"`
	Hash           string        `json:"hash"`
	Signature      string        `json:"signature"`
	ChangedFiles   []ChangedFile `json:"changed_files"`
}

// TODO: Add more features here such as renamed files, etc.
type ChangedFile struct {
	Name         string `json:"name"`
	NumAdditions int    `json:"num_additions"`
	NumDeletions int    `json:"num_deletions"`
}

type Config struct {
	FromDate           time.Time
	FromDateSet        bool
	RepositoryLocation string
	OutputFileName     string
}

var GlobalConfig Config
