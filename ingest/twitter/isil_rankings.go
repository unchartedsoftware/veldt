package twitter

import (
	"bufio"
	//"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/unchartedsoftware/prism/ingest/hdfs"
)

type RankingsMap map[string]uint64

var rankings = make(map[string]RankingsMap)

func isRankingFile(file os.FileInfo) bool {
	return !file.IsDir() &&
		file.Name() != "_SUCCESS" &&
		file.Size() > 0
}

func loadRankingFile(host string, port string, filepath string, rankingsId string) error {
	// get hdfs file reader
	reader, err := hdfs.Open(host, port, filepath)
	if err != nil {
		return err
	}
	// defer close reader
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		// split line
		line := strings.Split(scanner.Text(), ",")
		// get user and rank
		username := line[0]
		rank, err := strconv.ParseUint(line[1], 10, 64)
		if err != nil {
			return nil
		}
		// add user and rank to map
		rankings[rankingsId][username] = rank
	}
	return nil
}

// LoadRanking loads user rankings under the provided directory.
func LoadRanking(host string, port string, path string, rankingsDir string) error {
	rankingsId := strings.TrimSuffix(rankingsDir, filepath.Ext(rankingsDir))
	// add rankings file to map
	rankings[rankingsId] = make(RankingsMap)
	// read directory
	files, err := hdfs.ReadDir(host, port, path+"/"+rankingsId)
	if err != nil {
		return err
	}
	// load individual ranking files
	for _, file := range files {
		if isRankingFile(file) {
			err := loadRankingFile(host, port, path+"/"+rankingsId+"/"+file.Name(), rankingsId)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// GetUserRankings returns map containing the ranks for a particular username.
func GetUserRankings(username string) (map[string]uint64, error) {
	ranks := make(map[string]uint64)
	for rankingId, ranking := range rankings {
		rank, ok := ranking[username]
		if !ok {
			//return nil, errors.New("User has not been assigned a rank.")
		}
		ranks[rankingId] = rank
	}
	return ranks, nil
}
