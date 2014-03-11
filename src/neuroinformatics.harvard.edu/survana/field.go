package survana

type Field struct {
	Id          string      `bson:"id,omitempty" json:"id,omitempty"`
	Type        string      `bson:"type,omitempty" json:"type,omitempty"`               //text, input, number, etc
	Group       string      `bson:"group,omitempty" json:"group,omitempty"`             //radiogroup,checkboxgroup
	Align       string      `bson:"align,omitempty" json:"align,omitempty"`             //left, center, right
	Html        string      `bson:"html,omitempty" json:"html,omitempty"`               //for text fields
	Label       *LabelField `bson:"label,omitempty" json:"label,omitempty"`             //for labeled controls
	Placeholder string      `bson:"placeholder,omitempty" json:"placeholder,omitempty"` //for input controls
	Prefix      *Field      `bson:"prefix,omitempty" json:"prefix,omitempty"`           //any field
	Suffix      *Field      `bson:"suffix,omitempty" json:"suffix,omitempty"`           //any field
	Size        *FieldSize  `bson:"size,omitempty" json:"size,omitempty"`               //max: 12
	Note        string      `bson:"note,omitempty" json:"note,omitempty"`               //comment
	Fields      *[]Field    `bson:"fields,omitempty" json:"fields,omitempty"`           //radios in radiogroups, etc
	Value       interface{} `bson:"value,omitempty" json:"value,omitempty"`             //value (string or number)
	//matrix
	Matrix   string   `bson:"matrix,omitempty" json:"matrix,omitempty"`      //the type of matrix  (radio,checkbox)
	Striped  bool     `bson:"striped,omitempty" json:"striped,omitempty"`    //whether the matrix should be striped
	Hover    bool     `bson:"hover,omitempty" json:"hover,omitempty"`        //whether the matrix should react on mouse over
	NoAnswer bool     `bson:"no_answer,omitempty" json:"noanswer,omitempty"` //append Prefer Not to Answer
	Numbers  bool     `bson:"numbers,omitempty" json:"numbers,omitempty"`    //prepend index numbers to label
	Columns  *[]Field `bson:"columns,omitempty" json:"columns,omitempty"`    //table headers
	Rows     *[]Field `bson:"rows,omitempty" json:"rows,omitempty"`          //table rows
	Equalize bool     `bson:"equalize,omitempty" json:"equalize,omitempty"`  //make all rows equal height
    Validation  map[string]interface{} `bson:"validation,omitempty" json:"validation,omitempty"`
}

type LabelField struct {
	Html  string     `bson:"html,omitempty" json:"html,omitempty"`   //any valid html string
	Align string     `bson:"align,omitempty" json:"align,omitempty"` //left, center, right
	Size  *FieldSize `bson:"size,omitempty" json:"size,omitempty"`   //12 for vertical alignment
}

//12 is max
type FieldSize struct {
	Large      int `bson:"l,omitempty" json:"l,omitempty"`   //col-l-
	Medium     int `bson:"m,omitempty" json:"m,omitempty"`   //col-m-
	Small      int `bson:"s,omitempty" json:"s,omitempty"`   //col-s-
	ExtraSmall int `bson:"xs,omitempty" json:"xs,omitempty"` //col-xs-
}
