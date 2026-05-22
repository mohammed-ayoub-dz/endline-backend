package search

import (
    "math"
    "sort"
    "strings"
)

type SearchResult struct {
    Video   Video
    Score   float64
}

func editDistance(a, b string) int {
    f := make([]int, len(b)+1)
    for j := range f {
        f[j] = j
    }
    for _, ca := range a {
        j := 1
        nw := f[0]
        f[0]++
        for _, cb := range b {
            cur := f[j]
            if ca == cb {
                f[j] = nw
            } else {
                f[j] = int(math.Min(float64(f[j]+1), math.Min(float64(f[j-1]+1), float64(nw+1))))
            }
            nw = cur
            j++
        }
    }
    return f[len(b)]
}

func (idx *Index) Search(query string, limit int) []SearchResult {
    queryTerms := tokenize(query)
    if len(queryTerms) == 0 {
        return nil
    }

    scores := make(map[int]float64)
    totalDocs := float64(len(idx.docs))

    for _, term := range queryTerms {
        matchedTerms := make(map[string]float64) 

        if _, exists := idx.inverted[term]; exists {
            matchedTerms[term] = 1.0
        }

        if len(matchedTerms) == 0 && len(term) > 2 {
            for existingTerm := range idx.inverted {
                if editDistance(term, existingTerm) <= 1 || (strings.HasPrefix(existingTerm, term) && len(existingTerm)-len(term) <= 2) {
                    matchedTerms[existingTerm] = 0.6 
                }
            }
        }

        for matchedTerm, penalty := range matchedTerms {
            docMap := idx.inverted[matchedTerm]
            
            idf := math.Log(totalDocs / float64(len(docMap)))
            
            for docID, tf := range docMap {
                baseScore := float64(tf) * idf * penalty
                
                video := idx.docs[docID]
                if strings.Contains(strings.ToLower(video.Title), matchedTerm) {
                    baseScore *= 2.5 
                }

                scores[docID] += baseScore
            }
        }
    }

    results := make([]SearchResult, 0, len(scores))
    for docID, score := range scores {
        results = append(results, SearchResult{
            Video: idx.docs[docID],
            Score: score,
        })
    }
    
    sort.Slice(results, func(i, j int) bool {
        return results[i].Score > results[j].Score
    })

    if limit > 0 && len(results) > limit {
        results = results[:limit]
    }
    return results
}