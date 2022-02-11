package main

import (
	"crypto/sha1"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func walkGitRepositoryCommits(c *cli.Context) error { // nolint:funlen
	err := setConfigurationValues(c)
	if err != nil {
		return errors.WithMessage(err, "setConfigurationValues")
	}

	walkerOutput := WalkerOutput{
		RepositoryLocation: GlobalConfig.RepositoryLocation,
		Commits:            []CommitOutput{},
	}

	// Open the repository local as a git repo.
	repo, err := git.PlainOpen(GlobalConfig.RepositoryLocation)
	if err != nil {
		return err
	}

	// Grab the commit iterator.
	commitIter, err := repo.CommitObjects()
	if err != nil {
		return err
	}

	commit, err := commitIter.Next()

	// Iterate through the list of commits
	for commit != nil && err == nil {
		// If we have a from date set and the commit's date is before that date, skip to the next one.
		if GlobalConfig.FromDateSet && commit.Committer.When.Before(GlobalConfig.FromDate) {
			commit, err = commitIter.Next()
			if err != nil {
				break
			}

			continue
		}

		outputBlob := CommitOutput{
			AuthorDate:     commit.Author.When.String(),
			AuthorName:     commit.Author.Name,
			AuthorEmail:    commit.Author.Email,
			CommitterName:  commit.Committer.Name,
			CommitterEmail: commit.Committer.Email,
			CommitterDate:  commit.Author.When.String(),
			Hash:           commit.Hash.String(),
			Message:        commit.Message,
		}

		// Store the commit's PGP signature if one exists
		if commit.PGPSignature != "" {
			h := sha1.New()
			h.Write([]byte(commit.PGPSignature))
			bs := h.Sum(nil)

			outputBlob.Signature = fmt.Sprintf("%x", bs)
		}

		currTree, err := commit.Tree()
		if err != nil {
			return err
		}

		parentCommit, err := commit.Parent(0)
		if err != nil {
			commit, err = commitIter.Next()
			if err != nil {
				break
			}

			continue // This would be the first commit in the repo.
		}

		oldTree, err := parentCommit.Tree()
		if err != nil {
			return err
		}

		changes, err := object.DiffTreeWithOptions(c.Context, oldTree, currTree, &object.DiffTreeOptions{
			DetectRenames: true,
		})
		if err != nil {
			return err
		}

		patch, err := changes.Patch()
		if err != nil {
			return err
		}

		diffs := GetSignificantChangesFromPatch(patch)

		for k, v := range diffs {
			fileChange := ChangedFile{
				Name:         k,
				NumAdditions: v.NumAdditions,
				NumDeletions: v.NumDeletions,
			}

			outputBlob.ChangedFiles = append(outputBlob.ChangedFiles, fileChange)
		}

		walkerOutput.Commits = append(walkerOutput.Commits, outputBlob)

		commit, err = commitIter.Next()
		if err != nil {
			break
		}
	}

	err = writeOutputToFile(walkerOutput)
	if err != nil {
		return err
	}

	return nil
}

// GetSignificantChangesFromPatch returns the set of files that have changed in a given patch, if the changes are significant
func GetSignificantChangesFromPatch(patch *object.Patch) map[string]*FileDiff {
	changedFiles := map[string]*FileDiff{}

	for _, file := range patch.FilePatches() {
		numChanges := 0
		from, to := file.Files()

		fileDiff := FileDiff{}

		chunks := file.Chunks()
		for _, chunk := range chunks {
			if chunk.Type() != 0 { // nolint:nestif
				if chunk.Type() == 1 {
					fileDiff.NumAdditions++
					numChanges++
				} else if chunk.Type() == 2 {
					fileDiff.NumDeletions++
					numChanges++
				}
			}
		}

		if from == nil {
			fileDiff.FileAdded = true
		}

		if to == nil {
			fileDiff.FileDeleted = true
			fileDiff.FileName = from.Path()
		} else {
			fileDiff.FileName = to.Path()
		}

		if from != nil && to != nil && from.Path() != to.Path() {
			fileDiff.Renamed = true
			fileDiff.OldName = from.Path()
		}

		if numChanges > 0 {
			changedFiles[fileDiff.FileName] = &fileDiff
		}
	}

	return changedFiles
}
