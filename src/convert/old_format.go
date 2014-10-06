package main

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type OldForm struct {
	Title     string        `json:"title"`
	Published bool          `json:"published"`
	Id        string        `json:"id"`
	GroupName string        `json:"group"`
	GroupId   string        `json:"gid"`
	Code      string        `json:"code"`
	CreatedOn *UnixTime     `json:"created_on"`
	Version   string        `json:"version"`
	Data      QuestionArray `json:"data"`
}

type OldQuestion struct {
	SId         string          `json:"s-id"`
	SType       string          `json:"s-type"`
	DataTheme   string          `json:"data-theme"`
	Html        string          `json:"html"`
	SItems      json.RawMessage `json:"s-items"`
	SValidate   OldValidation   `json:"s-validate"`
	SDirection  string          `json:"s-direction"`
	SLabel      string          `json:"s-label"`
	SGroup      string          `json:"s-group"`
	SBlock      bool            `json:"s-block"`
	Placeholder string          `json:"placeholder"`
	MaxLength   IntString       `json:"maxlength"`
	Tag         string          `json:"tag"`
	SInline     bool            `json:"inline"`
	SMaximize   bool            `json:"maximize"`
	SSuffix     *OldQuestion    `json:"s-suffix"`
	Disabled    string          `json:"disabled"`
	Type        string          `json:"type"`
	Value       interface{}     `json:"value"`
	SEmpty      bool            `json:"s-empty"`
	SStore      string          `json:"s-store"`
	SSort       bool            `json:"s-sort"`
	SItem       *OldStoreItem   `json:"s-item"`
	SDepend     interface{}     `json:"s-depend"`
	Min         IntString       `json:"min"`
	Max         IntString       `json:"max"`
}

type OldValidation struct {
	Required bool `json:"required"`
	Skip     bool `json:"skip"`
}

type OldStoreItem struct {
	SType  string `json:"s-type"`
	SLabel string `json:"s-label"`
	SValue string `json:"s-value"`
}

type QuestionArray []OldQuestion

func (qa *QuestionArray) UnmarshalJSON(data []byte) (err error) {
	runes := bytes.Runes(data[0:1])

	var questions []OldQuestion

	if runes[0] == '[' {
		err = json.Unmarshal(data, &questions)
	} else {
		var old_question OldQuestion
		err = json.Unmarshal(data, &old_question)
		questions = []OldQuestion{old_question}
	}

	if err == nil {
		*qa = QuestionArray(questions)
	}

	return
}

type IntString int

func (ci *IntString) UnmarshalJSON(data []byte) (err error) {
	runes := bytes.Runes(data[0:1])
	var (
		str   string
		value int
	)

	if runes[0] == '"' {
		//convert string to int
		str = string(data[1 : len(data)-1])
	} else {
		str = string(data)
	}

	value, err = strconv.Atoi(str)
	*ci = IntString(value)

	return
}

type UnixTime time.Time

func (ux *UnixTime) UnmarshalJSON(data []byte) (err error) {
	var (
		int_value  int64
		time_value time.Time
		str_value  string = string(data)
	)

	int_value, err = strconv.ParseInt(str_value, 10, 64)
	if err != nil {
		return
	}

	time_value = time.Unix(int_value/1000, int_value%1000)

	*ux = UnixTime(time_value)

	log.Println("unix time = ", int_value, "time=", time_value.String())

	return
}
