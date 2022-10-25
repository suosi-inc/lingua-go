/*
 * Copyright © 2021 Peter M. Stahl pemistahl@gmail.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either expressed or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lingua

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math"
	"sync"
	"testing"
)

// ##############################
// MOCKS
// ##############################

type mockedTrainingDataLanguageModel struct {
	mock.Mock
}

func (m *mockedTrainingDataLanguageModel) getRelativeFrequency(ngram ngram) float64 {
	return m.Called(ngram).Get(0).(float64)
}

func createTrainingModelMock(data map[string]float64) *mockedTrainingDataLanguageModel {
	model := new(mockedTrainingDataLanguageModel)
	for ngram, probability := range data {
		model.On("getRelativeFrequency", newNgram(ngram)).Return(probability)
	}
	return model
}

// ##############################
// LANGUAGE MODELS FOR ENGLISH
// ##############################

func unigramModelForEnglish() languageModel {
	return createTrainingModelMock(map[string]float64{
		"a": 0.01,
		"l": 0.02,
		"t": 0.03,
		"e": 0.04,
		"r": 0.05,
		// unknown unigrams
		"w": 0.0,
	})

}

func bigramModelForEnglish() languageModel {
	return createTrainingModelMock(map[string]float64{
		"al": 0.11,
		"lt": 0.12,
		"te": 0.13,
		"er": 0.14,
		// unknown bigrams
		"aq": 0.0,
		"wx": 0.0,
	})
}

func trigramModelForEnglish() languageModel {
	return createTrainingModelMock(map[string]float64{
		"alt": 0.19,
		"lte": 0.2,
		"ter": 0.21,
		// unknown trigrams
		"aqu": 0.0,
		"tez": 0.0,
		"wxy": 0.0,
	})
}

func quadrigramModelForEnglish() languageModel {
	return createTrainingModelMock(map[string]float64{
		"alte": 0.25,
		"lter": 0.26,
		// unknown quadrigrams
		"aqua": 0.0,
		"wxyz": 0.0,
	})
}

func fivegramModelForEnglish() languageModel {
	return createTrainingModelMock(map[string]float64{
		"alter": 0.29,
		// unknown fivegrams
		"aquas": 0.0,
	})
}

// ##############################
// LANGUAGE MODELS FOR GERMAN
// ##############################

func unigramModelForGerman() languageModel {
	return createTrainingModelMock(map[string]float64{
		"a": 0.06,
		"l": 0.07,
		"t": 0.08,
		"e": 0.09,
		"r": 0.1,
		// unknown unigrams
		"w": 0.0,
	})
}

func bigramModelForGerman() languageModel {
	return createTrainingModelMock(map[string]float64{
		"al": 0.15,
		"lt": 0.16,
		"te": 0.17,
		"er": 0.18,
		// unknown bigrams
		"wx": 0.0,
	})
}

func trigramModelForGerman() languageModel {
	return createTrainingModelMock(map[string]float64{
		"alt": 0.22,
		"lte": 0.23,
		"ter": 0.24,
		// unknown trigrams
		"wxy": 0.0,
	})
}

func quadrigramModelForGerman() languageModel {
	return createTrainingModelMock(map[string]float64{
		"alte": 0.27,
		"lter": 0.28,
		// unknown quadrigrams
		"wxyz": 0.0,
	})
}

func fivegramModelForGerman() languageModel {
	return createTrainingModelMock(map[string]float64{
		"alter": 0.3,
	})
}

// ##############################
// TEST DATA MODELS
// ##############################

func testDataModel(strs []string) testDataLanguageModel {
	ngrams := make(map[ngram]bool)
	for _, s := range strs {
		ngrams[newNgram(s)] = true
	}
	return testDataLanguageModel{ngrams}
}

// ##############################
// DETECTORS
// ##############################

func newDetectorForEnglishAndGerman() languageDetector {
	var unigramLanguageModels sync.Map
	unigramLanguageModels.Store(English, unigramModelForEnglish())
	unigramLanguageModels.Store(German, unigramModelForGerman())

	var bigramLanguageModels sync.Map
	bigramLanguageModels.Store(English, bigramModelForEnglish())
	bigramLanguageModels.Store(German, bigramModelForGerman())

	var trigramLanguageModels sync.Map
	trigramLanguageModels.Store(English, trigramModelForEnglish())
	trigramLanguageModels.Store(German, trigramModelForGerman())

	var quadrigramLanguageModels sync.Map
	quadrigramLanguageModels.Store(English, quadrigramModelForEnglish())
	quadrigramLanguageModels.Store(German, quadrigramModelForGerman())

	var fivegramLanguageModels sync.Map
	fivegramLanguageModels.Store(English, fivegramModelForEnglish())
	fivegramLanguageModels.Store(German, fivegramModelForGerman())

	return languageDetector{
		languages:                     []Language{English, German},
		minimumRelativeDistance:       0.0,
		languagesWithUniqueCharacters: []Language{},
		oneLanguageAlphabets:          map[alphabet]Language{},
		unigramLanguageModels:         &unigramLanguageModels,
		bigramLanguageModels:          &bigramLanguageModels,
		trigramLanguageModels:         &trigramLanguageModels,
		quadrigramLanguageModels:      &quadrigramLanguageModels,
		fivegramLanguageModels:        &fivegramLanguageModels,
	}
}

func newDetectorForAllLanguages() languageDetector {
	languages := AllLanguages()
	var emptyLanguageModels sync.Map
	return languageDetector{
		languages:                     languages,
		minimumRelativeDistance:       0.0,
		languagesWithUniqueCharacters: collectLanguagesWithUniqueCharacters(languages),
		oneLanguageAlphabets:          collectOneLanguageAlphabets(languages),
		unigramLanguageModels:         &emptyLanguageModels,
		bigramLanguageModels:          &emptyLanguageModels,
		trigramLanguageModels:         &emptyLanguageModels,
		quadrigramLanguageModels:      &emptyLanguageModels,
		fivegramLanguageModels:        &emptyLanguageModels,
	}
}

var detectorForEnglishAndGerman = newDetectorForEnglishAndGerman()
var detectorForAllLanguages = newDetectorForAllLanguages()

// ##############################
// TESTS
// ##############################

var delta = 0.00000000000001

func TestCleanUpInputText(t *testing.T) {
	text := `Weltweit    gibt es ungefähr 6.000 Sprachen,
    wobei laut Schätzungen zufolge ungefähr 90  Prozent davon
    am Ende dieses Jahrhunderts verdrängt sein werden.`

	expectedCleanedText := "weltweit gibt es ungefähr sprachen wobei laut schätzungen " +
		"zufolge ungefähr prozent davon am ende dieses jahrhunderts verdrängt sein werden"

	assert.Equal(t, expectedCleanedText, detectorForAllLanguages.cleanUpInputText(text))
}

func TestSplitTextIntoWords(t *testing.T) {
	testCases := []struct {
		text          string
		expectedWords []string
	}{
		{
			"this is a sentence",
			[]string{"this", "is", "a", "sentence"},
		},
		{
			"sentence",
			[]string{"sentence"},
		},
		{
			"上海大学是一个好大学 this is a sentence",
			[]string{"上", "海", "大", "学", "是", "一", "个", "好", "大", "学", "this", "is", "a", "sentence"},
		},
	}
	for _, testCase := range testCases {
		assert.Equal(
			t,
			testCase.expectedWords,
			detectorForAllLanguages.splitTextIntoWords(testCase.text),
			fmt.Sprintf("unexpected tokenization for text '%s'", testCase.text),
		)
	}
}

func TestLookUpNgramProbability(t *testing.T) {
	testCases := []struct {
		language            Language
		ngram               string
		expectedProbability float64
	}{
		{English, "a", 0.01},
		{English, "lt", 0.12},
		{English, "ter", 0.21},
		{English, "alte", 0.25},
		{English, "alter", 0.29},
		{German, "t", 0.08},
		{German, "er", 0.18},
		{German, "alt", 0.22},
		{German, "lter", 0.28},
		{German, "alter", 0.3},
	}
	for _, testCase := range testCases {
		probability := detectorForEnglishAndGerman.lookUpNgramProbability(testCase.language, newNgram(testCase.ngram))
		message := fmt.Sprintf(
			"expected probability %v for language %v and ngram '%s', got %v",
			testCase.expectedProbability,
			testCase.language,
			testCase.ngram,
			probability,
		)
		assert.Equal(t, testCase.expectedProbability, probability, message)
	}

	assert.Panicsf(t, func() {
		detectorForEnglishAndGerman.lookUpNgramProbability(English, newNgram(""))
	}, "zerogram detected")
}

func TestComputeSumOfNgramProbabilities(t *testing.T) {
	testCases := []struct {
		ngrams                     []string
		expectedSumOfProbabilities float64
	}{
		{
			[]string{"a", "l", "t", "e", "r"},
			math.Log(0.01) + math.Log(0.02) + math.Log(0.03) + math.Log(0.04) + math.Log(0.05),
		},
		{
			// back off unknown Trigram("tez") to known Bigram("te")
			[]string{"alt", "lte", "tez"},
			math.Log(0.19) + math.Log(0.2) + math.Log(0.13),
		},
		{
			// back off unknown Fivegram("aquas") to known Unigram("a")
			[]string{"aquas"},
			math.Log(0.01),
		},
	}
	for _, testCase := range testCases {
		mappedNgrams := make(map[ngram]bool)
		for _, ngram := range testCase.ngrams {
			mappedNgrams[newNgram(ngram)] = true
		}
		sumOfProbabilities := detectorForEnglishAndGerman.computeSumOfNgramProbabilities(English, mappedNgrams)
		message := fmt.Sprintf(
			"expected sum %v for language %v and ngrams %v, got %v",
			testCase.expectedSumOfProbabilities,
			English,
			testCase.ngrams,
			sumOfProbabilities,
		)
		assert.InDelta(t, testCase.expectedSumOfProbabilities, sumOfProbabilities, delta, message)
	}
}

func TestComputeLanguageProbabilities(t *testing.T) {
	testCases := []struct {
		testDataModel         testDataLanguageModel
		expectedProbabilities map[Language]float64
	}{
		{
			testDataModel([]string{"a", "l", "t", "e", "r"}),
			map[Language]float64{
				English: math.Log(0.01) + math.Log(0.02) + math.Log(0.03) + math.Log(0.04) + math.Log(0.05),
				German:  math.Log(0.06) + math.Log(0.07) + math.Log(0.08) + math.Log(0.09) + math.Log(0.1),
			},
		},
		{
			testDataModel([]string{"alt", "lte", "ter", "wxy"}),
			map[Language]float64{
				English: math.Log(0.19) + math.Log(0.2) + math.Log(0.21),
				German:  math.Log(0.22) + math.Log(0.23) + math.Log(0.24),
			},
		},
		{
			testDataModel([]string{"alte", "lter", "wxyz"}),
			map[Language]float64{
				English: math.Log(0.25) + math.Log(0.26),
				German:  math.Log(0.27) + math.Log(0.28),
			},
		},
	}
	languages := []Language{English, German}
	for _, testCase := range testCases {
		probabilities := detectorForEnglishAndGerman.computeLanguageProbabilities(testCase.testDataModel, languages)

		for language, probability := range probabilities {
			expectedProbability := testCase.expectedProbabilities[language]
			message := fmt.Sprintf(
				"expected probability %v for language %v, got %v",
				expectedProbability,
				language,
				probability,
			)
			assert.InDelta(t, expectedProbability, probability, delta, message)
		}
	}
}

func TestComputeLanguageConfidenceValues(t *testing.T) {
	unigramCountForBothLanguages := 5.0
	totalProbabilityForGerman := (
	// unigrams
	math.Log(0.06) + math.Log(0.07) + math.Log(0.08) + math.Log(0.09) + math.Log(0.1) +
		// bigrams
		math.Log(0.15) + math.Log(0.16) + math.Log(0.17) + math.Log(0.18) +
		// trigrams
		math.Log(0.22) + math.Log(0.23) + math.Log(0.24) +
		// quadrigrams
		math.Log(0.27) + math.Log(0.28) +
		// fivegrams
		math.Log(0.3)) / unigramCountForBothLanguages

	totalProbabilityForEnglish := (
	// unigrams
	math.Log(0.01) + math.Log(0.02) + math.Log(0.03) + math.Log(0.04) + math.Log(0.05) +
		// bigrams
		math.Log(0.11) + math.Log(0.12) + math.Log(0.13) + math.Log(0.14) +
		// trigrams
		math.Log(0.19) + math.Log(0.2) + math.Log(0.21) +
		// quadrigrams
		math.Log(0.25) + math.Log(0.26) +
		// fivegrams
		math.Log(0.29)) / unigramCountForBothLanguages

	expectedConfidenceForGerman := 1.0
	expectedConfidenceForEnglish := totalProbabilityForGerman / totalProbabilityForEnglish

	confidenceValues := detectorForEnglishAndGerman.ComputeLanguageConfidenceValues("Alter")

	assert.Equal(
		t,
		2,
		len(confidenceValues),
		fmt.Sprintf("expected 2 confidence values, got %v", len(confidenceValues)),
	)

	first, second := confidenceValues[0], confidenceValues[1]

	assert.Equal(t, German, first.Language())
	assert.Equal(t, expectedConfidenceForGerman, first.Value())

	assert.Equal(t, English, second.Language())
	assert.InDelta(t, expectedConfidenceForEnglish, second.Value(), delta)
}

func TestDetectLanguage(t *testing.T) {
	language, exists := detectorForEnglishAndGerman.DetectLanguageOf("Alter")
	assert.Equal(t, German, language)
	assert.True(t, exists)

	language, exists = detectorForEnglishAndGerman.DetectLanguageOf("проарплап")
	assert.Equal(t, Unknown, language)
	assert.False(t, exists)
}

func TestDetectLanguageWithRules(t *testing.T) {
	testCases := []struct {
		word             string
		expectedLanguage Language
	}{
		// words with unique characters
		{"məhərrəm", Azerbaijani},
		{"substituïts", Catalan},
		{"rozdělit", Czech},
		{"tvořen", Czech},
		{"subjektů", Czech},
		{"nesufiĉecon", Esperanto},
		{"intermiksiĝis", Esperanto},
		{"monaĥinoj", Esperanto},
		{"kreitaĵoj", Esperanto},
		{"ŝpinante", Esperanto},
		{"apenaŭ", Esperanto},
		{"groß", German},
		{"σχέδια", Greek},
		{"fekvő", Hungarian},
		{"meggyűrűzni", Hungarian},
		{"ヴェダイヤモンド", Japanese},
		{"әлем", Kazakh},
		{"шаруашылығы", Kazakh},
		{"ақын", Kazakh},
		{"оның", Kazakh},
		{"шұрайлы", Kazakh},
		{"teoloģiska", Latvian},
		{"blaķene", Latvian},
		{"ceļojumiem", Latvian},
		{"numuriņu", Latvian},
		{"mergelės", Lithuanian},
		{"įrengus", Lithuanian},
		{"slegiamų", Lithuanian},
		{"припаѓа", Macedonian},
		{"ѕидови", Macedonian},
		{"ќерка", Macedonian},
		{"џамиите", Macedonian},
		{"मिळते", Marathi},
		{"үндсэн", Mongolian},
		{"дөхөж", Mongolian},
		{"zmieniły", Polish},
		{"państwowych", Polish},
		{"mniejszości", Polish},
		{"groźne", Polish},
		{"ialomiţa", Romanian},
		{"наслеђивања", Serbian},
		{"неисквареношћу", Serbian},
		{"podĺa", Slovak},
		{"pohľade", Slovak},
		{"mŕtvych", Slovak},
		{"ґрунтовому", Ukrainian},
		{"пропонує", Ukrainian},
		{"пристрої", Ukrainian},
		{"cằm", Vietnamese},
		{"thần", Vietnamese},
		{"chẳng", Vietnamese},
		{"quẩy", Vietnamese},
		{"sẵn", Vietnamese},
		{"nhẫn", Vietnamese},
		{"dắt", Vietnamese},
		{"chất", Vietnamese},
		{"đạp", Vietnamese},
		{"mặn", Vietnamese},
		{"hậu", Vietnamese},
		{"hiền", Vietnamese},
		{"lẻn", Vietnamese},
		{"biểu", Vietnamese},
		{"kẽm", Vietnamese},
		{"diễm", Vietnamese},
		{"phế", Vietnamese},
		{"việc", Vietnamese},
		{"chỉnh", Vietnamese},
		{"trĩ", Vietnamese},
		{"ravị", Vietnamese},
		{"thơ", Vietnamese},
		{"nguồn", Vietnamese},
		{"thờ", Vietnamese},
		{"sỏi", Vietnamese},
		{"tổng", Vietnamese},
		{"nhở", Vietnamese},
		{"mỗi", Vietnamese},
		{"bỡi", Vietnamese},
		{"tốt", Vietnamese},
		{"giới", Vietnamese},
		{"một", Vietnamese},
		{"hợp", Vietnamese},
		{"hưng", Vietnamese},
		{"từng", Vietnamese},
		{"của", Vietnamese},
		{"sử", Vietnamese},
		{"cũng", Vietnamese},
		{"những", Vietnamese},
		{"chức", Vietnamese},
		{"dụng", Vietnamese},
		{"thực", Vietnamese},
		{"kỳ", Vietnamese},
		{"kỷ", Vietnamese},
		{"mỹ", Vietnamese},
		{"mỵ", Vietnamese},
		{"aṣiwèrè", Yoruba},
		{"ṣaaju", Yoruba},
		{"والموضوع", Unknown},
		{"сопротивление", Unknown},
		{"house", Unknown},

		// words with unique alphabet
		{"ունենա", Armenian},
		{"জানাতে", Bengali},
		{"გარეუბან", Georgian},
		{"σταμάτησε", Greek},
		{"ઉપકરણોની", Gujarati},
		{"בתחרויות", Hebrew},
		{"びさ", Japanese},
		{"대결구도가", Korean},
		{"ਮੋਟਰਸਾਈਕਲਾਂ", Punjabi},
		{"துன்பங்களை", Tamil},
		{"కృష్ణదేవరాయలు", Telugu},
		{"ในทางหลวงหมายเลข", Thai},
	}
	for _, testCase := range testCases {
		detectedLanguage := detectorForAllLanguages.detectLanguageWithRules([]string{testCase.word})
		message := fmt.Sprintf(
			"expected %v for word '%s', got %v",
			testCase.expectedLanguage,
			testCase.word,
			detectedLanguage,
		)
		assert.Equal(t, testCase.expectedLanguage, detectedLanguage, message)
	}
}

func TestFilterLanguagesByRules(t *testing.T) {
	testCases := []struct {
		word              string
		expectedLanguages []Language
	}{
		{"والموضوع", []Language{Arabic, Persian, Urdu}},
		{"сопротивление", []Language{
			Belarusian, Bulgarian, Kazakh, Macedonian, Mongolian, Russian, Serbian, Ukrainian},
		},
		{"раскрывае", []Language{Belarusian, Kazakh, Mongolian, Russian}},
		{"этот", []Language{Belarusian, Kazakh, Mongolian, Russian}},
		{"огнём", []Language{Belarusian, Kazakh, Mongolian, Russian}},
		{"плаваща", []Language{Bulgarian, Kazakh, Mongolian, Russian}},
		{"довършат", []Language{Bulgarian, Kazakh, Mongolian, Russian}},
		{"павінен", []Language{Belarusian, Kazakh, Ukrainian}},
		{"затоплување", []Language{Macedonian, Serbian}},
		{"ректасцензија", []Language{Macedonian, Serbian}},
		{"набљудувач", []Language{Macedonian, Serbian}},
		{"aizklātā", []Language{Latvian, Maori, Yoruba}},
		{"sistēmas", []Language{Latvian, Maori, Yoruba}},
		{"palīdzi", []Language{Latvian, Maori, Yoruba}},
		{"nhẹn", []Language{Vietnamese, Yoruba}},
		{"chọn", []Language{Vietnamese, Yoruba}},
		{"prihvaćanju", []Language{Bosnian, Croatian, Polish}},
		{"nađete", []Language{Bosnian, Croatian, Vietnamese}},
		{"visão", []Language{Portuguese, Vietnamese}},
		{"wystąpią", []Language{Lithuanian, Polish}},
		{"budowę", []Language{Lithuanian, Polish}},
		{"nebūsime", []Language{Latvian, Lithuanian, Maori, Yoruba}},
		{"afişate", []Language{Azerbaijani, Romanian, Turkish}},
		{"kradzieżami", []Language{Polish, Romanian}},
		{"înviat", []Language{French, Romanian}},
		{"venerdì", []Language{Italian, Vietnamese, Yoruba}},
		{"años", []Language{Basque, Spanish}},
		{"rozohňuje", []Language{Czech, Slovak}},
		{"rtuť", []Language{Czech, Slovak}},
		{"pregătire", []Language{Romanian, Vietnamese}},
		{"jeďte", []Language{Czech, Romanian, Slovak}},
		{"minjaverðir", []Language{Icelandic, Turkish}},
		{"þagnarskyldu", []Language{Icelandic, Turkish}},
		{"nebûtu", []Language{French, Hungarian}},
		{"hashemidëve", []Language{Afrikaans, Albanian, Dutch, French}},
		{"forêt", []Language{Afrikaans, French, Portuguese, Vietnamese}},
		{"succèdent", []Language{French, Italian, Vietnamese, Yoruba}},
		{"où", []Language{French, Italian, Vietnamese, Yoruba}},
		{"tõeliseks", []Language{Estonian, Hungarian, Portuguese, Vietnamese}},
		{"viòiem", []Language{Catalan, Italian, Vietnamese, Yoruba}},
		{"contrôle", []Language{French, Portuguese, Slovak, Vietnamese}},
		{"direktør", []Language{Bokmal, Danish, Nynorsk}},
		{"vývoj", []Language{Czech, Icelandic, Slovak, Turkish, Vietnamese}},
		{"päralt", []Language{Estonian, Finnish, German, Slovak, Swedish}},
		{"labâk", []Language{French, Portuguese, Romanian, Turkish, Vietnamese}},
		{"pràctiques", []Language{Catalan, French, Italian, Portuguese, Vietnamese}},
		{"überrascht", []Language{
			Azerbaijani, Catalan, Estonian, German, Hungarian, Spanish, Turkish},
		},
		{"indebærer", []Language{Bokmal, Danish, Icelandic, Nynorsk}},
		{"måned", []Language{Bokmal, Danish, Nynorsk, Swedish}},
		{"zaručen", []Language{Bosnian, Czech, Croatian, Latvian, Lithuanian, Slovak, Slovene}},
		{"zkouškou", []Language{Bosnian, Czech, Croatian, Latvian, Lithuanian, Slovak, Slovene}},
		{"navržen", []Language{Bosnian, Czech, Croatian, Latvian, Lithuanian, Slovak, Slovene}},
		{"façonnage", []Language{
			Albanian, Azerbaijani, Basque, Catalan, French, Portuguese, Turkish},
		},
		{"höher", []Language{
			Azerbaijani, Estonian, Finnish, German, Hungarian, Icelandic, Swedish, Turkish},
		},
		{"catedráticos", []Language{
			Catalan, Czech, Icelandic, Irish, Hungarian, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},
		},
		{"política", []Language{
			Catalan, Czech, Icelandic, Irish, Hungarian, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},
		},
		{"música", []Language{
			Catalan, Czech, Icelandic, Irish, Hungarian, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},
		},
		{"contradicció", []Language{
			Catalan, Hungarian, Icelandic, Irish, Polish, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},
		},
		{"només", []Language{
			Catalan, Czech, French, Hungarian, Icelandic, Irish,
			Italian, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},
		},
		{"house", []Language{
			Afrikaans, Albanian, Azerbaijani, Basque, Bokmal, Bosnian, Catalan, Croatian, Czech, Danish, Dutch, English,
			Esperanto, Estonian, Finnish, French, Ganda, German, Hungarian, Icelandic, Indonesian, Irish, Italian,
			Latin, Latvian, Lithuanian, Malay, Maori, Nynorsk, Polish, Portuguese, Romanian, Shona, Slovak, Slovene,
			Somali, Sotho, Spanish, Swahili, Swedish, Tagalog, Tsonga, Tswana, Turkish, Vietnamese, Welsh, Xhosa,
			Yoruba, Zulu},
		},
	}
	for _, testCase := range testCases {
		filteredLanguages := detectorForAllLanguages.filterLanguagesByRules([]string{testCase.word})
		message := fmt.Sprintf("expected %v for word '%s', got %v", testCase.expectedLanguages, testCase.word, filteredLanguages)
		assert.ElementsMatch(t, testCase.expectedLanguages, filteredLanguages, message)
	}
}

func TestDetectionOfInvalidStrings(t *testing.T) {
	testCases := []string{"", " \n  \t;", "3<856%)§"}
	for _, testCase := range testCases {
		language, exists := detectorForAllLanguages.DetectLanguageOf(testCase)
		assert.Equal(t, Unknown, language)
		assert.False(t, exists)
	}
}

func TestLanguageDetectionIsDeterministic(t *testing.T) {
	testCases := []struct {
		text      string
		languages []Language
	}{
		{
			"ام وی با نیکی میناج تیزر داشت؟؟؟؟؟؟ i vote for bts ( _ ) as the _ via ( _ )",
			[]Language{English, Urdu},
		},
		{
			"Az elmúlt hétvégén 12-re emelkedett az elhunyt koronavírus-fertőzöttek száma Szlovákiában. Mindegyik szociális otthon dolgozóját letesztelik, Matovič szerint az ingázóknak még várniuk kellene a teszteléssel",
			[]Language{Hungarian, Slovak},
		},
	}
	for _, testCase := range testCases {
		detector := NewLanguageDetectorBuilder().
			FromLanguages(testCase.languages...).
			WithPreloadedLanguageModels().
			Build()
		detectedLanguages := make(map[Language]bool)
		for i := 0; i < 100; i++ {
			language, _ := detector.DetectLanguageOf(testCase.text)
			detectedLanguages[language] = true
		}
		assert.Len(
			t,
			detectedLanguages,
			1,
			fmt.Sprintf("language detector is non-deterministic for languages %v", testCase.languages),
		)
	}
}

func BenchmarkLanguageDetection(b *testing.B) {
	detector := NewLanguageDetectorBuilder().
		FromAllLanguages().
		WithPreloadedLanguageModels().
		Build()
	sentences := []string{
		"ربما يبتعد العقرب عن بعض الذين يخيبون أمله، أو يشعر بالحاجة إلى الانتقاء، وعدم البحث عن النشاطات التي ترهق أكثر مما تسعده.",
		"Επί της ουσίας τόσο οι υφιστάμενες ενισχύσεις που οφείλονται στους κτηνοτρόφους όσο και αυτές της νέας προγραμματικής περιόδου παραμένουν στον αέρα.",
		"It has three co-chairs, one from each of a provincial health and agriculture department, and a third from the federal government.",
		"અશ્વિની ભટ્ટની નવલકથામાંથી થોડુંક માણસ જ્યારે વેદનાની પરાકાષ્ટાની સીમા વટાવી જાય પછી એક એવી પરિસ્થિતિ આવે છે જ્યારે દર્દ-વેદના નથી રહેતી, વેદના છે કે નહિ તેનો પણ કોઇ ખ્યાલ નથી રહેતો.",
		"・京都大学施設に電離圏における電子数などの状況を取得可能なイオノゾンデ受信機（斜入射観測装置）を設置することで、新たな観測手法が地震先行現象検出に資するかを検証する。",
		"ამასთანავე წანარები სათავეში უდგებიან (თუ წარმართავენ?) კახეთის გაერთიანებისა და ერთიანი სამთავროს ჩამოყალიბების პროცესს.",
		"하지만 금융 전문가들은 “전체 대출 중 부동산 PF로의 쏠림 현상이 심각한 상태에서 각종 대출 규제로 자금 여력이 부족해질 경우 연체율이 높아질 수 있는데 당국이 안이하게 대응하는 측면이 있다”고 지적했다.",
		"И потому я должен возблагодарить провидение; если бы не провидение, то сердце твое, бедный сэр Пол, все конечно разбилось бы.",
		"ส.บัญชีรายชื่อ พรรคเพื่อไทย แต่อยู่ในระหว่างการตัดสินเรื่องการเป็นสมาชิกภาพของพรรคการเมือง เพราะถูกคุมขังโดยหมายศาล ระหว่างการสมัครรับเลือกตั้ง ซึ่งขณะนี้อยู่ในระหว่างการพิจารณาของ กกต.",
		"人们必须面对：遭受严重破坏的自然生态；大自然反扑所造成的天灾人祸；人口快速成长的沈重压力；生存竞争日异严峻的社会情况；传统家庭结构逐渐瓦解的隐忧，社会价值观念混淆等问题。",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, sentence := range sentences {
			detector.DetectLanguageOf(sentence)
		}
	}
}
