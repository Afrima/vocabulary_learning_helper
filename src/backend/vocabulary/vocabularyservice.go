package vocabulary

import (
	"errors"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GenerateTestRequest struct {
	ListIDs          []primitive.ObjectID `json:"listIds"`
	Limit            int8                 `json:"limit"`
	FirstValueField  string               `json:"firstValueField"`
	SecondValueField string               `json:"secondValueField"`
}

type UserDBVocabs struct {
	ID         primitive.ObjectID `json:"id"`
	UserFirst  Values             `json:"userFirst"`
	UserSecond Values             `json:"userSecond"`
	DBFirst    Values             `json:"dbFirst"`
	DBSecond   Values             `json:"dbSecond"`
}

type TestResult struct {
	Vocabs  []UserDBVocabs `json:"vocabs"`
	Correct int8           `json:"correct"`
}

type CheckTestRequest struct {
	Vocabularies     []Vocabulary `json:"vocabularies"`
	FirstValueField  string       `json:"firstValueField"`
	SecondValueField string       `json:"secondValueField"`
}

func generateTest(testReqBody GenerateTestRequest) ([]Vocabulary, error) {
	vocabs, err := GetRandomVocabularyByListIds(testReqBody.ListIDs, testReqBody.Limit)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return buildTestResponse(vocabs, testReqBody.FirstValueField, testReqBody.SecondValueField), nil
}

func buildTestResponse(vocabs []Vocabulary, firstValueField string, secondValueField string) []Vocabulary {
	responseVocabularies := make([]Vocabulary, 0, len(vocabs))
	for _, vocab := range vocabs {
		secondValue := vocab.GetValueByKey(secondValueField)
		if secondValue.Key != "" {
			firstValue := vocab.GetValueByKey(firstValueField)
			secondValue.Values = nil
			newValue := make([]Values, 0, 2)
			newValue = append(newValue, firstValue)
			newValue = append(newValue, secondValue)
			responseVocabularies = append(responseVocabularies,
				Vocabulary{ID: vocab.ID,
					ListID: vocab.ListID,
					Values: newValue})
		}
	}
	return responseVocabularies
}

func checkIfVocabEquals(vocab string, values []string) bool {
	for _, dbValue := range values {
		if strings.EqualFold(strings.TrimSpace(vocab), strings.TrimSpace(dbValue)) {
			return true
		}
	}
	return false
}

func getUserDBVocab(firstValueKey string, secondValueKey string, userVocab Vocabulary,
	dbVocab Vocabulary) (UserDBVocabs, error) {
	userDBVocab := UserDBVocabs{ID: dbVocab.ID,
		DBFirst:    dbVocab.GetValueByKey(firstValueKey),
		DBSecond:   dbVocab.GetValueByKey(secondValueKey),
		UserFirst:  userVocab.GetValueByKey(firstValueKey),
		UserSecond: userVocab.GetValueByKey(secondValueKey)}
	if userDBVocab.UserSecond.Key == "" || userDBVocab.DBSecond.Key == "" || userDBVocab.DBFirst.Key == "" || userDBVocab.UserFirst.Key == "" {
		return UserDBVocabs{}, errors.New("one field does not exist")
	}
	return userDBVocab, nil
}

func checkTest(correctVocabs []Vocabulary, checkRequestBody CheckTestRequest) (TestResult, error) {
	correct := int8(0)
	correctVocabMap := make(map[primitive.ObjectID]Vocabulary, len(correctVocabs))
	for _, correctVocab := range correctVocabs {
		correctVocabMap[correctVocab.ID] = correctVocab
	}
	userDBVocabs := make([]UserDBVocabs, 0, len(correctVocabs))
	for _, vocab := range checkRequestBody.Vocabularies {
		correctVocab := correctVocabMap[vocab.ID]
		userDBVocab, err := getUserDBVocab(checkRequestBody.FirstValueField, checkRequestBody.SecondValueField, vocab, correctVocab)
		if err != nil {
			return TestResult{}, err
		}
		valueCorrect := false
		for _, userValue := range userDBVocab.UserSecond.Values {
			if valueCorrect = checkIfVocabEquals(userValue, userDBVocab.DBSecond.Values); !valueCorrect {
				break
			}
		}
		if valueCorrect {
			correct++
		}
		userDBVocabs = append(userDBVocabs, userDBVocab)
	}
	correction := TestResult{Vocabs: userDBVocabs, Correct: correct}
	return correction, nil
}
